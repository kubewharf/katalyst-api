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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kubewharf/katalyst-api/pkg/consts"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=userwatermarkconfigurations,shortName=uwm
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="AGE",type=date,JSONPath=.metadata.creationTimestamp
// +kubebuilder:printcolumn:name="PAUSED",type=boolean,JSONPath=".spec.paused"
// +kubebuilder:printcolumn:name="SELECTOR",type=string,JSONPath=".spec.nodeLabelSelector"
// +kubebuilder:printcolumn:name="PRIORITY",type=string,JSONPath=".spec.priority"
// +kubebuilder:printcolumn:name="NODES",type=string,JSONPath=".spec.ephemeralSelector.nodeNames"
// +kubebuilder:printcolumn:name="DURATION",type=string,JSONPath=".spec.ephemeralSelector.lastDuration"
// +kubebuilder:printcolumn:name="TARGET",type=integer,JSONPath=".status.targetNodes"
// +kubebuilder:printcolumn:name="CANARY",type=integer,JSONPath=".status.canaryNodes"
// +kubebuilder:printcolumn:name="UPDATED-TARGET",type=integer,JSONPath=".status.updatedTargetNodes"
// +kubebuilder:printcolumn:name="HASH",type=string,JSONPath=".status.currentHash"
// +kubebuilder:printcolumn:name="VALID",type=string,JSONPath=".status.conditions[?(@.type==\"Valid\")].status"
// +kubebuilder:printcolumn:name="REASON",type=string,JSONPath=".status.conditions[?(@.type==\"Valid\")].reason"
// +kubebuilder:printcolumn:name="MESSAGE",type=string,JSONPath=".status.conditions[?(@.type==\"Valid\")].message"

