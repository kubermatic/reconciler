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

	"k8s.io/apimachinery/pkg/types"
	ctrlruntimeclient "sigs.k8s.io/controller-runtime/pkg/client"

	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	networkingv1 "k8s.io/api/networking/v1"
	policyv1 "k8s.io/api/policy/v1"
	rbacv1 "k8s.io/api/rbac/v1"
)

// NamespaceReconciler defines an interface to create/update Namespaces.
type NamespaceReconciler = func(existing *corev1.Namespace) (*corev1.Namespace, error)

// NamedNamespaceReconcilerFactory returns the name of the resource and the corresponding Reconciler function.
type NamedNamespaceReconcilerFactory = func() (name string, reconciler NamespaceReconciler)

// NamespaceObjectWrapper adds a wrapper so the NamespaceReconciler matches ObjectReconciler.
// This is needed as Go does not support function interface matching.
func NamespaceObjectWrapper(reconciler NamespaceReconciler) ObjectReconciler {
	return func(existing ctrlruntimeclient.Object) (ctrlruntimeclient.Object, error) {
		if existing != nil {
			return reconciler(existing.(*corev1.Namespace))
		}
		return reconciler(&corev1.Namespace{})
	}
}

// ReconcileNamespaces will create and update the Namespaces coming from the passed NamespaceReconciler slice.
func ReconcileNamespaces(ctx context.Context, namedFactories []NamedNamespaceReconcilerFactory, namespace string, client ctrlruntimeclient.Client, objectModifiers ...ObjectModifier) error {
	for _, factory := range namedFactories {
		name, reconciler := factory()
		reconcileObject := NamespaceObjectWrapper(reconciler)
		reconcileObject = CreateWithNamespace(reconcileObject, namespace)
		reconcileObject = CreateWithName(reconcileObject, name)

		for _, objectModifier := range objectModifiers {
			reconcileObject = objectModifier(reconcileObject)
		}

		if err := EnsureNamedObject(ctx, types.NamespacedName{Namespace: namespace, Name: name}, reconcileObject, client, &corev1.Namespace{}, false); err != nil {
			return fmt.Errorf("failed to ensure Namespace %s/%s: %w", namespace, name, err)
		}
	}

	return nil
}

// ServiceReconciler defines an interface to create/update Services.
type ServiceReconciler = func(existing *corev1.Service) (*corev1.Service, error)

// NamedServiceReconcilerFactory returns the name of the resource and the corresponding Reconciler function.
type NamedServiceReconcilerFactory = func() (name string, reconciler ServiceReconciler)

// ServiceObjectWrapper adds a wrapper so the ServiceReconciler matches ObjectReconciler.
// This is needed as Go does not support function interface matching.
func ServiceObjectWrapper(reconciler ServiceReconciler) ObjectReconciler {
	return func(existing ctrlruntimeclient.Object) (ctrlruntimeclient.Object, error) {
		if existing != nil {
			return reconciler(existing.(*corev1.Service))
		}
		return reconciler(&corev1.Service{})
	}
}

// ReconcileServices will create and update the Services coming from the passed ServiceReconciler slice.
func ReconcileServices(ctx context.Context, namedFactories []NamedServiceReconcilerFactory, namespace string, client ctrlruntimeclient.Client, objectModifiers ...ObjectModifier) error {
	for _, factory := range namedFactories {
		name, reconciler := factory()
		reconcileObject := ServiceObjectWrapper(reconciler)
		reconcileObject = CreateWithNamespace(reconcileObject, namespace)
		reconcileObject = CreateWithName(reconcileObject, name)

		for _, objectModifier := range objectModifiers {
			reconcileObject = objectModifier(reconcileObject)
		}

		if err := EnsureNamedObject(ctx, types.NamespacedName{Namespace: namespace, Name: name}, reconcileObject, client, &corev1.Service{}, false); err != nil {
			return fmt.Errorf("failed to ensure Service %s/%s: %w", namespace, name, err)
		}
	}

	return nil
}

// SecretReconciler defines an interface to create/update Secrets.
type SecretReconciler = func(existing *corev1.Secret) (*corev1.Secret, error)

// NamedSecretReconcilerFactory returns the name of the resource and the corresponding Reconciler function.
type NamedSecretReconcilerFactory = func() (name string, reconciler SecretReconciler)

// SecretObjectWrapper adds a wrapper so the SecretReconciler matches ObjectReconciler.
// This is needed as Go does not support function interface matching.
func SecretObjectWrapper(reconciler SecretReconciler) ObjectReconciler {
	return func(existing ctrlruntimeclient.Object) (ctrlruntimeclient.Object, error) {
		if existing != nil {
			return reconciler(existing.(*corev1.Secret))
		}
		return reconciler(&corev1.Secret{})
	}
}

