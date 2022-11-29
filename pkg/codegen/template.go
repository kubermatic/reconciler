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

package codegen

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"

	"k8c.io/reconciler/pkg/config"

	"k8s.io/apimachinery/pkg/util/sets"
)

var (
	//go:embed imports.go.template
	importsTemplate string

	//go:embed reconciler.go.template
	reconcilerTemplate string

	//go:embed file.go.template
	fileTemplate string

	extraTemplateFunctions = map[string]interface{}{
		"lowercaseFirst": lowercaseFirst,
	}
)

func lowercaseFirst(str string) string {
	return strings.ToLower(string(str[0])) + str[1:]
}

type importStatement struct {
	Alias   string
	Package string
}

type templateData struct {
	Configuration *config.Configuration
	Imports       []importStatement
}

func Render(cfg *config.Configuration) ([]byte, error) {
	imports, err := buildImports(cfg)
	if err != nil {
		return nil, err
	}

	for idx, rt := range cfg.ResourceTypes {
		if rt.ResourceNamePlural == "" {
			cfg.ResourceTypes[idx].ResourceNamePlural = rt.ResourceName + "s"
		}
	}

	tpl := template.
		New("").
		Funcs(sprig.TxtFuncMap()).
		Funcs(extraTemplateFunctions).
		Funcs(template.FuncMap{
			"isInternal": func() bool {
				return cfg.Internal
			},
		})

	if tpl, err = tpl.Parse(importsTemplate); err != nil {
		return nil, fmt.Errorf("failed to parse imports template: %w", err)
	}
	if tpl, err = tpl.Parse(reconcilerTemplate); err != nil {
		return nil, fmt.Errorf("failed to parse reconciler template: %w", err)
	}
	if tpl, err = tpl.Parse(fileTemplate); err != nil {
		return nil, fmt.Errorf("failed to parse file template: %w", err)
	}

	data := templateData{
		Configuration: cfg,
		Imports:       imports,
	}

	var b bytes.Buffer
	if err := tpl.Execute(&b, data); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func buildImports(cfg *config.Configuration) ([]importStatement, error) {
	var err error

	aliases := map[string]string{}
	imports := map[string]string{}

	for idx, rt := range cfg.ResourceTypes {
		if rt.Package == "" {
			return nil, fmt.Errorf("resource type %q has no package defined", rt.ResourceName)
		}

		alias := rt.ImportAlias
		if alias == "" {
			alias, err = deriveAlias(rt.Package)
			if err != nil {
				return nil, fmt.Errorf("cannot auto alias %q: %w", rt.Package, err)
			}

			rt.ImportAlias = alias
		}

		if existing, ok := aliases[alias]; ok && existing != rt.Package {
			return nil, fmt.Errorf("both %q and %q lead to the same auto-generated alias %q", existing, rt.Package, alias)
		}

		aliases[alias] = rt.Package
		imports[rt.Package] = alias
		cfg.ResourceTypes[idx] = rt
	}

	sortedImports := []importStatement{}

	for _, pkg := range sets.StringKeySet(imports).List() {
		sortedImports = append(sortedImports, importStatement{
			Alias:   imports[pkg],
			Package: pkg,
		})
	}

	return sortedImports, nil
}

func deriveAlias(packageName string) (string, error) {
	parts := strings.Split(packageName, "/")
	if len(parts) < 2 {
		return "", errors.New("package has less than 2 segments")
	}

	// remove dots, so that e.g. "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1"
	// results in "autoscalingk8siov1"
	pkg := strings.ReplaceAll(parts[len(parts)-2], ".", "")
	version := parts[len(parts)-1]

	return pkg + version, nil
}
