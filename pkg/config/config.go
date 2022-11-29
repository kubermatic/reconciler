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

package config

type Configuration struct {
	Package       string         `yaml:"package"`
	ResourceTypes []ResourceType `yaml:"resourceTypes"`
	Boilerplate   string         `yaml:"boilerplate"`

	// Internal should only be set if the generated code is placed in the same
	// package as all the other helpers (i.e. practically never outside of this
	// repository).
	Internal bool `yaml:"internal"`
}

type ResourceType struct {
	// The uppercase, singular resource name like "Deployment".
	ResourceName string `yaml:"resourceName"`
	// The uppercase, plural resource name like "Deployments"; if left
	// empty is defaulted to "<ResourceName>s".
	ResourceNamePlural string `yaml:"resourceNamePlural"`
	// The Go package that providers the given ResourceName, e.g. "k8s.io/api/apps/v1".
	Package string `yaml:"package"`
	// Optional: An Go import alias; if none is provided then a standard
	// Kubernetes structure like "pkg/apis/foo/v1" is assumed in the PackageName
	// and the alias is automatically derived from joining the last two path
	// elements (in this case "foov1").
	ImportAlias string `yaml:"importAlias"`
	// Optional: A defaulting func for the given object type
	// Must be defined inside the target package.
	DefaultingFunc string `yaml:"defaultingFunc"`
	// Whether the resource must be recreated instead of updated. Required
	// e.G. for PDBs
	Recreate bool `yaml:"recreate"`
	// Optional: adds an api version prefix to the generated functions to avoid duplication when different resources
	// have the same ResourceName
	APIVersionPrefix string `yaml:"apiVersionPrefix"`
}
