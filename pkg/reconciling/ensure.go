/*
Copyright 2022 The Kubermatic Kubernetes Platform contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package reconciling

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"go.uber.org/zap"

	"k8c.io/reconciler/pkg/compare"
	"k8c.io/reconciler/pkg/log"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	policyv1 "k8s.io/api/policy/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	ctrlruntimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// ObjectReconciler defines an interface to create/update a ctrlruntimeclient.Object.
type ObjectReconciler = func(existing ctrlruntimeclient.Object) (ctrlruntimeclient.Object, error)

// ObjectModifier is a wrapper function which modifies the object which gets returned by the passed in ObjectReconciler.
type ObjectModifier func(create ObjectReconciler) ObjectReconciler

// GenericObjectReconciler defines an interface to create/update a T.
type GenericObjectReconciler[T any] func(existing T) (T, error)

// GenericNamedObjectReconciler is a wrapper function which modifies the object which gets returned by the passed in GenericObjectReconciler.
type GenericNamedObjectReconciler[T any] func() (string, GenericObjectReconciler[T])

// GenericObjectModifier is a wrapper function which modifies the object which gets returned by the passed in ObjectReconciler.
type GenericObjectModifier[T any] func(create GenericObjectReconciler[T]) GenericObjectReconciler[T]

func objectLogger(obj ctrlruntimeclient.Object) *zap.SugaredLogger {
	// make sure we handle objects with broken typeMeta and still create a nice-looking kind name
	logger := log.Logger().With("kind", objectKind(obj))
	if ns := obj.GetNamespace(); ns != "" {
		logger = logger.With("namespace", ns)
	}

	// ensure name comes after namespace
	return logger.With("name", obj.GetName())
}

func requiresRecreate(x interface{}) bool {
	_, ok := x.(*policyv1.PodDisruptionBudget)
	return ok
}

func objectKind(obj ctrlruntimeclient.Object) string {
	if u, ok := obj.(*unstructured.Unstructured); ok {
		gvk := u.GetObjectKind().GroupVersionKind()

		return fmt.Sprintf("%s/%s.%s", gvk.Group, gvk.Version, gvk.Kind)
	}

	// make sure we handle objects with broken typeMeta and still create a nice-looking kind name
	return reflect.TypeOf(obj).Elem().String()
}

func objectName(name, namespace string) string {
	result := name
	if namespace != "" {
		result = fmt.Sprintf("%s/%s", namespace, name)
	}

	return result
}

// EnsureNamedObjects will call EnsureNamedObject for each of the given reconciler functions.
func EnsureNamedObjects[T ctrlruntimeclient.Object](
	ctx context.Context,
	client ctrlruntimeclient.Client,
	namespace string,
	reconcilerFactories []GenericNamedObjectReconciler[T],
	objectModifiers ...ObjectModifier,
) error {
	for _, factory := range reconcilerFactories {
		name, reconciler := factory()

		// create a new instance of the type represented by T
		var placeholder T
		emptyObject := reflect.New(reflect.TypeOf(placeholder).Elem()).Interface().(T)

		if err := EnsureNamedObject(ctx, client, types.NamespacedName{Namespace: namespace, Name: name}, emptyObject, reconciler, objectModifiers...); err != nil {
			return fmt.Errorf("failed to ensure %s %q: %w", objectKind(emptyObject), objectName(name, namespace), err)
		}
	}

	return nil
}

func makeGenericObjectReconciler[T ctrlruntimeclient.Object](wrapThis ObjectReconciler) GenericObjectReconciler[T] {
	return func(existing T) (T, error) {
		result, err := wrapThis(existing)
		return result.(T), err
	}
}

func makeInterfacedObjectReconciler[T ctrlruntimeclient.Object](wrapThis GenericObjectReconciler[T]) ObjectReconciler {
	return func(existing ctrlruntimeclient.Object) (ctrlruntimeclient.Object, error) {
		return wrapThis(existing.(T))
	}
}

func wrapGenericObjectModifier[T ctrlruntimeclient.Object](wrapThis GenericObjectReconciler[T], wrapper ObjectModifier) GenericObjectReconciler[T] {
	// convert the generic reconciler into an interfaced reconciler
	modifieddReconciler := makeInterfacedObjectReconciler(wrapThis)

	// wrap the interfaced reconciler with the object modifier
	modifieddReconciler = wrapper(modifieddReconciler)

	// convert back to a generic reconciler
	return makeGenericObjectReconciler[T](modifieddReconciler)
}

// EnsureNamedObject will generate the Object with the passed create function & create or update it in Kubernetes if necessary.
func EnsureNamedObject[T ctrlruntimeclient.Object](
	ctx context.Context,
	client ctrlruntimeclient.Client,
	namespacedName types.NamespacedName,
	existingObject T,
	reconciler GenericObjectReconciler[T],
	objectModifiers ...ObjectModifier,
) error {
	// ensure that whatever a reconciler does or does not, we will always
	// set the name and namespace (i.e. the reconciler cannot control these)
	reconcilerWithIdentity := func(obj T) (T, error) {
		created, err := reconciler(obj)
		if err != nil {
			return obj, err
		}

		created.SetName(namespacedName.Name)
		created.SetNamespace(namespacedName.Namespace)

		return created, nil
	}

	// ensure we always apply the default values
	defaultedReconciler := applyDefaultValues(reconcilerWithIdentity)

	// wrap the reconciler in the modifiers, for this it's necessary to convert
	// the interface-based modifiers into generic modifiers.
	for _, objectModifier := range objectModifiers {
		defaultedReconciler = wrapGenericObjectModifier(defaultedReconciler, objectModifier)
	}

	// check if the object exists already
	exists := true
	if err := client.Get(ctx, namespacedName, existingObject); err != nil {
		if !apierrors.IsNotFound(err) {
			return fmt.Errorf("failed to get %s: %w", objectKind(existingObject), err)
		}
		exists = false
	}

	if exists && pauseAnnotation != "" {
		annotations := existingObject.GetAnnotations()
		if v, ok := annotations[pauseAnnotation]; ok && v == "true" {
			objectLogger(existingObject).Warn("not touching paused resource")
			return nil
		}
	}

	// Object does not exist in lister -> Create the Object
	if !exists {
		obj, err := defaultedReconciler(existingObject)
		if err != nil {
			return fmt.Errorf("failed to generate object: %w", err)
		}
		if err := client.Create(ctx, obj); err != nil {
			return fmt.Errorf("failed to create %s %q: %w", objectKind(existingObject), namespacedName.String(), err)
		}

		// Wait until the object exists in the cache
		createdObjectIsInCache := WaitUntilObjectExistsInCacheConditionFunc(ctx, client, objectLogger(obj), namespacedName, obj)
		err = wait.PollImmediate(10*time.Millisecond, 10*time.Second, createdObjectIsInCache)
		if err != nil {
			return fmt.Errorf("failed waiting for the cache to contain our newly created object: %w", err)
		}

		objectLogger(obj).Info("created new resource")
		return nil
	}

	// Create a copy to make sure we don't compare the object onto itself
	// in case the reconciler returns the same pointer it got passed in
	obj, err := defaultedReconciler(existingObject.DeepCopyObject().(T))
	if err != nil {
		return fmt.Errorf("failed to build %s %q: %w", objectKind(existingObject), namespacedName.String(), err)
	}

	if compare.DeepEqual(obj, existingObject) {
		return nil
	}

	if !requiresRecreate(obj) {
		// We keep resetting the status here to avoid working on any outdated object
		// and all objects are up-to-date once a reconcile process starts.
		switch v := any(obj).(type) {
		case *appsv1.StatefulSet:
			v.Status.Reset()
		case *appsv1.Deployment:
			v.Status.Reset()
		}

		if err := client.Update(ctx, obj); err != nil {
			return fmt.Errorf("failed to update object %s %q: %w", objectKind(existingObject), namespacedName.String(), err)
		}
	} else {
		if err := client.Delete(ctx, obj.DeepCopyObject().(ctrlruntimeclient.Object)); err != nil {
			return fmt.Errorf("failed to delete object %s %q: %w", objectKind(existingObject), namespacedName.String(), err)
		}

		obj.SetResourceVersion("")
		obj.SetUID("")
		obj.SetGeneration(0)

		if err := client.Create(ctx, obj); err != nil {
			return fmt.Errorf("failed to create object %s %q: %w", objectKind(existingObject), namespacedName.String(), err)
		}
	}

	// Wait until the object we retrieve via "client.Get" has a different ResourceVersion than the old object
	updatedObjectIsInCache := waitUntilUpdateIsInCacheConditionFunc(ctx, client, objectLogger(obj), namespacedName, existingObject)
	err = wait.PollImmediate(10*time.Millisecond, 10*time.Second, updatedObjectIsInCache)
	if err != nil {
		return fmt.Errorf("failed waiting for the cache to contain our latest changes: %w", err)
	}

	objectLogger(obj).Info("updated resource")

	return nil
}

func applyDefaultValues[T ctrlruntimeclient.Object](reconciler GenericObjectReconciler[T]) GenericObjectReconciler[T] {
	return func(t T) (T, error) {
		oldObj := t.DeepCopyObject().(T)

		created, err := reconciler(t)
		if err != nil {
			return t, err
		}

		var defaulted ctrlruntimeclient.Object = created

		switch v := any(oldObj).(type) {
		case *appsv1.Deployment:
			defaulted, err = DefaultDeployment(v, any(created).(*appsv1.Deployment))

		case *appsv1.StatefulSet:
			defaulted, err = DefaultStatefulSet(v, any(created).(*appsv1.StatefulSet))

		case *appsv1.DaemonSet:
			defaulted, err = DefaultDaemonSet(v, any(created).(*appsv1.DaemonSet))

		case *batchv1.CronJob:
			defaulted, err = DefaultCronJob(v, any(created).(*batchv1.CronJob))
		}

		return defaulted.(T), err
	}
}

func waitUntilUpdateIsInCacheConditionFunc(
	ctx context.Context,
	client ctrlruntimeclient.Client,
	log *zap.SugaredLogger,
	namespacedName types.NamespacedName,
	oldObj ctrlruntimeclient.Object,
) wait.ConditionFunc {
	return func() (bool, error) {
		// Create a copy to have something which we can pass into the client
		currentObj := oldObj.DeepCopyObject().(ctrlruntimeclient.Object)

		if err := client.Get(ctx, namespacedName, currentObj); err != nil {
			log.Errorw("failed retrieving object while waiting for the cache to contain our latest changes", zap.Error(err))
			return false, nil
		}

		// Check if the object from the store differs the old object;
		// We are waiting for the resourceVersion/generation to change
		// and for this new version to be then present in our cache.
		if !compare.DeepEqual(currentObj.(metav1.Object), oldObj.(metav1.Object)) {
			return true, nil
		}

		return false, nil
	}
}

func WaitUntilObjectExistsInCacheConditionFunc(
	ctx context.Context,
	client ctrlruntimeclient.Client,
	log *zap.SugaredLogger,
	namespacedName types.NamespacedName,
	obj ctrlruntimeclient.Object,
) wait.ConditionFunc {
	return func() (bool, error) {
		newObj := obj.DeepCopyObject().(ctrlruntimeclient.Object)
		if err := client.Get(ctx, namespacedName, newObj); err != nil {
			if apierrors.IsNotFound(err) {
				return false, nil
			}

			log.Errorw("failed retrieving object while waiting for the cache to contain our newly created object", zap.Error(err))
			return false, nil
		}

		return true, nil
	}
}
