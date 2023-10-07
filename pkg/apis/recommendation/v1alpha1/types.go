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
	"k8s.io/apimachinery/pkg/runtime"
)

// ResourceRecommendSpec defines the desired state of ResourceRecommend
type ResourceRecommendSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// TargetRef points to the controller managing the set of pods for the
	// recommenders to control - e.g. Deployment.
	TargetRef CrossVersionObjectReference `json:"targetRef"`

	// Controls how the recommenders computes recommended resources.
	// The resource policy may be used to set constraints on the recommendations
	// for individual containers. If not specified, the recommenders computes recommended
	// resources for all containers in the pod, without additional constraints.
	// +optional
	ResourcePolicy PodResourcePolicy `json:"resourcePolicy,omitempty"`
}

// CrossVersionObjectReference contains enough information to let you identify the referred resource.
// +structType=atomic
type CrossVersionObjectReference struct {
	// Kind of the referent
	Kind string `json:"kind"`
	// Name of the referent
	Name string `json:"name"`
	// API version of the referent
	// +optional
	APIVersion string `json:"apiVersion,omitempty"`
}

// PodResourcePolicy controls how computes the recommended resources
// for containers belonging to the pod. There can be at most one entry for every
// named container.
type PodResourcePolicy struct {
	// policy of algorithm, if no algorithm is provided, using default
	// +optional
	AlgorithmPolicy AlgorithmPolicy `json:"algorithmPolicy,omitempty"`

	// Per-container resource policies.
	// +optional
	// +patchMergeKey=containerName
	// +patchStrategy=merge
	ContainerPolicies []ContainerResourcePolicy `json:"containerPolicies,omitempty" patchStrategy:"merge" patchMergeKey:"containerName"`
}

type AlgorithmPolicy struct {
	// Recommenders of the Algorithm used to this Pod;
	// if end-user wants to define their own recommenders algorithms,
	// they should manage this field to match their recommend implementations.
	// +optional
	Recommender string `json:"recommender,omitempty"`

	// Algorithm is the recommended algorithm of choice
	// +optional
	Algorithm Algorithm `json:"algorithm,omitempty"`

	// Extensions config by key-value format.
	// +optional
	// +kubebuilder:pruning:PreserveUnknownFields
	Extensions *runtime.RawExtension `json:"extensions,omitempty"`
}

// Algorithm is the recommended algorithm
type Algorithm string

// Algorithms
const (
	// AlgorithmPercentile is percentile recommender algorithm
	AlgorithmPercentile Algorithm = "percentile"
)

// ContainerResourcePolicy controls how computes the recommended
// resources for a specific container.
type ContainerResourcePolicy struct {
	// Name of the container or DefaultContainerResourcePolicy, in which
	// case the policy is used by the containers that don't have their own
	// policy specified.
	ContainerName string `json:"containerName"`

	// ControlledResourcesPolicy controls how the recommenders computes recommended resources
	// for the container. If not specified, the recommenders computes recommended resources
	// for none of the container resources.
	ControlledResourcesPolicy []ContainerControlledResourcesPolicy `json:"controlledResourcesPolicy" patchStrategy:"merge"`

	// Specifies which resource values should be controlled.
	// The default is "RequestsOnly".
	// +kubebuilder:default:=RequestsOnly
	ControlledValues ContainerControlledValues `json:"controlledValues"`
}

type ContainerControlledResourcesPolicy struct {
	// ResourceName is the name of the resource that is controlled.
	// +kubebuilder:validation:Enum=cpu;memory
	ResourceName v1.ResourceName `json:"resourceName"`

	// MinAllowed Specifies the minimal amount of resources that will be recommended
	// for the container. The default is no minimum.
	// +optional
	MinAllowed *resource.Quantity `json:"minAllowed,omitempty"`

	// MaxAllowed Specifies the maximum amount of resources that will be recommended
	// for the container. The default is no maximum.
	// +optional
	MaxAllowed *resource.Quantity `json:"maxAllowed,omitempty"`

	// BufferPercents is used to get extra resource buffer for the given containers
	// +optional
	BufferPercents *int32 `json:"bufferPercents,omitempty"`
}

// ContainerControlledValues controls which resource value should be autoscaled.
// +kubebuilder:validation:Enum=RequestsAndLimits;RequestsOnly;LimitsOnly
type ContainerControlledValues string

