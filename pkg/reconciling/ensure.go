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
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	ctrlruntimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// ObjectReconciler defines an interface to create/update a ctrlruntimeclient.Object.
type ObjectReconciler = func(existing ctrlruntimeclient.Object) (ctrlruntimeclient.Object, error)

// ObjectModifier is a wrapper function which modifies the object which gets returned by the passed in ObjectReconciler.
type ObjectModifier func(create ObjectReconciler) ObjectReconciler

func CreateWithNamespace(reconciler ObjectReconciler, namespace string) ObjectReconciler {
	return func(existing ctrlruntimeclient.Object) (ctrlruntimeclient.Object, error) {
		obj, err := reconciler(existing)
		if err != nil {
			return nil, err
		}
		obj.(metav1.Object).SetNamespace(namespace)
		return obj, nil
	}
}

func CreateWithName(reconciler ObjectReconciler, name string) ObjectReconciler {
	return func(existing ctrlruntimeclient.Object) (ctrlruntimeclient.Object, error) {
		obj, err := reconciler(existing)
		if err != nil {
			return nil, err
		}
		obj.(metav1.Object).SetName(name)
		return obj, nil
	}
}

func objectLogger(obj ctrlruntimeclient.Object) *zap.SugaredLogger {
	// make sure we handle objects with broken typeMeta and still create a nice-looking kind name
	logger := log.Logger().With("kind", reflect.TypeOf(obj).Elem())
	if ns := obj.GetNamespace(); ns != "" {
		logger = logger.With("namespace", ns)
	}

	// ensure name comes after namespace
	return logger.With("name", obj.GetName())
}

// EnsureNamedObject will generate the Object with the passed create function & create or update it in Kubernetes if necessary.
func EnsureNamedObject(ctx context.Context, namespacedName types.NamespacedName, rawReconciler ObjectReconciler, client ctrlruntimeclient.Client, emptyObject ctrlruntimeclient.Object, requiresRecreate bool) error {
	// A wrapper to ensure we always set the Namespace and Name. This is useful as we call reconciler twice
	reconciler := CreateWithNamespace(rawReconciler, namespacedName.Namespace)
	reconciler = CreateWithName(reconciler, namespacedName.Name)

	exists := true
	existingObject := emptyObject.DeepCopyObject().(ctrlruntimeclient.Object)
	if err := client.Get(ctx, namespacedName, existingObject); err != nil {
		if !apierrors.IsNotFound(err) {
			return fmt.Errorf("failed to get Object(%T): %w", existingObject, err)
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
		obj, err := reconciler(emptyObject)
		if err != nil {
			return fmt.Errorf("failed to generate object: %w", err)
		}
		if err := client.Create(ctx, obj); err != nil {
			return fmt.Errorf("failed to create %T '%s': %w", obj, namespacedName.String(), err)
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
	obj, err := reconciler(existingObject.DeepCopyObject().(ctrlruntimeclient.Object))
	if err != nil {
		return fmt.Errorf("failed to build Object(%T) '%s': %w", existingObject, namespacedName.String(), err)
	}

	if compare.DeepEqual(obj.(metav1.Object), existingObject.(metav1.Object)) {
		return nil
	}

	if !requiresRecreate {
		// We keep resetting the status here to avoid working on any outdated object
		// and all objects are up-to-date once a reconcile process starts.
		switch v := obj.(type) {
		case *appsv1.StatefulSet:
			v.Status.Reset()
		case *appsv1.Deployment:
			v.Status.Reset()
		}

		if err := client.Update(ctx, obj); err != nil {
			return fmt.Errorf("failed to update object %T %q: %w", obj, namespacedName.String(), err)
		}
	} else {
		if err := client.Delete(ctx, obj.DeepCopyObject().(ctrlruntimeclient.Object)); err != nil {
			return fmt.Errorf("failed to delete object %T %q: %w", obj, namespacedName.String(), err)
		}

		obj.SetResourceVersion("")
		obj.SetUID("")
		obj.SetGeneration(0)

		if err := client.Create(ctx, obj); err != nil {
			return fmt.Errorf("failed to create object %T %q: %w", obj, namespacedName.String(), err)
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
