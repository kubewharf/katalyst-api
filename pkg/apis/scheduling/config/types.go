// Copyright 2022 The Katalyst Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeschedulerconfig "k8s.io/kubernetes/pkg/scheduler/apis/config"

	"github.com/kubewharf/katalyst-api/pkg/consts"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// QoSAwareNodeResourcesFitArgs holds arguments used to configure the QoSAwareNodeResourcesFit plugin.
type QoSAwareNodeResourcesFitArgs struct {
	metav1.TypeMeta `json:",inline"`

	// ScoringStrategy selects the node resource scoring strategy.
	ScoringStrategy *ScoringStrategy `json:"scoringStrategy,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// QoSAwareNodeResourcesBalancedAllocationArgs holds arguments used to configure QoSAwareNodeResourcesBalancedAllocation plugin.
type QoSAwareNodeResourcesBalancedAllocationArgs struct {
	metav1.TypeMeta `json:",inline"`

	// Resources to be considered when scoring.
	// The default resource set includes "cpu" and "memory", only valid weight is 1.
	Resources []kubeschedulerconfig.ResourceSpec `json:"resources,omitempty"`

	// ReclaimedResources to be considered when scoring.
	// The default resource set includes "resource.katalyst.kubewharf.io/reclaimed_millicpu"
	// and "resource.katalyst.kubewharf.io/reclaimed_memory", only valid weight is 1.
	ReclaimedResources []kubeschedulerconfig.ResourceSpec `json:"reclaimedResources,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NodeResourceTopologyArgs holds arguments used to configure the NodeResourceTopologyMatch plugin
type NodeResourceTopologyArgs struct {
	metav1.TypeMeta `json:",inline"`

	// ScoringStrategy a scoring model that determine how the plugin will score the nodes.
	ScoringStrategy *ScoringStrategy `json:"scoringStrategy,omitempty"`

	// AlignedResources are resources should be aligned for dedicated pods.
	AlignedResources []string `json:"alignedResources,omitempty"`

	// ResourcePluginPolicy are QRMPlugin resource policy to allocate topology resource for containers.
	ResourcePluginPolicy consts.ResourcePluginPolicyName `json:"resourcePluginPolicy,omitempty"`
}

// ScoringStrategy define ScoringStrategyType for node resource plugin
type ScoringStrategy struct {
	// Type selects which strategy to run.
	Type kubeschedulerconfig.ScoringStrategyType `json:"type,omitempty"`

	// Resources to consider when scoring.
	// The default resource set includes "cpu" and "memory" with an equal weight.
	// Allowed weights go from 1 to 100.
	// Weight defaults to 1 if not specified or explicitly set to 0.
	Resources []kubeschedulerconfig.ResourceSpec `json:"resources,omitempty"`

	// ReclaimedResources to consider when scoring.
	// The default resource set includes "resource.katalyst.kubewharf.io/reclaimed_millicpu"
	// and "resource.katalyst.kubewharf.io/reclaimed_memory", only valid weight is 1.
	ReclaimedResources []kubeschedulerconfig.ResourceSpec `json:"reclaimedResources,omitempty"`

	// Arguments specific to RequestedToCapacityRatio strategy.
	RequestedToCapacityRatio *kubeschedulerconfig.RequestedToCapacityRatioParam `json:"requestedToCapacityRatio,omitempty"`

	// Arguments specific to RequestedToCapacityRatio strategy.
	ReclaimedRequestedToCapacityRatio *kubeschedulerconfig.RequestedToCapacityRatioParam `json:"reclaimedRequestedToCapacityRatio,omitempty"`
}