const (
	// ContainerControlledValuesRequestsAndLimits means resource request and limits
	// are scaled automatically.
	ContainerControlledValuesRequestsAndLimits ContainerControlledValues = "RequestsAndLimits"
	// ContainerControlledValuesRequestsOnly means only requested resource is autoscaled.
	ContainerControlledValuesRequestsOnly ContainerControlledValues = "RequestsOnly"
	// ContainerControlledValuesLimitsOnly means only requested resource is autoscaled.
	ContainerControlledValuesLimitsOnly ContainerControlledValues = "LimitsOnly"
)

// ResourceRecommendStatus defines the observed state of ResourceRecommend
type ResourceRecommendStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// LastRecommendationTime is last recommendation generation time
	// +optional
	LastRecommendationTime metav1.Time `json:"lastRecommendationTime"`

	// RecommendResources is the last recommendation of resources computed by
	// recommenders
	// +optional
	RecommendResources RecommendResources `json:"recommendResources,omitempty"`

	// Conditions is the set of conditions required for this recommender to recommend,
	// and indicates whether those conditions are met.
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []ResourceRecommendCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`

	// ObservedGeneration used to record the generation number when status is updated.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

// RecommendResources is the recommendation of resources computed by recommenders for
// the controlled pods. Respects the container resource policy if present in the spec.
type RecommendResources struct {
	// Resources recommended by the recommenders for specific pod.
	// +optional
	PodResources []PodResources `json:"podResources,omitempty"`

	// Resources recommended by the recommenders for each container.
	// +optional
	ContainerResources []ContainerResources `json:"containerRecommendations,omitempty"`
}

// PodResources is the recommendation of resources computed by
// recommenders for a specific pod. Respects the container resource policy
// if present in the spec.
type PodResources struct {
	// Name of the pod.
	PodName string `json:"podName"`
	// Resources recommended by the recommenders for each container.
	ContainerResources []ContainerResources `json:"containerRecommendations,omitempty"`
}

// ContainerResources is the recommendation of resources computed by
// recommenders for a specific container. Respects the container resource policy
// if present in the spec.
type ContainerResources struct {
	// Name of the container.
	ContainerName string `json:"containerName"`
	// Requests indicates the recommenders resources for requests of this container
	// +optional
	Requests *ContainerResourceList `json:"requests,omitempty"`
	// Limits indicates the recommendation resources for limits of this container
	// +optional
	Limits *ContainerResourceList `json:"limits,omitempty"`
}

// ContainerResourceList is used to represent the resourceLists
type ContainerResourceList struct {
	// Current indicates the real resource configuration from the view of CRI interface.
	// +optional
	Current v1.ResourceList `json:"current,omitempty"`
	// Recommended amount of resources. Observes ContainerResourcePolicy.
	// +optional
	Target v1.ResourceList `json:"target,omitempty"`
	// The most recent recommended resources target computed by the recommender
	// for the controlled pods, based only on actual resource usage, not taking
	// into account the ContainerResourcePolicy (UsageBuffer).
	// Used only as status indication, will not affect actual resource assignment.
	// +optional
	UncappedTarget v1.ResourceList `json:"uncappedTarget,omitempty"`
}

// ResourceRecommendCondition describes the state of
// a ResourceRecommender at a certain point.
type ResourceRecommendCondition struct {
	// type describes the current condition
	Type ResourceRecommendConditionType `json:"type"`
	// status is the status of the condition (True, False, Unknown)
	Status v1.ConditionStatus `json:"status"`
	// lastTransitionTime is the last time the condition transitioned from
	// one status to another
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// reason is the reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty"`
	// message is a human-readable explanation containing details about
	// the transition
	// +optional
	Message string `json:"message,omitempty"`
}

// ResourceRecommendConditionType are the valid conditions of a ResourceRecommend.
type ResourceRecommendConditionType string

const (
	// Validated indicates if initial validation is done
	Validated ResourceRecommendConditionType = "Validated"
	// Initialized indicates if the initialization is done
	Initialized ResourceRecommendConditionType = "Initialized"
	// RecommendationProvided indicates that recommender values have been provided
	RecommendationProvided ResourceRecommendConditionType = "RecommendationProvided"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=rec

// ResourceRecommend is the Schema for the resourcerecommends API
type ResourceRecommend struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ResourceRecommendSpec   `json:"spec,omitempty"`
	Status ResourceRecommendStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ResourceRecommendList contains a list of ResourceRecommend
type ResourceRecommendList struct {
	metav1.TypeMeta `json:",inline"`

	// Standard list metadata
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []ResourceRecommend `json:"items"`
}
