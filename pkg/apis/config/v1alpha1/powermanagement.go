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
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=powermanagementconfigurations,shortName=pmc
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

// PowerManagementConfiguration is the Schema for the configuration API used by power management
type PowerManagementConfiguration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PowerManagementConfigurationSpec `json:"spec,omitempty"`
	Status GenericConfigStatus              `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

// PowerManagementConfigurationList contains a list of PowerManagementConfiguration
type PowerManagementConfigurationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PowerManagementConfiguration `json:"items"`
}

// PowerManagementConfigurationSpec defines the desired state of PowerManagementConfiguration
type PowerManagementConfigurationSpec struct {
	GenericConfigSpec `json:",inline"`

	// Config is custom field for power management configuration
	Config PowerManagementConfig `json:"config"`
}

type PowerManagementConfig struct {
	// DisablePowerCapping indicates whether to disable power capping.
	// +optional
	DisablePowerCapping *bool `json:"disablePowerCapping,omitempty"`

	// DisablePowerAdvisor indicates whether to disable power advisor.
	// +optional
	DisablePowerAdvisor *bool `json:"disablePowerAdvisor,omitempty"`

	// PowerReductionRatio indicates the ratio of power reduction, expressed as a percentage (0-100).
	// For example, 10 means reducing power by 10%.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	// +optional
	PowerReductionRatio *int32 `json:"powerReductionRatio,omitempty"`
}
