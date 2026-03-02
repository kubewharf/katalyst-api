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

// Package pod defines the customized metric names related-to k8s-pod objects,
// those metrics are stored (and also can be referred) by custom metrics
// api-server provided by katalyst.
package pod

// real-time cpu related metric
const (
	CustomMetricPodCPULoad1Min   = "pod_cpu_load_1min"
	CustomMetricPodCPUUsage      = "pod_cpu_usage"
	CustomMetricPodCPUUsageRatio = "pod_cpu_usage_ratio"
	CustomMetricPodCPUCPI        = "pod_cpu_cpi"
)

const (
	CustomMetricPodMemoryRSS   = "pod_memory_rss"
	CustomMetricPodMemoryUsage = "pod_memory_usage"
)

const (
	CustomMetricPodGPUUsage = "pod_gpu_usage"
)

// real-time memory bandwidth related metric
const (
	CustomMetricPodTotalMemoryBandwidth  = "pod_mbm_total"
	CustomMetricPodLocalMemoryBandwidth  = "pod_mbm_local"
	CustomMetricPodVictimMemoryBandwidth = "pod_mbm_victim"
)
