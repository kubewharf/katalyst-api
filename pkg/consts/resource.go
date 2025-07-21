// Copyright 2022 The Katalyst Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package consts

import v1 "k8s.io/api/core/v1"

// const variables for resource names of reclaimed resource
const (
	ReclaimedResourceMilliCPU v1.ResourceName = "resource.katalyst.kubewharf.io/reclaimed_millicpu"
	ReclaimedResourceMemory   v1.ResourceName = "resource.katalyst.kubewharf.io/reclaimed_memory"
)

// const variables for resource names of guaranteed resource
const (
	ResourceNetBandwidth    v1.ResourceName = "resource.katalyst.kubewharf.io/net_bandwidth"
	ResourceMemoryBandwidth v1.ResourceName = "resource.katalyst.kubewharf.io/memory_bandwidth"
)

// const variables for resource attributes of resources
const (
	// ResourceAnnotationKeyResourceIdentifier nominated the key to override the default name
	// field in pod-resource-server (for qrm-related protocols); if the name field can't be
	// guaranteed to be unique in some cases, we can relay on this annotation to get unique keys
	// (to replace with the default name)
	ResourceAnnotationKeyResourceIdentifier = "katalyst.kubewharf.io/resource_identifier"

	// ResourceAnnotationKeyNICNetNSName nominated the key indicating net namespace name of the NIC
	ResourceAnnotationKeyNICNetNSName = "katalyst.kubewharf.io/netns_name"
)
