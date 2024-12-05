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
	"github.com/kubewharf/katalyst-api/pkg/consts"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=transparentmemoryoffloadingconfigurations,shortName=tmo
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

// TransparentMemoryOffloadingConfiguration is the Schema for the configuration API used by Transparent Memory Offloading
type TransparentMemoryOffloadingConfiguration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TransparentMemoryOffloadingConfigurationSpec `json:"spec,omitempty"`
	Status GenericConfigStatus                          `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

// TransparentMemoryOffloadingConfigurationList contains a list of TransparentMemoryOffloadingConfiguration
type TransparentMemoryOffloadingConfigurationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TransparentMemoryOffloadingConfiguration `json:"items"`
}

// TransparentMemoryOffloadingConfigurationSpec defines the desired state of TransparentMemoryOffloadingConfiguration
type TransparentMemoryOffloadingConfigurationSpec struct {
	GenericConfigSpec `json:",inline"`

	// Config is custom field for TMO configuration
	Config TransparentMemoryOffloadingConfig `json:"config"`
}

type TransparentMemoryOffloadingConfig struct {
	// QoSLevelConfig is a configuration for manipulating TMO on Different QoS Level
	// +optional
	// +listMapKey=qosLevel
	// +listType=map
	QoSLevelConfig []QoSLevelConfig `json:"qosLevelConfig,omitempty"`

	// CgroupConfig is a configuration for manipulating TMO on specified cgroups
	// +optional
	// +listMapKey=cgroupPath
	// +listType=map
	CgroupConfig []CgroupConfig `json:"CgroupConfig,omitempty"`

	// BlockConfig is a configuration for blocking tmo on specified pods.
	// +optional
	BlockConfig *BlockConfig `json:"blockConfig,omitempty"`
}

type QoSLevelConfig struct {
	// QoSLevel is either of reclaimed_cores, shared_cores, dedicated_cores, system_cores
	QoSLevel consts.QoSLevel `json:"qosLevel"`

	// ConfigDetail is configuration details of TMO
	ConfigDetail TMOConfigDetail `json:"configDetail"`
}

type CgroupConfig struct {
	// CgroupPath is an cgroupV2 absolute path, e.g. /sys/fs/cgroup/hdfs
	CgroupPath string `json:"cgroupPath"`

	// ConfigDetail is configuration details of TMO
	ConfigDetail TMOConfigDetail `json:"configDetail"`
}

type TMOConfigDetail struct {
	// EnableTMO is whether to enable TMO on target objective
	// +optional
	EnableTMO *bool `json:"enableTMO,omitempty"`

	// EnableSwap is whether to enable swap to offloading anon pages
	// +optional
	EnableSwap *bool `json:"enableSwap,omitempty"`

	// Interval is the minimum duration the objectives got memory reclaimed by TMO
	// +optional
	Interval *metav1.Duration `json:"interval,omitempty"`

	// PolicyName is used to specify the policy for calculating memory offloading size
	// +optional
	PolicyName *TMOPolicyName `json:"policyName,omitempty"`

	// PSIPolicyConf is configurations of a TMO policy which reclaim memory by PSI
	// +optional
	PSIPolicyConf *PSIPolicyConf `json:"psiPolicy,omitempty"`

	// RefaultPolicy is configurations of a TMO policy which reclaim memory by refault
	// +optional
	RefaultPolicConf *RefaultPolicyConf `json:"refaultPolicy,omitempty"`
}

type TMOPolicyName string

const (
	TMOPolicyNamePSI        TMOPolicyName = "PSI"
	TMOPolicyNameRefault    TMOPolicyName = "Refault"
	TMOPolicyNameIntegrated TMOPolicyName = "Integrated"
)

type PSIPolicyConf struct {
	// MaxProbe limits the memory offloading size in one cycle, it's a ratio of memory usage.
	MaxProbe *float64 `json:"maxProbe,omitempty"`

	// PSIAvg60Threshold indicates the threshold of memory pressure. If observed pressure exceeds
	// this threshold, memory offloading will be paused.
	PSIAvg60Threshold *float64 `json:"psiAvg60Threshold,omitempty"`
}

type RefaultPolicyConf struct {
	// MaxProbe limits the memory offloading size in one cycle, it's a ratio of memory usage.
	MaxProbe *float64 `json:"maxProbe,omitempty"`
	// ReclaimAccuracyTarget indicates the desired level of precision or accuracy in offloaded pages.
	ReclaimAccuracyTarget *float64 `json:"reclaimAccuracyTarget,omitempty"`
	// ReclaimScanEfficiencyTarget indicates the desired level of efficiency in scanning and
	// identifying memory pages that can be offloaded.
	ReclaimScanEfficiencyTarget *float64 `json:"reclaimScanEfficiencyTarget,omitempty"`
}

type BlockConfig struct {
	// Labels indicates disable tmo if pods with these labels. The requirements are ORed.
	// +optional
	Labels []metav1.LabelSelectorRequirement `json:"labels,omitempty"`

	// Annotations indicates disable tmo if pods with these annotations. The requirements are ORed.
	// +optional
	Annotations []metav1.LabelSelectorRequirement `json:"annotations,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// TransparentMemoryOffloadingIndicators is indicator for transparent memory offloading
type TransparentMemoryOffloadingIndicators struct {
	metav1.TypeMeta `json:",inline"`

	// ConfigDetail is configuration details of TMO
	ConfigDetail *TMOConfigDetail `json:"configDetail,omitempty"`
}
