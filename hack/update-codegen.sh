#!/usr/bin/env bash

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

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

# generate the code with:
# --output-base    because this script should also be able to run inside the vendor dir of
#                  k8s.io/kubernetes. The output-base is needed for the generators to output into the vendor dir
#                  instead of the $GOPATH directly. For normal projects this can be dropped.
./hack/generate-groups.sh all \
   github.com/kubewharf/katalyst-api/pkg/client  github.com/kubewharf/katalyst-api/pkg/apis \
  "node:v1alpha1 autoscaling:v1alpha1 config:v1alpha1 workload:v1alpha1 overcommit:v1alpha1" \
  --output-base "$(dirname "${BASH_SOURCE[0]}")/../../../../" \
  --go-header-file "${SCRIPT_ROOT}"/hack/boilerplate.go.txt

./hack/generate-internal-groups.sh \
  "deepcopy,defaulter,conversion" \
  github.com/kubewharf/katalyst-api/pkg/apis/scheduling/generated \
  github.com/kubewharf/katalyst-api/pkg/apis/scheduling \
  github.com/kubewharf/katalyst-api/pkg/apis/scheduling \
  "config:v1beta3" \
  --go-header-file "${SCRIPT_ROOT}"/hack/boilerplate.go.txt