// ReconcileSecrets will create and update the Secrets coming from the passed SecretReconciler slice.
func ReconcileSecrets(ctx context.Context, namedFactories []NamedSecretReconcilerFactory, namespace string, client ctrlruntimeclient.Client, objectModifiers ...ObjectModifier) error {
	for _, factory := range namedFactories {
		name, reconciler := factory()
		reconcileObject := SecretObjectWrapper(reconciler)
		reconcileObject = CreateWithNamespace(reconcileObject, namespace)
		reconcileObject = CreateWithName(reconcileObject, name)

		for _, objectModifier := range objectModifiers {
			reconcileObject = objectModifier(reconcileObject)
		}

		if err := EnsureNamedObject(ctx, types.NamespacedName{Namespace: namespace, Name: name}, reconcileObject, client, &corev1.Secret{}, false); err != nil {
			return fmt.Errorf("failed to ensure Secret %s/%s: %w", namespace, name, err)
		}
	}

	return nil
}

// ConfigMapReconciler defines an interface to create/update ConfigMaps.
type ConfigMapReconciler = func(existing *corev1.ConfigMap) (*corev1.ConfigMap, error)

// NamedConfigMapReconcilerFactory returns the name of the resource and the corresponding Reconciler function.
type NamedConfigMapReconcilerFactory = func() (name string, reconciler ConfigMapReconciler)

// ConfigMapObjectWrapper adds a wrapper so the ConfigMapReconciler matches ObjectReconciler.
// This is needed as Go does not support function interface matching.
func ConfigMapObjectWrapper(reconciler ConfigMapReconciler) ObjectReconciler {
	return func(existing ctrlruntimeclient.Object) (ctrlruntimeclient.Object, error) {
		if existing != nil {
			return reconciler(existing.(*corev1.ConfigMap))
		}
		return reconciler(&corev1.ConfigMap{})
	}
}

// ReconcileConfigMaps will create and update the ConfigMaps coming from the passed ConfigMapReconciler slice.
func ReconcileConfigMaps(ctx context.Context, namedFactories []NamedConfigMapReconcilerFactory, namespace string, client ctrlruntimeclient.Client, objectModifiers ...ObjectModifier) error {
	for _, factory := range namedFactories {
		name, reconciler := factory()
		reconcileObject := ConfigMapObjectWrapper(reconciler)
		reconcileObject = CreateWithNamespace(reconcileObject, namespace)
		reconcileObject = CreateWithName(reconcileObject, name)

		for _, objectModifier := range objectModifiers {
			reconcileObject = objectModifier(reconcileObject)
		}

		if err := EnsureNamedObject(ctx, types.NamespacedName{Namespace: namespace, Name: name}, reconcileObject, client, &corev1.ConfigMap{}, false); err != nil {
			return fmt.Errorf("failed to ensure ConfigMap %s/%s: %w", namespace, name, err)
		}
	}

	return nil
}

// ServiceAccountReconciler defines an interface to create/update ServiceAccounts.
type ServiceAccountReconciler = func(existing *corev1.ServiceAccount) (*corev1.ServiceAccount, error)

// NamedServiceAccountReconcilerFactory returns the name of the resource and the corresponding Reconciler function.
type NamedServiceAccountReconcilerFactory = func() (name string, reconciler ServiceAccountReconciler)

// ServiceAccountObjectWrapper adds a wrapper so the ServiceAccountReconciler matches ObjectReconciler.
// This is needed as Go does not support function interface matching.
func ServiceAccountObjectWrapper(reconciler ServiceAccountReconciler) ObjectReconciler {
	return func(existing ctrlruntimeclient.Object) (ctrlruntimeclient.Object, error) {
		if existing != nil {
			return reconciler(existing.(*corev1.ServiceAccount))
		}
		return reconciler(&corev1.ServiceAccount{})
	}
}

// ReconcileServiceAccounts will create and update the ServiceAccounts coming from the passed ServiceAccountReconciler slice.
func ReconcileServiceAccounts(ctx context.Context, namedFactories []NamedServiceAccountReconcilerFactory, namespace string, client ctrlruntimeclient.Client, objectModifiers ...ObjectModifier) error {
	for _, factory := range namedFactories {
		name, reconciler := factory()
		reconcileObject := ServiceAccountObjectWrapper(reconciler)
		reconcileObject = CreateWithNamespace(reconcileObject, namespace)
		reconcileObject = CreateWithName(reconcileObject, name)

		for _, objectModifier := range objectModifiers {
			reconcileObject = objectModifier(reconcileObject)
		}

		if err := EnsureNamedObject(ctx, types.NamespacedName{Namespace: namespace, Name: name}, reconcileObject, client, &corev1.ServiceAccount{}, false); err != nil {
			return fmt.Errorf("failed to ensure ServiceAccount %s/%s: %w", namespace, name, err)
		}
	}

	return nil
}

// EndpointsReconciler defines an interface to create/update Endpoints.
type EndpointsReconciler = func(existing *corev1.Endpoints) (*corev1.Endpoints, error)

// NamedEndpointsReconcilerFactory returns the name of the resource and the corresponding Reconciler function.
type NamedEndpointsReconcilerFactory = func() (name string, reconciler EndpointsReconciler)

// EndpointsObjectWrapper adds a wrapper so the EndpointsReconciler matches ObjectReconciler.
// This is needed as Go does not support function interface matching.
func EndpointsObjectWrapper(reconciler EndpointsReconciler) ObjectReconciler {
	return func(existing ctrlruntimeclient.Object) (ctrlruntimeclient.Object, error) {
		if existing != nil {
			return reconciler(existing.(*corev1.Endpoints))
		}
		return reconciler(&corev1.Endpoints{})
	}
}

