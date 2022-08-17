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
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:scope=Cluster,shortName=kcnr
// +kubebuilder:subresource:status

// CustomNodeResource captures information about a custom defined node resource
// CustomNodeResource objects are non-namespaced.
type CustomNodeResource struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the behavior of a CustomNodeResource.
	// +optional
	Spec CustomNodeResourceSpec `json:"spec,omitempty"`

	// Status represents the current information about a CustomNodeResource.
	// This data may not be up-to-date.
	// +optional
	Status CustomNodeResourceStatus `json:"status,omitempty"`
}

type CustomNodeResourceSpec struct {
	// +optional
	NodeResourceProperties []*Property `json:"nodeResourceProperties,omitempty"`

	// customized taint for katalyst, which may affect partial tasks
	// +optional
	Taints []*Taint `json:"taints,omitempty"`
}

type Taint struct {
	// Required. The taint key to be applied to a node.
	Key string `json:"key,omitempty"`
	// Required. The taint value corresponding to the taint key.
	// +optional
	Value string `json:"value,omitempty"`
	// Required. The effect of the taint on pods
	// that do not tolerate the taint.
	// Valid effects are NoScheduleForReclaimedTasks.
	Effect TaintEffect `json:"effect,omitempty"`
}

type TaintEffect string

const (
	// TaintEffectNoScheduleForReclaimedTasks
	// Do not allow new pods using reclaimed resources to schedule onto the node unless they tolerate the taint,
	// but allow all pods submitted to Kubelet without going through the scheduler
	// to start, and allow all already-running pods to continue running.
	// Enforced by the scheduler.
	TaintEffectNoScheduleForReclaimedTasks TaintEffect = "NoScheduleForReclaimedTasks"
)

type Property struct {
	// property name
	PropertyName string `json:"propertyName"`

	// values of the specific property
	// +optional
	PropertyValues []string `json:"propertyValues,omitempty"`

	// values of the quantity-types property
	// +optional
	PropertyQuantity *resource.Quantity `json:"propertyQuantity,omitempty"`
}

type CustomNodeResourceStatus struct {
	// +optional
	ResourceAllocatable *v1.ResourceList `json:"resourceAllocatable,omitempty"`

	// +optional
	ResourceCapacity *v1.ResourceList `json:"resourceCapacity,omitempty"`

	// +optional
	TopologyStatus *TopologyStatus `json:"topologyStatus,omitempty"`

	// Conditions is an array of current observed cnr conditions.
	// +optional
	Conditions []CNRCondition `json:"conditions,omitempty"`
}

// CNRCondition contains condition information for a cnr.
type CNRCondition struct {
	// Type is the type of the condition.
	Type CNRConditionType `json:"type"`
	// Status of the condition, one of True, False, Unknown.
	Status v1.ConditionStatus `json:"status" `
	// Last time we got an update on a given condition.
	// +optional
	LastHeartbeatTime metav1.Time `json:"lastHeartbeatTime,omitempty"`
	// (brief) reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty"`
	// Human-readable message indicating details about last transition.
	// +optional
	Message string `json:"message,omitempty"`
}

type CNRConditionType string

const (
	CNRAgentReady    CNRConditionType = "AgentReady"
	CNRAgentNotFound CNRConditionType = "AgentNotFound"
)

type TopologyStatus struct {
	// +optional
	Sockets []*SocketStatus `json:"sockets,omitempty"`
}

type SocketStatus struct {
	SocketID int `json:"socketID"`
	// +optional
	Numas []*NumaStatus `json:"numas,omitempty"`
}

type NumaStatus struct {
	NumaID int `json:"numaID"`
	// +optional
	Allocatable *v1.ResourceList `json:"allocatable,omitempty"`
	// +optional
	Capacity *v1.ResourceList `json:"capacity,omitempty"`
	// +optional
	Allocations []*Allocation `json:"allocations,omitempty"`
}

type Allocation struct {
	Consumer string `json:"consumer"`
	// +optional
	Requests *v1.ResourceList `json:"requests,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CustomNodeResourceList is a collection of CustomNodeResource objects.
type CustomNodeResourceList struct {
	metav1.TypeMeta `json:",inline"`

	// Standard list metadata
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	// items is the list of CNRs
	Items []CustomNodeResource `json:"items"`
}
