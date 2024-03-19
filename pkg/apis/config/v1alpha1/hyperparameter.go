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

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=hyperparameterconfigurations,shortName=hpc
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="AGE",type=date,JSONPath=.metadata.creationTimestamp
// +kubebuilder:printcolumn:name="SELECTOR",type=string,JSONPath=".spec.nodeLabelSelector"
// +kubebuilder:printcolumn:name="PRIORITY",type=string,JSONPath=".spec.priority"
// +kubebuilder:printcolumn:name="NODES",type=string,JSONPath=".spec.ephemeralSelector.nodeNames"
// +kubebuilder:printcolumn:name="DURATION",type=string,JSONPath=".spec.ephemeralSelector.lastDuration"
// +kubebuilder:printcolumn:name="VALID",type=string,JSONPath=".status.conditions[?(@.type==\"Valid\")].status"
// +kubebuilder:printcolumn:name="REASON",type=string,JSONPath=".status.conditions[?(@.type==\"Valid\")].reason"
// +kubebuilder:printcolumn:name="MESSAGE",type=string,JSONPath=".status.conditions[?(@.type==\"Valid\")].message"

// HyperParameterConfiguration is the Schema for the configuration API used by Katalyst policies
type HyperParameterConfiguration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HyperParameterConfigurationSpec `json:"spec,omitempty"`
	Status GenericConfigStatus             `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

// HyperParameterConfigurationList contains a list of HyperParameterConfiguration
type HyperParameterConfigurationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []HyperParameterConfiguration `json:"items"`
}

// HyperParameterConfigurationSpec defines the desired state of HyperParameterConfiguration
type HyperParameterConfigurationSpec struct {
	GenericConfigSpec `json:",inline"`

	Config HyperParameterConfig `json:"config,omitempty"`
}

type HyperParameterConfig struct {
	// BorweinProvisionConfig is configuration related to Borwein provision policy
	// +optional
	BorweinProvisionConfig *BorweinProvisionConfig `json:"borweinProvisionConfig,omitempty"`
}

type BorweinProvisionConfig struct {
	// BorweinProvisionParams is the list of parameter set for Borwein provision policy
	// +optional
	BorweinProvisionParams []BorweinProvisionParam `json:"borweinProvisionParams,omitempty"`
}

type BorweinProvisionParam struct {
	// +optional
	ModelAbnormalRatioThreshold float64 `json:"modelAbnormalRatioThreshold,omitempty"`
	// +optional
	IndicatorTargetOffsetMax float64 `json:"indicatorTargetOffsetMax,omitempty"`
	// +optional
	IndicatorTargetOffsetMin float64 `json:"indicatorTargetOffsetMin,omitempty"`
	// +optional
	IndicatorTargetRampUpStep float64 `json:"indicatorTargetRampUpStep,omitempty"`
	// +optional
	IndicatorTargetRampDownStep float64 `json:"indicatorTargetRampDownStep,omitempty"`
	// +optional
	ParamVersion int `json:"paramVersion,omitempty"`
	// +optional
	ParamID int `json:"paramID,omitempty"`
}