// ReconcileEndpoints will create and update the Endpoints coming from the passed EndpointsReconciler slice.
func ReconcileEndpoints(ctx context.Context, namedFactories []NamedEndpointsReconcilerFactory, namespace string, client ctrlruntimeclient.Client, objectModifiers ...ObjectModifier) error {
	for _, factory := range namedFactories {
		name, reconciler := factory()
		reconcileObject := EndpointsObjectWrapper(reconciler)
		reconcileObject = CreateWithNamespace(reconcileObject, namespace)
		reconcileObject = CreateWithName(reconcileObject, name)

		for _, objectModifier := range objectModifiers {
			reconcileObject = objectModifier(reconcileObject)
		}

		if err := EnsureNamedObject(ctx, types.NamespacedName{Namespace: namespace, Name: name}, reconcileObject, client, &corev1.Endpoints{}, false); err != nil {
			return fmt.Errorf("failed to ensure Endpoints %s/%s: %w", namespace, name, err)
		}
	}

	return nil
}

// EndpointSliceReconciler defines an interface to create/update EndpointSlices.
type EndpointSliceReconciler = func(existing *discoveryv1.EndpointSlice) (*discoveryv1.EndpointSlice, error)

// NamedEndpointSliceReconcilerFactory returns the name of the resource and the corresponding Reconciler function.
type NamedEndpointSliceReconcilerFactory = func() (name string, reconciler EndpointSliceReconciler)

// EndpointSliceObjectWrapper adds a wrapper so the EndpointSliceReconciler matches ObjectReconciler.
// This is needed as Go does not support function interface matching.
func EndpointSliceObjectWrapper(reconciler EndpointSliceReconciler) ObjectReconciler {
	return func(existing ctrlruntimeclient.Object) (ctrlruntimeclient.Object, error) {
		if existing != nil {
			return reconciler(existing.(*discoveryv1.EndpointSlice))
		}
		return reconciler(&discoveryv1.EndpointSlice{})
	}
}

// ReconcileEndpointSlices will create and update the EndpointSlices coming from the passed EndpointSliceReconciler slice.
func ReconcileEndpointSlices(ctx context.Context, namedFactories []NamedEndpointSliceReconcilerFactory, namespace string, client ctrlruntimeclient.Client, objectModifiers ...ObjectModifier) error {
	for _, factory := range namedFactories {
		name, reconciler := factory()
		reconcileObject := EndpointSliceObjectWrapper(reconciler)
		reconcileObject = CreateWithNamespace(reconcileObject, namespace)
		reconcileObject = CreateWithName(reconcileObject, name)

		for _, objectModifier := range objectModifiers {
			reconcileObject = objectModifier(reconcileObject)
		}

		if err := EnsureNamedObject(ctx, types.NamespacedName{Namespace: namespace, Name: name}, reconcileObject, client, &discoveryv1.EndpointSlice{}, false); err != nil {
			return fmt.Errorf("failed to ensure EndpointSlice %s/%s: %w", namespace, name, err)
		}
	}

	return nil
}

// JobReconciler defines an interface to create/update Jobs.
type JobReconciler = func(existing *batchv1.Job) (*batchv1.Job, error)

// NamedJobReconcilerFactory returns the name of the resource and the corresponding Reconciler function.
type NamedJobReconcilerFactory = func() (name string, reconciler JobReconciler)

// JobObjectWrapper adds a wrapper so the JobReconciler matches ObjectReconciler.
// This is needed as Go does not support function interface matching.
func JobObjectWrapper(reconciler JobReconciler) ObjectReconciler {
	return func(existing ctrlruntimeclient.Object) (ctrlruntimeclient.Object, error) {
		if existing != nil {
			return reconciler(existing.(*batchv1.Job))
		}
		return reconciler(&batchv1.Job{})
	}
}

// ReconcileJobs will create and update the Jobs coming from the passed JobReconciler slice.
func ReconcileJobs(ctx context.Context, namedFactories []NamedJobReconcilerFactory, namespace string, client ctrlruntimeclient.Client, objectModifiers ...ObjectModifier) error {
	for _, factory := range namedFactories {
		name, reconciler := factory()
		reconcileObject := JobObjectWrapper(reconciler)
		reconcileObject = CreateWithNamespace(reconcileObject, namespace)
		reconcileObject = CreateWithName(reconcileObject, name)

		for _, objectModifier := range objectModifiers {
			reconcileObject = objectModifier(reconcileObject)
		}

		if err := EnsureNamedObject(ctx, types.NamespacedName{Namespace: namespace, Name: name}, reconcileObject, client, &batchv1.Job{}, false); err != nil {
			return fmt.Errorf("failed to ensure Job %s/%s: %w", namespace, name, err)
		}
	}

	return nil
}

// CronJobReconciler defines an interface to create/update CronJobs.
type CronJobReconciler = func(existing *batchv1.CronJob) (*batchv1.CronJob, error)

// NamedCronJobReconcilerFactory returns the name of the resource and the corresponding Reconciler function.
type NamedCronJobReconcilerFactory = func() (name string, reconciler CronJobReconciler)

