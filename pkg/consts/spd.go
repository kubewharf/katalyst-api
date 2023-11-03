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
	PodAnnotationSPDNameKey = "spd.katalyst.kubewharf.io/name"
)

// const variables for workload annotations about spd.
const (
	// WorkloadAnnotationSPDEnableKey disables for workload means that we should not
	// maintain spd CR and much less to calculate service profiling automatically
	WorkloadAnnotationSPDEnableKey = "spd.katalyst.kubewharf.io/enable"
	WorkloadAnnotationSPDEnabled   = "true"

	WorkloadAnnotationSPDNameKey = "spd.katalyst.kubewharf.io/name"
)

// const variables for spd annotations.
const (
	// SPDAnnotationBaselinePercentileKey is updated by the SPD controller. It represents
	// the baseline percentile across all pods managed by this SPD. Agents or controllers
	// can use this key to determine if a pod falls within the baseline by comparing it
	// with the pod's baseline coefficient.
	SPDAnnotationBaselinePercentileKey = "spd.katalyst.kubewharf.io/baselinePercentile"
)
