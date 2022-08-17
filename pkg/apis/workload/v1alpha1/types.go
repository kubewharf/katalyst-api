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
	metrics "k8s.io/metrics/pkg/apis/metrics/v1beta1"

	"github.com/kubewharf/katalyst-api/pkg/apis/autoscaling/v1alpha1"
)

// +kubebuilder:validation:Enum=avg;max
type Aggregator string

const (
	Avg    Aggregator      = "avg"
	Max    Aggregator      = "max"
	Load1m v1.ResourceName = "load1m"
	Load1h v1.ResourceName = "load1h"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=spd

// ServiceProfileDescriptor captures information about a VPA object
type ServiceProfileDescriptor struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the behavior of a ServiceProfileDescriptor.
	// +optional
	Spec ServiceProfileDescriptorSpec `json:"spec,omitempty"`

	// Status represents the concrete samples of ServiceProfileData with multiple resources.
	// +optional
	Status ServiceProfileDescriptorStatus `json:"status,omitempty"`
}

// ServiceProfileDescriptorSpec is the specification of the behavior of the SPD.
type ServiceProfileDescriptorSpec struct {
	// TargetRef points to the controller managing the set of pods for the
	// spd-controller to control - e.g. Deployment, StatefulSet.
	TargetRef v1alpha1.CrossVersionObjectReference `json:"targetRef"`
}

type AggPodMetrics struct {
	Aggregator Aggregator           `json:"aggregator"`
	Items      []metrics.PodMetrics `json:"items"`
}

// ServiceProfileDescriptorStatus describes the aggregated metrics of the spd.
type ServiceProfileDescriptorStatus struct {
	AggMetrics []AggPodMetrics `json:"aggMetrics"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ServiceProfileDescriptorList is a collection of SPD objects.
type ServiceProfileDescriptorList struct {
	metav1.TypeMeta `json:",inline"`

	// Standard list metadata
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	// items is the list of SPDs
	Items []ServiceProfileDescriptor `json:"items"`
}
