/*
Copyright 2022 The Katalyst Authors.

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

package consts

// const variables for pod annotations about vpa in-place resource update.
const (
	PodAnnotationInplaceUpdateResourcesKey = "pod.kubernetes.io/resizeResources"

	PodAnnotationInplaceUpdateResizePolicyKey     = "pod.kubernetes.io/resizePolicy"
	PodAnnotationInplaceUpdateResizePolicyRestart = "Restart"

	PodAnnotationInplaceUpdateResizingKey = "pod.kubernetes.io/inplace-update-resizing"
	PodAnnotationAggregatedRequestsKey    = "pod.kubernetes.io/pod-aggregated-requests"
)

// PodAnnotationNetClassKey is a const variable for pod annotation about net class.
const (
	PodAnnotationNetClassKey = "katalyst.kubewharf.io/net_class_id"
)

// PodAnnotationNUMABindResultKey is a const variable for pod annotation about numa bind result.
const (
	PodAnnotationNUMABindResultKey = "katalyst.kubewharf.io/numa_bind_result"
)

// PodAnnotationNICSelectionResultKey is a const variable for pod annotation about a nic selection result.
const (
	PodAnnotationNICSelectionResultKey = "katalyst.kubewharf.io/nic_selection_result"
)

const (
	// PodAnnotationResourcePoolKey is a const variable for pod annotation about resource pool name
	PodAnnotationResourcePoolKey = "katalyst.kubewharf.io/resource_pool"
)
