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

package v1beta3

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/kube-scheduler/config/v1beta3"
	"k8s.io/utils/pointer"

	"github.com/kubewharf/katalyst-api/pkg/consts"
)

var defaultResourceSpec = []v1beta3.ResourceSpec{
	{Name: string(v1.ResourceCPU), Weight: 1},
	{Name: string(v1.ResourceMemory), Weight: 1},
}

var defaultReclaimedResourceSpec = []v1beta3.ResourceSpec{
	{Name: fmt.Sprintf("%s", consts.ReclaimedResourceMilliCPU), Weight: 1},
	{Name: fmt.Sprintf("%s", consts.ReclaimedResourceMemory), Weight: 1},
}

var defaultAlignedResourceSpec = []string{v1.ResourceCPU.String(), v1.ResourceMemory.String()}

var (
	defaultNodeMonitorExpiredSeconds int64 = 180
	defaultResourceToWeightMap             = map[v1.ResourceName]int64{
		v1.ResourceCPU:    1,
		v1.ResourceMemory: 1,
	}

	defaultResourceToThresholdMap = map[v1.ResourceName]int64{
		v1.ResourceCPU:    70, // 70%
		v1.ResourceMemory: 95, // 95%
	}

	defaultResourceToScalingFactorMap = map[v1.ResourceName]int64{
		v1.ResourceCPU:    85, // 85%
		v1.ResourceMemory: 70, // 70%
	}
	defaultCalculateIndicatorWeight = map[IndicatorType]int64{
		consts.Usage15MinAvgKey: 30, //30%
		consts.Usage1HourMaxKey: 30, //30%
		consts.Usage1DayMaxKey:  40, //40%
	}
	defaultResourceToTargetMap = map[v1.ResourceName]int64{
		v1.ResourceCPU:    50,
		v1.ResourceMemory: 70,
	}
)

// SetDefaults_QoSAwareNodeResourcesFitArgs sets the default parameters for QoSAwareNodeResourcesFit plugin.
func SetDefaults_QoSAwareNodeResourcesFitArgs(obj *QoSAwareNodeResourcesFitArgs) {
	if obj.ScoringStrategy == nil {
		obj.ScoringStrategy = &ScoringStrategy{
			Type:               v1beta3.LeastAllocated,
			Resources:          defaultResourceSpec,
			ReclaimedResources: defaultReclaimedResourceSpec,
		}
	}
	if len(obj.ScoringStrategy.Resources) == 0 {
		// If no resources specified, use the default set.
		obj.ScoringStrategy.Resources = append(obj.ScoringStrategy.Resources, defaultResourceSpec...)
	}
	if len(obj.ScoringStrategy.ReclaimedResources) == 0 {
		obj.ScoringStrategy.ReclaimedResources = append(obj.ScoringStrategy.ReclaimedResources, defaultReclaimedResourceSpec...)
	}
	for i := range obj.ScoringStrategy.Resources {
		if obj.ScoringStrategy.Resources[i].Weight == 0 {
			obj.ScoringStrategy.Resources[i].Weight = 1
		}
	}
	for i := range obj.ScoringStrategy.ReclaimedResources {
		if obj.ScoringStrategy.ReclaimedResources[i].Weight == 0 {
			obj.ScoringStrategy.ReclaimedResources[i].Weight = 1
		}
	}
}

// SetDefaults_QoSAwareNodeResourcesBalancedAllocationArgs sets the default parameters for QoSAwareNodeResourcesBalancedAllocation plugin.
func SetDefaults_QoSAwareNodeResourcesBalancedAllocationArgs(obj *QoSAwareNodeResourcesBalancedAllocationArgs) {
	if len(obj.Resources) == 0 {
		obj.Resources = append(obj.Resources, defaultResourceSpec...)
	}
	if len(obj.ReclaimedResources) == 0 {
		obj.ReclaimedResources = append(obj.ReclaimedResources, defaultReclaimedResourceSpec...)
	}
	// If the weight is not set or it is explicitly set to 0, then apply the default weight(1) instead.
	for i := range obj.Resources {
		if obj.Resources[i].Weight == 0 {
			obj.Resources[i].Weight = 1
		}
	}
	for i := range obj.ReclaimedResources {
		if obj.ReclaimedResources[i].Weight == 0 {
			obj.ReclaimedResources[i].Weight = 1
		}
	}
}

func SetDefaults_NodeResourceTopologyArgs(obj *NodeResourceTopologyArgs) {
	if len(obj.AlignedResources) == 0 {
		obj.AlignedResources = append(obj.AlignedResources, defaultAlignedResourceSpec...)
	}

	if obj.ScoringStrategy == nil {
		obj.ScoringStrategy = &ScoringStrategy{
			Type:      v1beta3.LeastAllocated,
			Resources: defaultResourceSpec,
		}
	}
	if len(obj.ScoringStrategy.Resources) == 0 {
		obj.ScoringStrategy.Resources = append(obj.ScoringStrategy.Resources, defaultResourceSpec...)
	}
	for i := range obj.ScoringStrategy.Resources {
		if obj.ScoringStrategy.Resources[i].Weight == 0 {
			obj.ScoringStrategy.Resources[i].Weight = 1
		}
	}
	if obj.ResourcePluginPolicy == "" {
		obj.ResourcePluginPolicy = consts.ResourcePluginPolicyNameDynamic
	}
}

func SetDefaults_LoadAwareArgs(obj *LoadAwareArgs) {
	if obj.FilterExpiredNodeMetrics == nil {
		obj.FilterExpiredNodeMetrics = pointer.BoolPtr(true)
	}
	if obj.NodeMetricsExpiredSeconds == nil {
		obj.NodeMetricsExpiredSeconds = pointer.Int64Ptr(defaultNodeMonitorExpiredSeconds)
	}
	if len(obj.ResourceToWeightMap) == 0 {
		obj.ResourceToWeightMap = defaultResourceToWeightMap
	}
	if len(obj.ResourceToThresholdMap) == 0 {
		obj.ResourceToThresholdMap = defaultResourceToThresholdMap
	}
	if len(obj.ResourceToScalingFactorMap) == 0 {
		obj.ResourceToScalingFactorMap = defaultResourceToScalingFactorMap
	}
	if len(obj.CalculateIndicatorWeight) == 0 {
		obj.CalculateIndicatorWeight = defaultCalculateIndicatorWeight
	}
	if len(obj.ResourceToTargetMap) == 0 {
		obj.ResourceToTargetMap = defaultResourceToTargetMap
	}
	if obj.EnablePortrait == nil {
		obj.EnablePortrait = pointer.Bool(false)
	}
}