// CronJobObjectWrapper adds a wrapper so the CronJobReconciler matches ObjectReconciler.
// This is needed as Go does not support function interface matching.
func CronJobObjectWrapper(reconciler CronJobReconciler) ObjectReconciler {
	return func(existing ctrlruntimeclient.Object) (ctrlruntimeclient.Object, error) {
		if existing != nil {
			return reconciler(existing.(*batchv1.CronJob))
		}
		return reconciler(&batchv1.CronJob{})
	}
}

// ReconcileCronJobs will create and update the CronJobs coming from the passed CronJobReconciler slice.
func ReconcileCronJobs(ctx context.Context, namedFactories []NamedCronJobReconcilerFactory, namespace string, client ctrlruntimeclient.Client, objectModifiers ...ObjectModifier) error {
	for _, factory := range namedFactories {
		name, reconciler := factory()
		reconciler = DefaultCronJob(reconciler)
		reconcileObject := CronJobObjectWrapper(reconciler)
		reconcileObject = CreateWithNamespace(reconcileObject, namespace)
		reconcileObject = CreateWithName(reconcileObject, name)

		for _, objectModifier := range objectModifiers {
			reconcileObject = objectModifier(reconcileObject)
		}

		if err := EnsureNamedObject(ctx, types.NamespacedName{Namespace: namespace, Name: name}, reconcileObject, client, &batchv1.CronJob{}, false); err != nil {
			return fmt.Errorf("failed to ensure CronJob %s/%s: %w", namespace, name, err)
		}
	}

	return nil
}

// ValidatingWebhookConfigurationReconciler defines an interface to create/update ValidatingWebhookConfigurations.
type ValidatingWebhookConfigurationReconciler = func(existing *admissionregistrationv1.ValidatingWebhookConfiguration) (*admissionregistrationv1.ValidatingWebhookConfiguration, error)

// NamedValidatingWebhookConfigurationReconcilerFactory returns the name of the resource and the corresponding Reconciler function.
type NamedValidatingWebhookConfigurationReconcilerFactory = func() (name string, reconciler ValidatingWebhookConfigurationReconciler)

// ValidatingWebhookConfigurationObjectWrapper adds a wrapper so the ValidatingWebhookConfigurationReconciler matches ObjectReconciler.
// This is needed as Go does not support function interface matching.
func ValidatingWebhookConfigurationObjectWrapper(reconciler ValidatingWebhookConfigurationReconciler) ObjectReconciler {
	return func(existing ctrlruntimeclient.Object) (ctrlruntimeclient.Object, error) {
		if existing != nil {
			return reconciler(existing.(*admissionregistrationv1.ValidatingWebhookConfiguration))
		}
		return reconciler(&admissionregistrationv1.ValidatingWebhookConfiguration{})
	}
}

// ReconcileValidatingWebhookConfigurations will create and update the ValidatingWebhookConfigurations coming from the passed ValidatingWebhookConfigurationReconciler slice.
func ReconcileValidatingWebhookConfigurations(ctx context.Context, namedFactories []NamedValidatingWebhookConfigurationReconcilerFactory, namespace string, client ctrlruntimeclient.Client, objectModifiers ...ObjectModifier) error {
	for _, factory := range namedFactories {
		name, reconciler := factory()
		reconcileObject := ValidatingWebhookConfigurationObjectWrapper(reconciler)
		reconcileObject = CreateWithNamespace(reconcileObject, namespace)
		reconcileObject = CreateWithName(reconcileObject, name)

		for _, objectModifier := range objectModifiers {
			reconcileObject = objectModifier(reconcileObject)
		}

		if err := EnsureNamedObject(ctx, types.NamespacedName{Namespace: namespace, Name: name}, reconcileObject, client, &admissionregistrationv1.ValidatingWebhookConfiguration{}, false); err != nil {
			return fmt.Errorf("failed to ensure ValidatingWebhookConfiguration %s/%s: %w", namespace, name, err)
		}
	}

	return nil
}

// MutatingWebhookConfigurationReconciler defines an interface to create/update MutatingWebhookConfigurations.
type MutatingWebhookConfigurationReconciler = func(existing *admissionregistrationv1.MutatingWebhookConfiguration) (*admissionregistrationv1.MutatingWebhookConfiguration, error)

// NamedMutatingWebhookConfigurationReconcilerFactory returns the name of the resource and the corresponding Reconciler function.
type NamedMutatingWebhookConfigurationReconcilerFactory = func() (name string, reconciler MutatingWebhookConfigurationReconciler)

// MutatingWebhookConfigurationObjectWrapper adds a wrapper so the MutatingWebhookConfigurationReconciler matches ObjectReconciler.
// This is needed as Go does not support function interface matching.
func MutatingWebhookConfigurationObjectWrapper(reconciler MutatingWebhookConfigurationReconciler) ObjectReconciler {
	return func(existing ctrlruntimeclient.Object) (ctrlruntimeclient.Object, error) {
		if existing != nil {
			return reconciler(existing.(*admissionregistrationv1.MutatingWebhookConfiguration))
		}
		return reconciler(&admissionregistrationv1.MutatingWebhookConfiguration{})
	}
}

