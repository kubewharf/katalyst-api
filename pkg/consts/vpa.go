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

// const variables for workload annotations about vpa.
const (
	// WorkloadAnnotationVPAEnabledKey disables for workload means that
	// we won't apply the recommended resources for pod belonging to this workload;
	// However, we may still do this calculation logic and update to status if vpa
	// CR is created for this workload
	WorkloadAnnotationVPAEnabledKey = "vpa.katalyst.kubewharf.io/enable"
	WorkloadAnnotationVPAEnabled    = "true"

	WorkloadAnnotationVPANameKey = "vpa.katalyst.kubewharf.io/name"

	// WorkloadAnnotationVPASelectorKey is pod label selector for non-native workload
	WorkloadAnnotationVPASelectorKey = "vpa.katalyst.kubewharf.io/selector"
)

const (
	VPAAnnotationVPARecNameKey = "vpa.katalyst.kubewharf.io/recName"

	VPAAnnotationWorkloadRetentionPolicyKey    = "vpa.katalyst.kubewharf.io/retentionPolicy"
	VPAAnnotationWorkloadRetentionPolicyRetain = "retain"
	VPAAnnotationWorkloadRetentionPolicyDelete = "delete"
)
