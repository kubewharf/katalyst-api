# Copyright 2022 The Katalyst Authors.
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

GOPROXY := $(shell go env GOPROXY)
ifeq ($(GOPROXY),)
GOPROXY := https://proxy.golang.org
endif
export GOPROXY
GOPATH ?= $(shell go env GOPATH)

# Directories
TOOLS_DIR := hack/tools
BIN_DIR := bin
TOOLS_BIN_DIR := $(TOOLS_DIR)/$(BIN_DIR)

# Binaries
CONTROLLER_GEN := $(abspath $(TOOLS_BIN_DIR)/controller-gen)

# Protocols
Protocol_PATH = pkg/protocol
PROTO_INCLUDE_DIR := $(abspath .protoinclude)
GOGO_PROTOBUF_DIR := $(shell go list -m -f '{{.Dir}}' github.com/gogo/protobuf)
K8S_API_DIR := $(shell go list -m -f '{{.Dir}}' k8s.io/api)
K8S_APIMACHINERY_DIR := $(shell go list -m -f '{{.Dir}}' k8s.io/apimachinery)

all: generate
crd: generate-manifests generate-go
pb: generate-pb

## --------------------------------------
## Binaries
## --------------------------------------

$(CONTROLLER_GEN): $(TOOLS_DIR)/go.mod ## Build controller-gen from tools folder.
	cd $(TOOLS_DIR); \
		go build -tags=tools -o $(BIN_DIR)/controller-gen \
		sigs.k8s.io/controller-tools/cmd/controller-gen

## --------------------------------------
## Generate / Manifests
## --------------------------------------

.PHONY: generate
generate: $(CONTROLLER_GEN)
	$(MAKE) generate-manifests
	$(MAKE) generate-go
	$(MAKE) generate-pb

.PHONY: generate-go ## Generate go client codes
generate-go: hack/update-codegen.sh
	./hack/update-codegen.sh

.PHONY: generate-manifests ## Generate CRD manifests
generate-manifests: $(CONTROLLER_GEN)
	$(CONTROLLER_GEN) \
		paths=./pkg/apis/... \
		crd:crdVersions=v1,allowDangerousTypes=true \
		output:crd:dir=./config/crd/bases

## --------------------------------------
## Generate / Protocols
## --------------------------------------

.PHONY: generate-pb
generate-pb: prepare-proto-includes generate-reporter-plugin generate-eviction-plugin

.PHONY: prepare-proto-includes
prepare-proto-includes:
	mkdir -p "$(PROTO_INCLUDE_DIR)/github.com/gogo" "$(PROTO_INCLUDE_DIR)/k8s.io"
	ln -sfn "$(GOGO_PROTOBUF_DIR)" "$(PROTO_INCLUDE_DIR)/github.com/gogo/protobuf"
	ln -sfn "$(K8S_API_DIR)" "$(PROTO_INCLUDE_DIR)/k8s.io/api"
	ln -sfn "$(K8S_APIMACHINERY_DIR)" "$(PROTO_INCLUDE_DIR)/k8s.io/apimachinery"

.PHONY: generate-reporter-plugin ## Generate Protocol for reporter manager
generate-reporter-plugin:
	(protoc -I=$(Protocol_PATH)/reporterplugin/ -I=$(PROTO_INCLUDE_DIR) -I=$(GOPATH)/src/ --gogo_out=plugins=grpc:$(Protocol_PATH)/reporterplugin/ $(Protocol_PATH)/reporterplugin/v1alpha1/api.proto)
	cat ./hack/boilerplate.go.txt "$(Protocol_PATH)/reporterplugin/v1alpha1/api.pb.go" > tmpfile && mv tmpfile "$(Protocol_PATH)/reporterplugin/v1alpha1/api.pb.go"

.PHONY: generate-eviction-plugin ## Generate Protocol for eviction manager
generate-eviction-plugin:
	(protoc -I=$(Protocol_PATH)/evictionplugin/ -I=$(PROTO_INCLUDE_DIR) -I=$(GOPATH)/src/ --gogo_out=plugins=grpc:$(Protocol_PATH)/evictionplugin/ $(Protocol_PATH)/evictionplugin/v1alpha1/api.proto)
	cat ./hack/boilerplate.go.txt "$(Protocol_PATH)/evictionplugin/v1alpha1/api.pb.go" > tmpfile && mv tmpfile "$(Protocol_PATH)/evictionplugin/v1alpha1/api.pb.go"


## --------------------------------------
## Cleanup / Verification
## --------------------------------------

.PHONY: clean
clean: ## Remove all generated files
	$(MAKE) clean-bin

.PHONY: clean-bin
clean-bin: ## Remove all generated binaries
	rm -rf bin
	rm -rf hack/tools/bin

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...
