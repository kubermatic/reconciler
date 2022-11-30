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
type NamespaceReconciler = GenericObjectReconciler[*corev1.Namespace]

// NamedNamespaceReconcilerFactory returns the name of the resource and the corresponding Reconciler function.
type NamedNamespaceReconcilerFactory = GenericNamedObjectReconciler[*corev1.Namespace]

// ServiceReconciler defines an interface to create/update Services.
type ServiceReconciler = GenericObjectReconciler[*corev1.Service]

// NamedServiceReconcilerFactory returns the name of the resource and the corresponding Reconciler function.
type NamedServiceReconcilerFactory = GenericNamedObjectReconciler[*corev1.Service]

// SecretReconciler defines an interface to create/update Secrets.
type SecretReconciler = GenericObjectReconciler[*corev1.Secret]

// NamedSecretReconcilerFactory returns the name of the resource and the corresponding Reconciler function.
type NamedSecretReconcilerFactory = GenericNamedObjectReconciler[*corev1.Secret]

// ConfigMapReconciler defines an interface to create/update ConfigMaps.
type ConfigMapReconciler = GenericObjectReconciler[*corev1.ConfigMap]

// NamedConfigMapReconcilerFactory returns the name of the resource and the corresponding Reconciler function.
type NamedConfigMapReconcilerFactory = GenericNamedObjectReconciler[*corev1.ConfigMap]

// ServiceAccountReconciler defines an interface to create/update ServiceAccounts.
type ServiceAccountReconciler = GenericObjectReconciler[*corev1.ServiceAccount]

// NamedServiceAccountReconcilerFactory returns the name of the resource and the corresponding Reconciler function.
type NamedServiceAccountReconcilerFactory = GenericNamedObjectReconciler[*corev1.ServiceAccount]

// EndpointsReconciler defines an interface to create/update Endpoints.
type EndpointsReconciler = GenericObjectReconciler[*corev1.Endpoints]

// NamedEndpointsReconcilerFactory returns the name of the resource and the corresponding Reconciler function.
type NamedEndpointsReconcilerFactory = GenericNamedObjectReconciler[*corev1.Endpoints]

// EndpointSliceReconciler defines an interface to create/update EndpointSlices.
type EndpointSliceReconciler = GenericObjectReconciler[*discoveryv1.EndpointSlice]

// NamedEndpointSliceReconcilerFactory returns the name of the resource and the corresponding Reconciler function.
type NamedEndpointSliceReconcilerFactory = GenericNamedObjectReconciler[*discoveryv1.EndpointSlice]

// JobReconciler defines an interface to create/update Jobs.
type JobReconciler = GenericObjectReconciler[*batchv1.Job]

// NamedJobReconcilerFactory returns the name of the resource and the corresponding Reconciler function.
type NamedJobReconcilerFactory = GenericNamedObjectReconciler[*batchv1.Job]

// CronJobReconciler defines an interface to create/update CronJobs.
type CronJobReconciler = GenericObjectReconciler[*batchv1.CronJob]

// NamedCronJobReconcilerFactory returns the name of the resource and the corresponding Reconciler function.
type NamedCronJobReconcilerFactory = GenericNamedObjectReconciler[*batchv1.CronJob]

// ValidatingWebhookConfigurationReconciler defines an interface to create/update ValidatingWebhookConfigurations.
type ValidatingWebhookConfigurationReconciler = GenericObjectReconciler[*admissionregistrationv1.ValidatingWebhookConfiguration]

// NamedValidatingWebhookConfigurationReconcilerFactory returns the name of the resource and the corresponding Reconciler function.
type NamedValidatingWebhookConfigurationReconcilerFactory = GenericNamedObjectReconciler[*admissionregistrationv1.ValidatingWebhookConfiguration]

// MutatingWebhookConfigurationReconciler defines an interface to create/update MutatingWebhookConfigurations.
type MutatingWebhookConfigurationReconciler = GenericObjectReconciler[*admissionregistrationv1.MutatingWebhookConfiguration]

