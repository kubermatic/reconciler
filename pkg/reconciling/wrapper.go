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
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	ctrlruntimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// OwnerRefWrapper is responsible for wrapping a ObjectReconciler function, solely to set the OwnerReference to the cluster object.
func OwnerRefWrapper(ref metav1.OwnerReference) ObjectModifier {
	return func(create ObjectReconciler) ObjectReconciler {
		return func(existing ctrlruntimeclient.Object) (ctrlruntimeclient.Object, error) {
			obj, err := create(existing)
			if err != nil {
				return obj, err
			}

			obj.(metav1.Object).SetOwnerReferences([]metav1.OwnerReference{ref})
			return obj, nil
		}
	}
}

// ImagePullSecretsWrapper is generating a new ObjectModifier that wraps an ObjectReconciler
// and takes care of adding the secret names provided to the ImagePullSecrets.
//
// TODO At the moment only Deployments are supported, but
// this can be extended to whatever Object carrying a PodSpec.
func ImagePullSecretsWrapper(secretNames ...string) ObjectModifier {
	return func(create ObjectReconciler) ObjectReconciler {
		return func(existing ctrlruntimeclient.Object) (ctrlruntimeclient.Object, error) {
			obj, err := create(existing)
			if err != nil {
				return obj, err
			}
			if len(secretNames) == 0 {
				return obj, nil
			}
			switch o := obj.(type) {
			case *appsv1.Deployment:
				configureImagePullSecrets(&o.Spec.Template.Spec, secretNames)
				return o, nil
			default:
				return o, fmt.Errorf(`type %q is not supported by ImagePullSecretModifier`, o.GetObjectKind().GroupVersionKind())
			}
		}
	}
}

func configureImagePullSecrets(podSpec *corev1.PodSpec, secretNames []string) {
	// Only configure image pull secrets when provided in the configuration.
	currentSecretNames := sets.NewString()
	for _, ips := range podSpec.ImagePullSecrets {
		currentSecretNames.Insert(ips.Name)
	}
	for _, s := range secretNames {
		if !currentSecretNames.Has(s) {
			podSpec.ImagePullSecrets = append(podSpec.ImagePullSecrets, corev1.LocalObjectReference{Name: s})
		}
	}
}