// UserWatermarkConfiguration is the Schema for the configuration API used by User Watermark
type UserWatermarkConfiguration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   UserWatermarkConfigurationSpec `json:"spec,omitempty"`
	Status GenericConfigStatus            `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

// UserWatermarkConfigurationList contains a list of UserWatermarkConfiguration
type UserWatermarkConfigurationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []UserWatermarkConfiguration `json:"items"`
}

// UserWatermarkConfigurationSpec defines the desired state of UserWatermarkConfiguration
type UserWatermarkConfigurationSpec struct {
	GenericConfigSpec `json:",inline"`

	// Config is custom field for UserWatermark configuration
	Config UserWatermarkConfig `json:"config"`
}

type UserWatermarkConfig struct {
	// EnableReclaimer is whether to enable reclaimer on target objective
	// +optional
	EnableReclaimer *bool `json:"enableReclaimer,omitempty"`
	// ReconcileInterval is the minimum duration the objectives got memory reclaimed by reclaimer
	// +optional
	ReconcileInterval *int64 `json:"reconcileInterval,omitempty"`
	// ServiceLabel is the label selector to filter services that need to be configured
	// +optional
	ServiceLabel string `json:"serviceLabel,omitempty"`

	// ServiceConfig is a configuration used to reclaim user watermark on specified services
	// +optional
	// +listMapKey=serviceName
	// +listType=map
	ServiceConfig []UserWatermarkServiceConfig `json:"serviceConfig,omitempty"`

	// QoSLevelConfig is a configuration used to reclaim user watermark on specified qos levels
	// +optional
	// +listMapKey=qosLevel
	// +listType=map
	QoSLevelConfig []UserWatermarkQoSLevelConfig `json:"qosLevelConfig,omitempty"`

	// CgroupConfig is a configuration used to reclaim user watermark on specified cgroups
	// +optional
	// +listMapKey=cgroupPath
	// +listType=map
	CgroupConfig []UserWatermarkCgroupConfig `json:"cgroupConfig,omitempty"`
}

type UserWatermarkServiceConfig struct {
	// ServiceName is the name of the service to be configured
	ServiceName string `json:"serviceName"`

	// ConfigDetail is configuration details of UserWatermark
	ConfigDetail ReclaimConfigDetail `json:"configDetail"`
}

type UserWatermarkQoSLevelConfig struct {
	// QoSLevel is either of reclaimed_cores, shared_cores, dedicated_cores, system_cores
	QoSLevel consts.QoSLevel `json:"qosLevel"`

	// ConfigDetail is configuration details of UserWatermark
	ConfigDetail ReclaimConfigDetail `json:"configDetail"`
}

type UserWatermarkCgroupConfig struct {
	// CgroupPath is an cgroupV2 absolute path, e.g. /sys/fs/cgroup/hdfs
	CgroupPath string `json:"cgroupPath"`

	// ConfigDetail is configuration details of UserWatermark
	ConfigDetail ReclaimConfigDetail `json:"configDetail"`
}

type ReclaimConfigDetail struct {
	// EnableMemoryReclaim is whether to enable memory reclaim on target objective
	// +optional
	EnableMemoryReclaim *bool `json:"enableMemoryReclaim,omitempty"`
	// ReclaimInterval is the minimum duration the objectives got memory reclaimed by reclaimer
	// +optional
	ReclaimInterval *int64 `json:"reclaimInterval,omitempty"`
	// ScaleFactor is the scale factor of memory reclaim size
	// +optional
	ScaleFactor *uint64 `json:"scaleFactor,omitempty"`
	// SingleReclaimFactor is the max memory reclaim size ratio in one reclaim cycle
	// +optional
	SingleReclaimFactor *float64 `json:"singleReclaimFactor,omitempty"`
	// SingleReclaimSize is the max memory reclaim size in one reclaim cycle
	// +optional
	SingleReclaimSize *uint64 `json:"singleReclaimSize,omitempty"`
	// BackoffDuration is the duration the reclaimer will wait before next reclaim cycle
	// +optional
	BackoffDuration *metav1.Duration `json:"backoffDuration,omitempty"`
	// FeedbackPolicy is the policy used to feedback memory reclaim status to reclaimer
	// +optional
	FeedbackPolicy *UserWatermarkPolicyName `json:"feedbackPolicy,omitempty"`
	// ReclaimFailedThreshold is the threshold of consecutive reclaim failed times. If the threshold
	// is reached, the reclaimer will pause reclaiming.
	// +optional
	ReclaimFailedThreshold *uint64 `json:"reclaimFailedThreshold,omitempty"`
	// FailureFreezePeriod is the duration the reclaimer will wait before next reclaim cycle after
	// consecutive reclaim failed times reach the threshold.
	// +optional
	FailureFreezePeriod *metav1.Duration `json:"failureFreezePeriod,omitempty"`

	// PSIPolicyConf is a configuration that uses the psi indicator to report memory reclamation status
	// +optional
	PSIPolicyConf *UserWatermarkPSIPolicyConf `json:"psiPolicy,omitempty"`

	// RefaultPolicy is a configuration that uses the refault indicator to report memory reclamation status
	// +optional
	RefaultPolicConf *UserWatermarkRefaultPolicyConf `json:"refaultPolicy,omitempty"`
}

type UserWatermarkPolicyName string

const (
	UserWatermarkPolicyNamePSI        UserWatermarkPolicyName = "PSI"
	UserWatermarkPolicyNameRefault    UserWatermarkPolicyName = "Refault"
	UserWatermarkPolicyNameIntegrated UserWatermarkPolicyName = "Integrated"
)

type UserWatermarkPSIPolicyConf struct {
	// PSIAvg60Threshold indicates the threshold of memory pressure. If observed pressure exceeds
	// this threshold, memory reclaiming will be paused.
	PSIAvg60Threshold *float64 `json:"psiAvg60Threshold,omitempty"`
}

type UserWatermarkRefaultPolicyConf struct {
	// ReclaimAccuracyTarget indicates the desired level of precision or accuracy in offloaded pages.
	ReclaimAccuracyTarget *float64 `json:"reclaimAccuracyTarget,omitempty"`
	// ReclaimScanEfficiencyTarget indicates the desired level of efficiency in scanning and
	// identifying memory pages that can be offloaded.
	ReclaimScanEfficiencyTarget *float64 `json:"reclaimScanEfficiencyTarget,omitempty"`
}
