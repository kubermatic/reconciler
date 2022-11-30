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

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	ctrlruntimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// UnstructuredReconciler defines an interface to create/update Unstructureds.
type UnstructuredReconciler = func(existing *unstructured.Unstructured) (*unstructured.Unstructured, error)

// NamedUnstructuredReconcilerFactory returns the name, kind and API version of the resource and the
// corresponding reconciler function.
type NamedUnstructuredReconcilerFactory = func() (name, kind, apiVersion string, create UnstructuredReconciler)

// EnsureNamedUnstructuredObjects will call EnsureNamedObject for each of the given reconciler functions.
// This is a dedicated function (the sister of EnsureNamedObjects) because unstructured objects require
// us to manually inject the API version and kind, which is not possible with the regular ObjectReconcilers,
// as they cannot provide this information.
func EnsureNamedUnstructuredObjects(
	ctx context.Context,
	client ctrlruntimeclient.Client,
	namespace string,
	creatorFactories []NamedUnstructuredReconcilerFactory,
	objectModifiers ...ObjectModifier,
) error {
	for _, factory := range creatorFactories {
		name, kind, apiVersion, reconciler := factory()

		emptyObj := &unstructured.Unstructured{}
		emptyObj.SetAPIVersion(apiVersion)
		emptyObj.SetKind(kind)

		if err := EnsureNamedObject(ctx, client, types.NamespacedName{Namespace: namespace, Name: name}, emptyObj, reconciler, objectModifiers...); err != nil {
			return fmt.Errorf("failed to ensure %s %q: %w", objectKind(emptyObj), objectName(name, namespace), err)
		}
	}

	return nil
}
