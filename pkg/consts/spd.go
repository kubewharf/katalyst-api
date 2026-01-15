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

const (
	// PodAnnotationSPDNameKey is used to maintain corresponding spdName in pod
	// annotation to make metaServer to target its spd more conveniently.
	PodAnnotationSPDNameKey = "spd.katalyst.kubewharf.io/name"
)

// const variables for workload annotations about spd.
const (
	// WorkloadAnnotationSPDEnableKey provides a mechanism for white list when enabling spd,
	// if it's set as false, we should not maintain spd CR or calculate service profiling automatically.
	WorkloadAnnotationSPDEnableKey = "spd.katalyst.kubewharf.io/enable"
	WorkloadAnnotationSPDEnabled   = "true"
)

// const variables for spd.
const (
	// SPDAnnotationBaselineSentinelKey and SPDAnnotationExtendedBaselineSentinelKey is
	// updated by the SPD controller. It represents the sentinel pod among all pods managed
	// by this SPD. Agents or controllers can use this key to determine if a pod falls within
	// the baseline by comparing it with the pod's createTime and podName.
	SPDAnnotationBaselineSentinelKey         = "spd.katalyst.kubewharf.io/baselineSentinel"
	SPDAnnotationExtendedBaselineSentinelKey = "spd.katalyst.kubewharf.io/extendedBaselineSentinel"

	// SPDAnnotationKeyCustomCompareKey holds annotation for spd baseline compare key
	SPDAnnotationKeyCustomCompareKey = "spd.katalyst.kubewharf.io/customCompareKey"

	SPDBaselinePercentMax = 100
	SPDBaselinePercentMin = 0
)

// metric names for aggregate metric
const (
	// SPDAggMetricNameMemoryBandwidth is per core memory bandwidth
	SPDAggMetricNameMemoryBandwidth = "memory_bandwidth"
)
