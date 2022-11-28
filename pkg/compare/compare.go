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

package compare

import (
	"encoding/json"

	"k8c.io/reconciler/pkg/equality"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type DiffReporterFunc func(a, b metav1.Object)

var (
	// DiffReporter is called whenever the old and new objects are not
	// semantically equivalent; its main purpose is to aid in debugging
	// and to spot reconciling loops (often caused by server-side
	// defaulting).
	DiffReporter DiffReporterFunc
)

// DeepEqual compares both objects for equality.
func DeepEqual(a, b metav1.Object) bool {
	if equality.Semantic.DeepEqual(a, b) {
		return true
	}

	// For some reason unstructured objects returned from the api have types for their fields
	// that are not map[string]interface{} and don't even exist in our codebase like
	// `openshift.infrastructureStatus`, so we have to compare the wire format here.
	// We only do this for unstrucutred as this comparison is pretty expensive.
	if _, isUnstructured := a.(*unstructured.Unstructured); isUnstructured && jsonEqual(a, b) {
		return true
	}

	if DiffReporter != nil {
		DiffReporter(a, b)
	}

	return false
}

func jsonEqual(a, b interface{}) bool {
	aJSON, err := json.Marshal(a)
	if err != nil {
		return false
	}
	bJSON, err := json.Marshal(b)
	if err != nil {
		return false
	}

	return string(aJSON) == string(bJSON)
}
