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

	PodAnnotationMemoryEnhancementRssOverUseThreshold = "rss_overuse_threshold"

	PodAnnotationMemoryEnhancementNumaBinding       = "numa_binding"
	PodAnnotationMemoryEnhancementNumaBindingEnable = "true"

	PodAnnotationMemoryEnhancementNumaExclusive       = "numa_exclusive"
	PodAnnotationMemoryEnhancementNumaExclusiveEnable = "true"

	// PodAnnotationMemoryEnhancementOOMPriority provides a mechanism to specify
	// the OOM priority for pods. Higher priority values indicate a higher likelihood
	// of surviving OOM events.
	//
	// For different QoS levels, the acceptable value ranges are as follows:
	// - reclaimed_cores: [-100, 0)
	// - shared_cores: [0, 100)
	// - dedicated_cores: [100, 200)
	// - system_cores: [200, 300)
	// Additionally, there are two predefined values for any pod:
	// - -300: Indicates that the OOM priority is ignored, and the pod does not
	//   participate in priority comparison.
	// - 300: Indicates that the OOM priority is set to the highest level, the pod
	//   will never be terminated due to OOM events from the perspective of OOM enhancement
	PodAnnotationMemoryEnhancementOOMPriority = "oom_priority"
)

// const variables for pod annotations about qos level enhancement in cpu
const (
	PodAnnotationCPUEnhancementKey = "katalyst.kubewharf.io/cpu_enhancement"

	PodAnnotationCPUEnhancementCPUSet = "cpuset_pool"

	PodAnnotationCPUEnhancementSuppressionToleranceRate = "suppression_tolerance_rate"
)

// const variables for pod annotations about qos level enhancement in network
const (
	PodAnnotationNetworkEnhancementKey = "katalyst.kubewharf.io/network_enhancement"

	PodAnnotationNetworkEnhancementNamespaceType              = "namespace_type"
	PodAnnotationNetworkEnhancementNamespaceTypeHost          = "host_ns"
	PodAnnotationNetworkEnhancementNamespaceTypeHostPrefer    = "host_ns_preferred"
	PodAnnotationNetworkEnhancementNamespaceTypeNotHost       = "anti_host_ns"
	PodAnnotationNetworkEnhancementNamespaceTypeNotHostPrefer = "anti_host_ns_preferred"

	PodAnnotationNetworkEnhancementAffinityRestricted     = "topology_affinity_restricted"
	PodAnnotationNetworkEnhancementAffinityRestrictedTrue = "true"
)
