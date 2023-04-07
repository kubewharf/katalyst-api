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
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=adminqosconfigurations,shortName=aqc
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="AGE",type=date,JSONPath=.metadata.creationTimestamp
// +kubebuilder:printcolumn:name="SELECTOR",type=string,JSONPath=".spec.nodeLabelSelector"
// +kubebuilder:printcolumn:name="NODES",type=string,JSONPath=".spec.ephemeralSelector.nodeNames"
// +kubebuilder:printcolumn:name="DURATION",type=string,JSONPath=".spec.ephemeralSelector.lastDuration"
// +kubebuilder:printcolumn:name="VALID",type=string,JSONPath=".status.conditions[?(@.type==\"Valid\")].status"
// +kubebuilder:printcolumn:name="REASON",type=string,JSONPath=".status.conditions[?(@.type==\"Valid\")].reason"
// +kubebuilder:printcolumn:name="MESSAGE",type=string,JSONPath=".status.conditions[?(@.type==\"Valid\")].message"

// AdminQoSConfiguration is the Schema for the configuration API used by admin QoS policy
type AdminQoSConfiguration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AdminQoSConfigurationSpec `json:"spec,omitempty"`
	Status GenericConfigStatus       `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

// AdminQoSConfigurationList contains a list of AdminQoSConfiguration
type AdminQoSConfigurationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AdminQoSConfiguration `json:"items"`
}

// AdminQoSConfigurationSpec defines the desired state of AdminQoSConfiguration
type AdminQoSConfigurationSpec struct {
	GenericConfigSpec `json:",inline"`

	// Config is custom field for admin qos configuration
	Config AdminQoSConfig `json:"config,omitempty"`
}

type AdminQoSConfig struct {
	// ReclaimedResourceConfig is a configuration for reclaim resource
	ReclaimedResourceConfig ReclaimedResourceConfig `json:"reclaimedResourceConfig,omitempty"`
}

type ReclaimedResourceConfig struct {
	// EnableReclaim is a flag to enable reclaim resource, if true, reclaim resource will be enabled,
	// which means reclaim resource will be reported to custom node resource and support colocation between
	// reclaimed_cores pod and other pods, otherwise, reclaim resource will be disabled.
	// +optional
	EnableReclaim *bool `json:"enableReclaim,omitempty"`

	// ReservedResourceForReport is a reserved resource for report to custom node resource, which is used to
	// prevent reclaim resource from being requested by reclaimed_cores pods.
	// +optional
	ReservedResourceForReport *v1.ResourceList `json:"reservedResourceForReport,omitempty"`

	// MinReclaimedResourceForReport is a minimum reclaimed resource for report to custom node resource, which means
	// if reclaimed resource is less than MinReclaimedResourceForReport, then reclaimed resource will be reported as
	// MinReclaimedResourceForReport.
	// +optional
	MinReclaimedResourceForReport *v1.ResourceList `json:"minReclaimedResourceForReport,omitempty"`

	// ReservedResourceForAllocate is a resource reserved for non-reclaimed_cores pods that are not allocated to
	// reclaimed_cores pods. It is used to set aside some buffer resources to avoid sudden increase in resource
	// requirements.
	// +optional
	ReservedResourceForAllocate *v1.ResourceList `json:"reservedResourceForAllocate,omitempty"`
}