// ReconcileMutatingWebhookConfigurations will create and update the MutatingWebhookConfigurations coming from the passed MutatingWebhookConfigurationReconciler slice.
func ReconcileMutatingWebhookConfigurations(ctx context.Context, namedFactories []NamedMutatingWebhookConfigurationReconcilerFactory, namespace string, client ctrlruntimeclient.Client, objectModifiers ...ObjectModifier) error {
	for _, factory := range namedFactories {
		name, reconciler := factory()
		reconcileObject := MutatingWebhookConfigurationObjectWrapper(reconciler)
		reconcileObject = CreateWithNamespace(reconcileObject, namespace)
		reconcileObject = CreateWithName(reconcileObject, name)

		for _, objectModifier := range objectModifiers {
			reconcileObject = objectModifier(reconcileObject)
		}

		if err := EnsureNamedObject(ctx, types.NamespacedName{Namespace: namespace, Name: name}, reconcileObject, client, &admissionregistrationv1.MutatingWebhookConfiguration{}, false); err != nil {
			return fmt.Errorf("failed to ensure MutatingWebhookConfiguration %s/%s: %w", namespace, name, err)
		}
	}

	return nil
}

// StatefulSetReconciler defines an interface to create/update StatefulSets.
type StatefulSetReconciler = func(existing *appsv1.StatefulSet) (*appsv1.StatefulSet, error)

// NamedStatefulSetReconcilerFactory returns the name of the resource and the corresponding Reconciler function.
type NamedStatefulSetReconcilerFactory = func() (name string, reconciler StatefulSetReconciler)

// StatefulSetObjectWrapper adds a wrapper so the StatefulSetReconciler matches ObjectReconciler.
// This is needed as Go does not support function interface matching.
func StatefulSetObjectWrapper(reconciler StatefulSetReconciler) ObjectReconciler {
	return func(existing ctrlruntimeclient.Object) (ctrlruntimeclient.Object, error) {
		if existing != nil {
			return reconciler(existing.(*appsv1.StatefulSet))
		}
		return reconciler(&appsv1.StatefulSet{})
	}
}

// ReconcileStatefulSets will create and update the StatefulSets coming from the passed StatefulSetReconciler slice.
func ReconcileStatefulSets(ctx context.Context, namedFactories []NamedStatefulSetReconcilerFactory, namespace string, client ctrlruntimeclient.Client, objectModifiers ...ObjectModifier) error {
	for _, factory := range namedFactories {
		name, reconciler := factory()
		reconciler = DefaultStatefulSet(reconciler)
		reconcileObject := StatefulSetObjectWrapper(reconciler)
		reconcileObject = CreateWithNamespace(reconcileObject, namespace)
		reconcileObject = CreateWithName(reconcileObject, name)

		for _, objectModifier := range objectModifiers {
			reconcileObject = objectModifier(reconcileObject)
		}

		if err := EnsureNamedObject(ctx, types.NamespacedName{Namespace: namespace, Name: name}, reconcileObject, client, &appsv1.StatefulSet{}, false); err != nil {
			return fmt.Errorf("failed to ensure StatefulSet %s/%s: %w", namespace, name, err)
		}
	}

	return nil
}

// DeploymentReconciler defines an interface to create/update Deployments.
type DeploymentReconciler = func(existing *appsv1.Deployment) (*appsv1.Deployment, error)

// NamedDeploymentReconcilerFactory returns the name of the resource and the corresponding Reconciler function.
type NamedDeploymentReconcilerFactory = func() (name string, reconciler DeploymentReconciler)

// DeploymentObjectWrapper adds a wrapper so the DeploymentReconciler matches ObjectReconciler.
// This is needed as Go does not support function interface matching.
func DeploymentObjectWrapper(reconciler DeploymentReconciler) ObjectReconciler {
	return func(existing ctrlruntimeclient.Object) (ctrlruntimeclient.Object, error) {
		if existing != nil {
			return reconciler(existing.(*appsv1.Deployment))
		}
		return reconciler(&appsv1.Deployment{})
	}
}

// ReconcileDeployments will create and update the Deployments coming from the passed DeploymentReconciler slice.
func ReconcileDeployments(ctx context.Context, namedFactories []NamedDeploymentReconcilerFactory, namespace string, client ctrlruntimeclient.Client, objectModifiers ...ObjectModifier) error {
	for _, factory := range namedFactories {
		name, reconciler := factory()
		reconciler = DefaultDeployment(reconciler)
		reconcileObject := DeploymentObjectWrapper(reconciler)
		reconcileObject = CreateWithNamespace(reconcileObject, namespace)
		reconcileObject = CreateWithName(reconcileObject, name)

		for _, objectModifier := range objectModifiers {
			reconcileObject = objectModifier(reconcileObject)
		}

		if err := EnsureNamedObject(ctx, types.NamespacedName{Namespace: namespace, Name: name}, reconcileObject, client, &appsv1.Deployment{}, false); err != nil {
			return fmt.Errorf("failed to ensure Deployment %s/%s: %w", namespace, name, err)
		}
	}

	return nil
}

// DaemonSetReconciler defines an interface to create/update DaemonSets.
type DaemonSetReconciler = func(existing *appsv1.DaemonSet) (*appsv1.DaemonSet, error)

