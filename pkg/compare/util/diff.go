/*
Copyright 2020 The Kubermatic Kubernetes Platform contributors.

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

package util

import (
	"log"

	"github.com/go-test/deep"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func LogDiff(a, b metav1.Object) {
	maxDepth := deep.MaxDepth
	logErrors := deep.LogErrors

	// Kubernetes Objects can be deeper than the default 10 levels.
	deep.MaxDepth = 50
	deep.LogErrors = true

	// For informational purpose we use deep.equal as it tells us what the difference is.
	// We need to calculate the difference in both ways as deep.equal only does a one-way comparison
	diff := deep.Equal(a, b)
	if diff == nil {
		diff = deep.Equal(b, a)
	}

	deep.MaxDepth = maxDepth
	deep.LogErrors = logErrors

	log.Printf("Object differs from generated one:\n%s", diff)
}
