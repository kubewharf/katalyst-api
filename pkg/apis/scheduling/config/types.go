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
	corev1 "k8s.io/api/core/v1"
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

// IndicatorType indicator participate in calculate score
type IndicatorType string

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LoadAwareArgs holds arguments used to configure the LoadAwareScheduling plugin.
type LoadAwareArgs struct {
	metav1.TypeMeta

	EnablePortrait               *bool   `json:"enablePortrait,omitempty"`
	PodAnnotationLoadAwareEnable *string `json:"podAnnotationLoadAwareEnable,omitempty"`
	// FilterExpiredNodeMetrics indicates whether to filter nodes where  fails to update NPD.
	FilterExpiredNodeMetrics *bool `json:"filterExpiredNodeMetrics,omitempty"`
	// NodeMetricsExpiredSeconds indicates the NodeMetrics in NPD expiration in seconds.
	// When NodeMetrics expired, the node is considered abnormal.
	// default 5 minute
	NodeMetricsExpiredSeconds *int64 `json:"NodeMetricsExpiredSeconds,omitempty"`
	// ResourceToWeightMap contains resource name and weight.
	ResourceToWeightMap map[corev1.ResourceName]int64 `json:"resourceToWeightMap,omitempty"`
	// ResourceToThresholdMap contains resource name and threshold. Node can not be scheduled
	// if usage of it is more than threshold.
	ResourceToThresholdMap map[corev1.ResourceName]int64 `json:"resourceToThresholdMap,omitempty"`
	// ResourceToScalingFactorMap contains resource name and scaling factor, which are used to estimate pod usage
	// if usage of pod is not exists in node monitor.
	ResourceToScalingFactorMap map[corev1.ResourceName]int64 `json:"resourceToScalingFactorMap,omitempty"`
	// ResourceToTargetMap contains resource name and node usage target which are used in loadAware score
	ResourceToTargetMap map[corev1.ResourceName]int64 `json:"resourceToTargetMap,omitempty"`
	// CalculateIndicatorWeight indicates the participates in calculate indicator weight
	// The default avg_15min 30, max_1hour 30, max_1day 40
	CalculateIndicatorWeight map[IndicatorType]int64 `json:"calculateIndicatorWeight,omitempty"`

	KubeConfigPath string `json:"kubeConfigPath,omitempty"`
}