// NamedDaemonSetReconcilerFactory returns the name of the resource and the corresponding Reconciler function.
type NamedDaemonSetReconcilerFactory = func() (name string, reconciler DaemonSetReconciler)

// DaemonSetObjectWrapper adds a wrapper so the DaemonSetReconciler matches ObjectReconciler.
// This is needed as Go does not support function interface matching.
func DaemonSetObjectWrapper(reconciler DaemonSetReconciler) ObjectReconciler {
	return func(existing ctrlruntimeclient.Object) (ctrlruntimeclient.Object, error) {
		if existing != nil {
			return reconciler(existing.(*appsv1.DaemonSet))
		}
		return reconciler(&appsv1.DaemonSet{})
	}
}

// ReconcileDaemonSets will create and update the DaemonSets coming from the passed DaemonSetReconciler slice.
func ReconcileDaemonSets(ctx context.Context, namedFactories []NamedDaemonSetReconcilerFactory, namespace string, client ctrlruntimeclient.Client, objectModifiers ...ObjectModifier) error {
	for _, factory := range namedFactories {
		name, reconciler := factory()
		reconciler = DefaultDaemonSet(reconciler)
		reconcileObject := DaemonSetObjectWrapper(reconciler)
		reconcileObject = CreateWithNamespace(reconcileObject, namespace)
		reconcileObject = CreateWithName(reconcileObject, name)

		for _, objectModifier := range objectModifiers {
			reconcileObject = objectModifier(reconcileObject)
		}

		if err := EnsureNamedObject(ctx, types.NamespacedName{Namespace: namespace, Name: name}, reconcileObject, client, &appsv1.DaemonSet{}, false); err != nil {
			return fmt.Errorf("failed to ensure DaemonSet %s/%s: %w", namespace, name, err)
		}
	}

	return nil
}

// RoleReconciler defines an interface to create/update Roles.
type RoleReconciler = func(existing *rbacv1.Role) (*rbacv1.Role, error)

// NamedRoleReconcilerFactory returns the name of the resource and the corresponding Reconciler function.
type NamedRoleReconcilerFactory = func() (name string, reconciler RoleReconciler)

// RoleObjectWrapper adds a wrapper so the RoleReconciler matches ObjectReconciler.
// This is needed as Go does not support function interface matching.
func RoleObjectWrapper(reconciler RoleReconciler) ObjectReconciler {
	return func(existing ctrlruntimeclient.Object) (ctrlruntimeclient.Object, error) {
		if existing != nil {
			return reconciler(existing.(*rbacv1.Role))
		}
		return reconciler(&rbacv1.Role{})
	}
}

// ReconcileRoles will create and update the Roles coming from the passed RoleReconciler slice.
func ReconcileRoles(ctx context.Context, namedFactories []NamedRoleReconcilerFactory, namespace string, client ctrlruntimeclient.Client, objectModifiers ...ObjectModifier) error {
	for _, factory := range namedFactories {
		name, reconciler := factory()
		reconcileObject := RoleObjectWrapper(reconciler)
		reconcileObject = CreateWithNamespace(reconcileObject, namespace)
		reconcileObject = CreateWithName(reconcileObject, name)

		for _, objectModifier := range objectModifiers {
			reconcileObject = objectModifier(reconcileObject)
		}

		if err := EnsureNamedObject(ctx, types.NamespacedName{Namespace: namespace, Name: name}, reconcileObject, client, &rbacv1.Role{}, false); err != nil {
			return fmt.Errorf("failed to ensure Role %s/%s: %w", namespace, name, err)
		}
	}

	return nil
}

// RoleBindingReconciler defines an interface to create/update RoleBindings.
type RoleBindingReconciler = func(existing *rbacv1.RoleBinding) (*rbacv1.RoleBinding, error)

// NamedRoleBindingReconcilerFactory returns the name of the resource and the corresponding Reconciler function.
type NamedRoleBindingReconcilerFactory = func() (name string, reconciler RoleBindingReconciler)

// RoleBindingObjectWrapper adds a wrapper so the RoleBindingReconciler matches ObjectReconciler.
// This is needed as Go does not support function interface matching.
func RoleBindingObjectWrapper(reconciler RoleBindingReconciler) ObjectReconciler {
	return func(existing ctrlruntimeclient.Object) (ctrlruntimeclient.Object, error) {
		if existing != nil {
			return reconciler(existing.(*rbacv1.RoleBinding))
		}
		return reconciler(&rbacv1.RoleBinding{})
	}
}

// ReconcileRoleBindings will create and update the RoleBindings coming from the passed RoleBindingReconciler slice.
func ReconcileRoleBindings(ctx context.Context, namedFactories []NamedRoleBindingReconcilerFactory, namespace string, client ctrlruntimeclient.Client, objectModifiers ...ObjectModifier) error {
	for _, factory := range namedFactories {
		name, reconciler := factory()
		reconcileObject := RoleBindingObjectWrapper(reconciler)
		reconcileObject = CreateWithNamespace(reconcileObject, namespace)
		reconcileObject = CreateWithName(reconcileObject, name)

		for _, objectModifier := range objectModifiers {
			reconcileObject = objectModifier(reconcileObject)
		}

		if err := EnsureNamedObject(ctx, types.NamespacedName{Namespace: namespace, Name: name}, reconcileObject, client, &rbacv1.RoleBinding{}, false); err != nil {
			return fmt.Errorf("failed to ensure RoleBinding %s/%s: %w", namespace, name, err)
		}
	}

	return nil
}

