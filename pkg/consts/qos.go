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

// const variables for pod annotations about qos level
const (
	PodAnnotationQoSLevelKey = "katalyst.kubewharf.io/qos_level"

	PodAnnotationQoSLevelReclaimedCores = "reclaimed_cores"
	PodAnnotationQoSLevelSharedCores    = "shared_cores"
	PodAnnotationQoSLevelDedicatedCores = "dedicated_cores"
	PodAnnotationQoSLevelSystemCores    = "system_cores"
)

// const variables for pod annotations about qos level enhancement in memory
const (
	PodAnnotationMemoryEnhancementKey = "katalyst.kubewharf.io/memory_enhancement"

	PodAnnotationMemoryEnhancementNumaBinding       = "numa_binding"
	PodAnnotationMemoryEnhancementNumaBindingEnable = "true"
)

// const variables for pod annotations about qos level enhancement in cpu
const (
	PodAnnotationCPUEnhancementKey = "katalyst.kubewharf.io/cpu_enhancement"

	PodAnnotationCPUEnhancementCPUSet = "cpuset_pool"

	PodAnnotationCPUEnhancementSuppressionToleranceRate = "suppression_tolerance_rate"
)
