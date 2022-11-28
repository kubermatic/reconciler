# Copyright 2022 The Kubermatic Kubernetes Platform contributors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

export CGO_ENABLED ?= 0
export GO111MODULE = on
export GOFLAGS ?= -mod=readonly -trimpath

CMD ?= $(filter-out OWNERS, $(notdir $(wildcard ./cmd/*)))
GOBUILDFLAGS ?= -v
GOOS ?= $(shell go env GOOS)
LDFLAGS += -extldflags '-static'
GOTOOLFLAGS ?= $(GOBUILDFLAGS) -ldflags '$(LDFLAGS_EXTRA) $(LDFLAGS)' $(GOTOOLFLAGS_EXTRA)
BUILD_DEST ?= _build

.PHONY: build
build: $(CMD)

.PHONY: $(CMD)
$(CMD): %: $(BUILD_DEST)/%

$(BUILD_DEST)/%: cmd/%
	GOOS=$(GOOS) go build $(GOTOOLFLAGS) -o $@ ./cmd/$*

.PHONY:  test
test:
	go test ./pkg/...

lint:
	golangci-lint run \
		--verbose \
		--print-resources-usage \
		./pkg/... ./cmd/...

.PHONY: verify-imports
verify-imports:
	./hack/verify-import-order.sh

.PHONY: check-dependencies
check-dependencies:
	go mod tidy
	go mod verify
	git diff --exit-code

.PHONY: spellcheck
spellcheck:
	./hack/verify-spelling.sh

.PHONY: clean
clean:
	rm -rf $(BUILD_DEST)
	@echo "Cleaned $(BUILD_DEST)"
