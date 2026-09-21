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

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=strategygroupconfigurations,shortName=sgc,scope=Cluster
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

// StrategyGroupConfiguration is the Schema for the configuration API used to group nodes by enabled policies
type StrategyGroupConfiguration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StrategyGroupConfigurationSpec `json:"spec,omitempty"`
	Status GenericConfigStatus            `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

// StrategyGroupConfigurationList contains a list of StrategyGroupConfiguration
type StrategyGroupConfigurationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []StrategyGroupConfiguration `json:"items"`
}

// StrategyGroupConfigurationSpec defines the desired state of StrategyGroupConfiguration
type StrategyGroupConfigurationSpec struct {
	GenericConfigSpec `json:",inline"`

	// Config is custom field for StrategyGroupConfiguration
	Config StrategyGroupConfig `json:"config"`
}

type StrategyGroupConfig struct {
	// GroupConfigs contains configs for every group.
	// Ratios must sum to at most 1.0 across all groups in legacy mode, or within each pool in pool mode.
	GroupConfigs []GroupConfig `json:"groupConfigs,omitempty"`
}

type GroupConfig struct {
	// NodePool selects a pool matching a value of the label key in spec.nodeLabelSelector.
	// Any nonempty NodePool enables pool mode for the entire StrategyGroupConfiguration.
	// In pool mode, empty selects nodes without that label; nodes are never borrowed across pools.
	// +optional
	NodePool string `json:"nodePool,omitempty"`

	// ExperimentDimension selects a specification.<field> or KCNR nodeResourceProperties property for this pool.
	// Requires pool mode. All groups in the same pool must use the same dimension, including an empty value.
	// Empty keeps ordinary strategy grouping; nonempty varies the selected hardware dimension across groups.
	// +optional
	ExperimentDimension string `json:"experimentDimension,omitempty"`

	// ExperimentValue selects this group's value of ExperimentDimension.
	// If omitted, the allocator chooses an unused observed value and persists the binding.
	// +optional
	ExperimentValue string `json:"experimentValue,omitempty"`

	// Ratio is the desired fraction of nodes for this group; the resulting node count is rounded down.
	// Legacy mode uses the cluster node count; pool mode uses the eligible node count in the selected pool.
	Ratio *float64 `json:"ratio"`

	// EnabledStrategies are strategies enabled in this group
	EnabledStrategies []Strategy `json:"enabledStrategies"`
}
