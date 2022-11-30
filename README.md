# reconciler

This repository contains a simple library to make reconciling Kubernetes
resources a bit less tedious. It wraps the typical get-update-else-create
loop in a convenient interface and relies on code generation. A set of
default reconciler functions are shipped with this library.

## Usage: Library

For common resources like Deployments, ConfigMaps or Secrets you can use
the reconcilers that ship with this library directly.

```go
func myNamedConfigMapReconcilerFactory() (name string, reconciler reconciling.ConfigMapReconciler) {
   return "my-configmap", func myConfigMapReconciler(cm *corev1.ConfigMap) (*corev1.ConfigMap, error) {
      cm.Data["field"] = "new-value"

      return cm, nil
   }
}

reconcilers := []reconciling.NamedConfigMapReconcilerFactory{
   myNamedConfigMapReconcilerFactory(),
}

if err := reconciling.ReconcileConfigMaps(ctx, client, "namespace", reconcilers); err != nil {
   log.Fatalf("Failed to reconcile ConfigMap: %v", err)
}
```

## Usage: Code Generation

For non-standard APIs you can use `reconciler-gen` to generate your own reconciler code.
Create a YAML config file like so:

```yaml
package: reconciling
boilerplate: boilerplate.go.txt
resourceTypes:
  - { package: k8s.io/kube-aggregator/pkg/apis/apiregistration/v1, resourceName: APIService }
  - { package: kubevirt.io/api/instancetype/v1alpha1, resourceName: VirtualMachineInstancetype }
```

Then run `reconciler-gen`:

```bash
go run k8c.io/reconciler/cmd/reconciler-gen --config reconciling.yaml > zz_generated_reconcile.go
```