// ClusterRoleReconciler defines an interface to create/update ClusterRoles.
type ClusterRoleReconciler = func(existing *rbacv1.ClusterRole) (*rbacv1.ClusterRole, error)

// NamedClusterRoleReconcilerFactory returns the name of the resource and the corresponding Reconciler function.
type NamedClusterRoleReconcilerFactory = func() (name string, reconciler ClusterRoleReconciler)

// ClusterRoleObjectWrapper adds a wrapper so the ClusterRoleReconciler matches ObjectReconciler.
// This is needed as Go does not support function interface matching.
func ClusterRoleObjectWrapper(reconciler ClusterRoleReconciler) ObjectReconciler {
	return func(existing ctrlruntimeclient.Object) (ctrlruntimeclient.Object, error) {
		if existing != nil {
			return reconciler(existing.(*rbacv1.ClusterRole))
		}
		return reconciler(&rbacv1.ClusterRole{})
	}
}

// ReconcileClusterRoles will create and update the ClusterRoles coming from the passed ClusterRoleReconciler slice.
func ReconcileClusterRoles(ctx context.Context, namedFactories []NamedClusterRoleReconcilerFactory, namespace string, client ctrlruntimeclient.Client, objectModifiers ...ObjectModifier) error {
	for _, factory := range namedFactories {
		name, reconciler := factory()
		reconcileObject := ClusterRoleObjectWrapper(reconciler)
		reconcileObject = CreateWithNamespace(reconcileObject, namespace)
		reconcileObject = CreateWithName(reconcileObject, name)

		for _, objectModifier := range objectModifiers {
			reconcileObject = objectModifier(reconcileObject)
		}

		if err := EnsureNamedObject(ctx, types.NamespacedName{Namespace: namespace, Name: name}, reconcileObject, client, &rbacv1.ClusterRole{}, false); err != nil {
			return fmt.Errorf("failed to ensure ClusterRole %s/%s: %w", namespace, name, err)
		}
	}

	return nil
}

// ClusterRoleBindingReconciler defines an interface to create/update ClusterRoleBindings.
type ClusterRoleBindingReconciler = func(existing *rbacv1.ClusterRoleBinding) (*rbacv1.ClusterRoleBinding, error)

// NamedClusterRoleBindingReconcilerFactory returns the name of the resource and the corresponding Reconciler function.
type NamedClusterRoleBindingReconcilerFactory = func() (name string, reconciler ClusterRoleBindingReconciler)

// ClusterRoleBindingObjectWrapper adds a wrapper so the ClusterRoleBindingReconciler matches ObjectReconciler.
// This is needed as Go does not support function interface matching.
func ClusterRoleBindingObjectWrapper(reconciler ClusterRoleBindingReconciler) ObjectReconciler {
	return func(existing ctrlruntimeclient.Object) (ctrlruntimeclient.Object, error) {
		if existing != nil {
			return reconciler(existing.(*rbacv1.ClusterRoleBinding))
		}
		return reconciler(&rbacv1.ClusterRoleBinding{})
	}
}

// ReconcileClusterRoleBindings will create and update the ClusterRoleBindings coming from the passed ClusterRoleBindingReconciler slice.
func ReconcileClusterRoleBindings(ctx context.Context, namedFactories []NamedClusterRoleBindingReconcilerFactory, namespace string, client ctrlruntimeclient.Client, objectModifiers ...ObjectModifier) error {
	for _, factory := range namedFactories {
		name, reconciler := factory()
		reconcileObject := ClusterRoleBindingObjectWrapper(reconciler)
		reconcileObject = CreateWithNamespace(reconcileObject, namespace)
		reconcileObject = CreateWithName(reconcileObject, name)

		for _, objectModifier := range objectModifiers {
			reconcileObject = objectModifier(reconcileObject)
		}

		if err := EnsureNamedObject(ctx, types.NamespacedName{Namespace: namespace, Name: name}, reconcileObject, client, &rbacv1.ClusterRoleBinding{}, false); err != nil {
			return fmt.Errorf("failed to ensure ClusterRoleBinding %s/%s: %w", namespace, name, err)
		}
	}

	return nil
}

// IngressReconciler defines an interface to create/update Ingresses.
type IngressReconciler = func(existing *networkingv1.Ingress) (*networkingv1.Ingress, error)

// NamedIngressReconcilerFactory returns the name of the resource and the corresponding Reconciler function.
type NamedIngressReconcilerFactory = func() (name string, reconciler IngressReconciler)

// IngressObjectWrapper adds a wrapper so the IngressReconciler matches ObjectReconciler.
// This is needed as Go does not support function interface matching.
func IngressObjectWrapper(reconciler IngressReconciler) ObjectReconciler {
	return func(existing ctrlruntimeclient.Object) (ctrlruntimeclient.Object, error) {
		if existing != nil {
			return reconciler(existing.(*networkingv1.Ingress))
		}
		return reconciler(&networkingv1.Ingress{})
	}
}

