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

// Package node defines the customized metric names related-to k8s-node objects,
// those metrics are stored (and also can be referred) by custom metrics
// api-server provided by katalyst.
package node

const (
	CustomMetricNodeCPUTotal      = "node_cpu_total"
	CustomMetricNodeCPUUsage      = "node_cpu_usage"
	CustomMetricNodeCPUUsageRatio = "node_cpu_usage_ratio"
	CustomMetricNodeCPULoad1Min   = "node_cpu_load_system_1min"
)

// real-time memory related metric
const (
	CustomMetricNodeMemoryFree      = "node_system_memory_free"
	CustomMetricNodeMemoryAvailable = "node_system_memory_available"
)

// real-time advisor-related metric
const (
	CustomMetricNodeAdvisorPoolLoad1Min = "node_advisor_pool_load_1min"
	CustomMetricNodeAdvisorKnobStatus   = "node_advisor_knob_status"
)

// real-time numa level memory bandwidth related metric
const (
	CustomMetricNUMAMemoryBandwidthTotal  = "numa_mbm_total"
	CustomMetricNUMAMemoryBandwidthLocal  = "numa_mbm_local"
	CustomMetricNUMAMemoryBandwidthVictim = "numa_mbm_victim"
)
