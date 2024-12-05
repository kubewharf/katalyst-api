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

type QoSLevel string

const (
	QoSLevelReclaimedCores QoSLevel = "reclaimed_cores"
	QoSLevelSharedCores    QoSLevel = "shared_cores"
	QoSLevelDedicatedCores QoSLevel = "dedicated_cores"
	QoSLevelSystemCores    QoSLevel = "system_cores"
)

// const variables for pod annotations about qos level
const (
	PodAnnotationQoSLevelKey = "katalyst.kubewharf.io/qos_level"

	PodAnnotationQoSLevelReclaimedCores = string(QoSLevelReclaimedCores)
	PodAnnotationQoSLevelSharedCores    = string(QoSLevelSharedCores)
	PodAnnotationQoSLevelDedicatedCores = string(QoSLevelDedicatedCores)
	PodAnnotationQoSLevelSystemCores    = string(QoSLevelSystemCores)
)

// const variables for pod annotations about qos level enhancement in memory
const (
	PodAnnotationMemoryEnhancementKey = "katalyst.kubewharf.io/memory_enhancement"

	// PodAnnotationMemoryEnhancementRssOverUseThreshold provides a mechanism to enable
	// the ability of overcommit for memory, and we will relay on this enhancement to ensure
	// memory protection if rss usage exceeds requests (based on this given ratio)
	PodAnnotationMemoryEnhancementRssOverUseThreshold = "rss_overuse_threshold"

	// PodAnnotationMemoryEnhancementNumaBinding provides a mechanism to enable numa-binding
	// for workload to provide more ultimate running performances.
	//
	// With PodAnnotationMemoryEnhancementNumaBinding but without PodAnnotationMemoryEnhancementNumaExclusive,
	// we have several constraints below:
	// 1. different workloads may still share the same numa
	//   - these workloads may still have contentions on memory bandwidth
	// 2. the request for pod can be settled in a single numa node
	//   - this to avoid complicated cross numa memory capacity/bandwidth control
	//
	// todo: this enhancement is only supported for dedicated-cores now,
	//  the community if to support shared-cores in the short future.
	PodAnnotationMemoryEnhancementNumaBinding       = "numa_binding"
	PodAnnotationMemoryEnhancementNumaBindingEnable = "true"

	// PodAnnotationMemoryEnhancementNumaExclusive provides a mechanism to enable numa-exclusive
	// for A SINGLE Pod to avoid contention on memory bandwidth and so on.
	//
	// - this enhancement is only supported for dedicated-cores, for now and foreseeable future
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

	// PodAnnotationCPUEnhancementCPUSet provides a mechanism separate cpuset into
	// several orthogonal pools to avoid cpu contentions for different types of workloads,
	// i.e. spark batch, flink streaming, web service may fall into three pools.
	// and, each individual pod should be put into only one pool.
	//
	// - this enhancement is only supported for shared-cores, for now and foreseeable future
	// - all pods will be settled in `default` pool if not specified
	PodAnnotationCPUEnhancementCPUSet = "cpuset_pool"

	// PodAnnotationCPUEnhancementSuppressionToleranceRate provides a mechanism to ensure
	// the quality for reclaimed resources. since reclaimed resources will always change
	// dynamically according to running states of none-reclaimed services, it may reach to
	// a point that the resource contention is still be tolerable for none-reclaimed services,
	// but the reclaimed services runs too slow and would rather be killed and rescheduled.
	// in this case, the workload can use this enhancement to trigger eviction.
	//
	// - this enhancement is only supported for shared-cores, for now and foreseeable future
	PodAnnotationCPUEnhancementSuppressionToleranceRate = "suppression_tolerance_rate"
)

// const variables for pod annotations about qos level enhancement in network
const (
	PodAnnotationNetworkEnhancementKey = "katalyst.kubewharf.io/network_enhancement"

	// PodAnnotationNetworkEnhancementNamespaceType provides a mechanism to select nic in different namespaces
	// - PodAnnotationNetworkEnhancementNamespaceTypeHost
	//   - only select nic device in host namespace
	//   - admit failed if not possible
	// - PodAnnotationNetworkEnhancementNamespaceTypeHostPrefer
	//   - prefer tp select nic device in non-host namespace
	//   - also accept nic device in non-host namespace if not possible
	// - PodAnnotationNetworkEnhancementNamespaceTypeNotHost
	//   - only select nic device in non-host namespace
	//   - admit failed if not possible
	// - PodAnnotationNetworkEnhancementNamespaceTypeNotHostPrefer
	//   - only select nic device in non-host namespace
	//	 - also accept nic device in host namespace if not possible
	PodAnnotationNetworkEnhancementNamespaceType              = "namespace_type"
	PodAnnotationNetworkEnhancementNamespaceTypeHost          = "host_ns"
	PodAnnotationNetworkEnhancementNamespaceTypeHostPrefer    = "host_ns_preferred"
	PodAnnotationNetworkEnhancementNamespaceTypeNotHost       = "anti_host_ns"
	PodAnnotationNetworkEnhancementNamespaceTypeNotHostPrefer = "anti_host_ns_preferred"

	// PodAnnotationNetworkEnhancementAffinityRestricted sets as true to indicate
	// we must ensure the numa affinity for nic devices, and we should admit failed if not possible
	PodAnnotationNetworkEnhancementAffinityRestricted     = "topology_affinity_restricted"
	PodAnnotationNetworkEnhancementAffinityRestrictedTrue = "true"
)

// ResourcePluginPolicyName is a string type for QosResourceManager plugin policy
type ResourcePluginPolicyName string

// const variables for QRM plugin policy name
const (
	// ResourcePluginPolicyNameDynamic is the name of the dynamic policy.
	ResourcePluginPolicyNameDynamic ResourcePluginPolicyName = "dynamic"
	// ResourcePluginPolicyNameNative is the name of the native policy.
	ResourcePluginPolicyNameNative ResourcePluginPolicyName = "native"
	// ResourcePluginPolicyNameStatic is the name of the static policy.
	ResourcePluginPolicyNameStatic ResourcePluginPolicyName = "static"
)