// ReconcileIngresses will create and update the Ingresses coming from the passed IngressReconciler slice.
func ReconcileIngresses(ctx context.Context, namedFactories []NamedIngressReconcilerFactory, namespace string, client ctrlruntimeclient.Client, objectModifiers ...ObjectModifier) error {
	for _, factory := range namedFactories {
		name, reconciler := factory()
		reconcileObject := IngressObjectWrapper(reconciler)
		reconcileObject = CreateWithNamespace(reconcileObject, namespace)
		reconcileObject = CreateWithName(reconcileObject, name)

		for _, objectModifier := range objectModifiers {
			reconcileObject = objectModifier(reconcileObject)
		}

		if err := EnsureNamedObject(ctx, types.NamespacedName{Namespace: namespace, Name: name}, reconcileObject, client, &networkingv1.Ingress{}, false); err != nil {
			return fmt.Errorf("failed to ensure Ingress %s/%s: %w", namespace, name, err)
		}
	}

	return nil
}

// NetworkPolicyReconciler defines an interface to create/update NetworkPolicies.
type NetworkPolicyReconciler = func(existing *networkingv1.NetworkPolicy) (*networkingv1.NetworkPolicy, error)

// NamedNetworkPolicyReconcilerFactory returns the name of the resource and the corresponding Reconciler function.
type NamedNetworkPolicyReconcilerFactory = func() (name string, reconciler NetworkPolicyReconciler)

// NetworkPolicyObjectWrapper adds a wrapper so the NetworkPolicyReconciler matches ObjectReconciler.
// This is needed as Go does not support function interface matching.
func NetworkPolicyObjectWrapper(reconciler NetworkPolicyReconciler) ObjectReconciler {
	return func(existing ctrlruntimeclient.Object) (ctrlruntimeclient.Object, error) {
		if existing != nil {
			return reconciler(existing.(*networkingv1.NetworkPolicy))
		}
		return reconciler(&networkingv1.NetworkPolicy{})
	}
}

// ReconcileNetworkPolicies will create and update the NetworkPolicies coming from the passed NetworkPolicyReconciler slice.
func ReconcileNetworkPolicies(ctx context.Context, namedFactories []NamedNetworkPolicyReconcilerFactory, namespace string, client ctrlruntimeclient.Client, objectModifiers ...ObjectModifier) error {
	for _, factory := range namedFactories {
		name, reconciler := factory()
		reconcileObject := NetworkPolicyObjectWrapper(reconciler)
		reconcileObject = CreateWithNamespace(reconcileObject, namespace)
		reconcileObject = CreateWithName(reconcileObject, name)

		for _, objectModifier := range objectModifiers {
			reconcileObject = objectModifier(reconcileObject)
		}

		if err := EnsureNamedObject(ctx, types.NamespacedName{Namespace: namespace, Name: name}, reconcileObject, client, &networkingv1.NetworkPolicy{}, false); err != nil {
			return fmt.Errorf("failed to ensure NetworkPolicy %s/%s: %w", namespace, name, err)
		}
	}

	return nil
}

// PodDisruptionBudgetReconciler defines an interface to create/update PodDisruptionBudgets.
type PodDisruptionBudgetReconciler = func(existing *policyv1.PodDisruptionBudget) (*policyv1.PodDisruptionBudget, error)

// NamedPodDisruptionBudgetReconcilerFactory returns the name of the resource and the corresponding Reconciler function.
type NamedPodDisruptionBudgetReconcilerFactory = func() (name string, reconciler PodDisruptionBudgetReconciler)

// PodDisruptionBudgetObjectWrapper adds a wrapper so the PodDisruptionBudgetReconciler matches ObjectReconciler.
// This is needed as Go does not support function interface matching.
func PodDisruptionBudgetObjectWrapper(reconciler PodDisruptionBudgetReconciler) ObjectReconciler {
	return func(existing ctrlruntimeclient.Object) (ctrlruntimeclient.Object, error) {
		if existing != nil {
			return reconciler(existing.(*policyv1.PodDisruptionBudget))
		}
		return reconciler(&policyv1.PodDisruptionBudget{})
	}
}

// ReconcilePodDisruptionBudgets will create and update the PodDisruptionBudgets coming from the passed PodDisruptionBudgetReconciler slice.
func ReconcilePodDisruptionBudgets(ctx context.Context, namedFactories []NamedPodDisruptionBudgetReconcilerFactory, namespace string, client ctrlruntimeclient.Client, objectModifiers ...ObjectModifier) error {
	for _, factory := range namedFactories {
		name, reconciler := factory()
		reconcileObject := PodDisruptionBudgetObjectWrapper(reconciler)
		reconcileObject = CreateWithNamespace(reconcileObject, namespace)
		reconcileObject = CreateWithName(reconcileObject, name)

		for _, objectModifier := range objectModifiers {
			reconcileObject = objectModifier(reconcileObject)
		}

		if err := EnsureNamedObject(ctx, types.NamespacedName{Namespace: namespace, Name: name}, reconcileObject, client, &policyv1.PodDisruptionBudget{}, true); err != nil {
			return fmt.Errorf("failed to ensure PodDisruptionBudget %s/%s: %w", namespace, name, err)
		}
	}

	return nil
}
