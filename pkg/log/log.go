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

package log

import "go.uber.org/zap"

var (
	// logger is the logging instance used by the reconciling functions.
	logger *zap.SugaredLogger

	// nopLogger is used in case logger is nil. It's defaulted at runtime
	// instead of an init() function to not cause issues with other package's
	// init functions.
	nopLogger = zap.NewNop().Sugar()
)

func SetLogger(log *zap.SugaredLogger) {
	logger = log
}

func Logger() *zap.SugaredLogger {
	if logger == nil {
		return nopLogger
	}

	return logger
}