// NamedMutatingWebhookConfigurationReconcilerFactory returns the name of the resource and the corresponding Reconciler function.
type NamedMutatingWebhookConfigurationReconcilerFactory = GenericNamedObjectReconciler[*admissionregistrationv1.MutatingWebhookConfiguration]

// StatefulSetReconciler defines an interface to create/update StatefulSets.
type StatefulSetReconciler = GenericObjectReconciler[*appsv1.StatefulSet]

// NamedStatefulSetReconcilerFactory returns the name of the resource and the corresponding Reconciler function.
type NamedStatefulSetReconcilerFactory = GenericNamedObjectReconciler[*appsv1.StatefulSet]

// DeploymentReconciler defines an interface to create/update Deployments.
type DeploymentReconciler = GenericObjectReconciler[*appsv1.Deployment]

// NamedDeploymentReconcilerFactory returns the name of the resource and the corresponding Reconciler function.
type NamedDeploymentReconcilerFactory = GenericNamedObjectReconciler[*appsv1.Deployment]

// DaemonSetReconciler defines an interface to create/update DaemonSets.
type DaemonSetReconciler = GenericObjectReconciler[*appsv1.DaemonSet]

// NamedDaemonSetReconcilerFactory returns the name of the resource and the corresponding Reconciler function.
type NamedDaemonSetReconcilerFactory = GenericNamedObjectReconciler[*appsv1.DaemonSet]

// RoleReconciler defines an interface to create/update Roles.
type RoleReconciler = GenericObjectReconciler[*rbacv1.Role]

// NamedRoleReconcilerFactory returns the name of the resource and the corresponding Reconciler function.
type NamedRoleReconcilerFactory = GenericNamedObjectReconciler[*rbacv1.Role]

// RoleBindingReconciler defines an interface to create/update RoleBindings.
type RoleBindingReconciler = GenericObjectReconciler[*rbacv1.RoleBinding]

// NamedRoleBindingReconcilerFactory returns the name of the resource and the corresponding Reconciler function.
type NamedRoleBindingReconcilerFactory = GenericNamedObjectReconciler[*rbacv1.RoleBinding]

// ClusterRoleReconciler defines an interface to create/update ClusterRoles.
type ClusterRoleReconciler = GenericObjectReconciler[*rbacv1.ClusterRole]

// NamedClusterRoleReconcilerFactory returns the name of the resource and the corresponding Reconciler function.
type NamedClusterRoleReconcilerFactory = GenericNamedObjectReconciler[*rbacv1.ClusterRole]

// ClusterRoleBindingReconciler defines an interface to create/update ClusterRoleBindings.
type ClusterRoleBindingReconciler = GenericObjectReconciler[*rbacv1.ClusterRoleBinding]

// NamedClusterRoleBindingReconcilerFactory returns the name of the resource and the corresponding Reconciler function.
type NamedClusterRoleBindingReconcilerFactory = GenericNamedObjectReconciler[*rbacv1.ClusterRoleBinding]

// IngressReconciler defines an interface to create/update Ingresses.
type IngressReconciler = GenericObjectReconciler[*networkingv1.Ingress]

// NamedIngressReconcilerFactory returns the name of the resource and the corresponding Reconciler function.
type NamedIngressReconcilerFactory = GenericNamedObjectReconciler[*networkingv1.Ingress]

// NetworkPolicyReconciler defines an interface to create/update NetworkPolicies.
type NetworkPolicyReconciler = GenericObjectReconciler[*networkingv1.NetworkPolicy]

// NamedNetworkPolicyReconcilerFactory returns the name of the resource and the corresponding Reconciler function.
type NamedNetworkPolicyReconcilerFactory = GenericNamedObjectReconciler[*networkingv1.NetworkPolicy]

// PodDisruptionBudgetReconciler defines an interface to create/update PodDisruptionBudgets.
type PodDisruptionBudgetReconciler = GenericObjectReconciler[*policyv1.PodDisruptionBudget]

// NamedPodDisruptionBudgetReconcilerFactory returns the name of the resource and the corresponding Reconciler function.
type NamedPodDisruptionBudgetReconcilerFactory = GenericNamedObjectReconciler[*policyv1.PodDisruptionBudget]
