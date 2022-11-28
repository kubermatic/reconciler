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

package main

import (
	"fmt"
	"go/format"
	"log"
	"os"

	"github.com/spf13/pflag"

	"k8c.io/reconciler/pkg/codegen"
	"k8c.io/reconciler/pkg/config"

	"k8s.io/apimachinery/pkg/util/yaml"
)

func main() {
	opt := options{}
	opt.AddFlags()

	pflag.Parse()

	if err := opt.Validate(); err != nil {
		log.Fatalf("Invalid flags: %v", err)
	}

	cfg, err := loadConfiguration(opt.configFile)
	if err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	boilerplate, err := loadBoilerplate(cfg.Boilerplate)
	if err != nil {
		log.Fatalf("Failed to load boilerplate: %v", err)
	}

	rendered, err := codegen.Render(cfg)
	if err != nil {
		log.Fatalf("Failed to render code: %v", err)
	}

	formatted, err := format.Source(rendered)
	if err != nil {
		log.Fatalf("Failed to format generated source code: %v", err)
	}

	fmt.Println(boilerplate)
	fmt.Print(string(formatted))
}

func loadBoilerplate(filename string) (string, error) {
	if filename == "" {
		return "# This file was generated, DO NOT EDIT.", nil
	}

	contents, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return string(contents), nil
}

func loadConfiguration(filename string) (*config.Configuration, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := &config.Configuration{}
	if err = yaml.UnmarshalStrict(content, config); err != nil {
		return nil, fmt.Errorf("%s is not a valid Configuration structure: %w", filename, err)
	}

	return config, nil
}
