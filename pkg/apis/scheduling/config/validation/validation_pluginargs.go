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

package validation

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"
	kubeschedulerconfig "k8s.io/kubernetes/pkg/scheduler/apis/config"

	"github.com/kubewharf/katalyst-api/pkg/apis/scheduling/config"
	"github.com/kubewharf/katalyst-api/pkg/consts"
)

// ValidateQoSAwareNodeResourcesFitArgs validates that QoSAwareNodeResourcesFitArgs are set correctly.
func ValidateQoSAwareNodeResourcesFitArgs(path *field.Path, args *config.QoSAwareNodeResourcesFitArgs) error {
	var allErrs field.ErrorList

	if args.ScoringStrategy != nil {
		allErrs = append(allErrs, validateResources(args.ScoringStrategy.Resources, path.Child("resources"))...)
		allErrs = append(allErrs, validateResources(args.ScoringStrategy.ReclaimedResources, path.Child("reclaimedResources"))...)
		if args.ScoringStrategy.RequestedToCapacityRatio != nil {
			allErrs = append(allErrs, validateFunctionShape(args.ScoringStrategy.RequestedToCapacityRatio.Shape, path.Child("shape"))...)
		}
	}

	if len(allErrs) == 0 {
		return nil
	}
	return allErrs.ToAggregate()
}

// ValidateQoSAwareNodeResourcesBalancedAllocationArgs validates that QoSAwareNodeResourcesBalancedAllocationArgs are set correctly.
func ValidateQoSAwareNodeResourcesBalancedAllocationArgs(path *field.Path, args *config.QoSAwareNodeResourcesBalancedAllocationArgs) error {
	var allErrs field.ErrorList
	seenResources := sets.NewString()
	seenReclaimedResources := sets.NewString()
	for i, resource := range append(args.Resources) {
		if seenResources.Has(resource.Name) {
			allErrs = append(allErrs, field.Duplicate(path.Child("resources").Index(i).Child("name"), resource.Name))
		} else {
			seenResources.Insert(resource.Name)
		}
		if resource.Weight != 1 {
			allErrs = append(allErrs, field.Invalid(path.Child("resources").Index(i).Child("weight"), resource.Weight, "must be 1"))
		}
	}
	for i, resource := range append(args.ReclaimedResources) {
		if seenReclaimedResources.Has(resource.Name) {
			allErrs = append(allErrs, field.Duplicate(path.Child("reclaimedResources").Index(i).Child("name"), resource.Name))
		} else {
			seenReclaimedResources.Insert(resource.Name)
		}
		if resource.Weight != 1 {
			allErrs = append(allErrs, field.Invalid(path.Child("reclaimedResources").Index(i).Child("weight"), resource.Weight, "must be 1"))
		}
	}
	return allErrs.ToAggregate()
}

func validateResources(resources []kubeschedulerconfig.ResourceSpec, p *field.Path) field.ErrorList {
	var allErrs field.ErrorList
	for i, resource := range resources {
		if resource.Weight <= 0 || resource.Weight > 100 {
			msg := fmt.Sprintf("resource weight of %v not in valid range (0, 100]", resource.Name)
			allErrs = append(allErrs, field.Invalid(p.Index(i).Child("weight"), resource.Weight, msg))
		}
	}
	return allErrs
}

func validateFunctionShape(shape []kubeschedulerconfig.UtilizationShapePoint, path *field.Path) field.ErrorList {
	const (
		minUtilization = 0
		maxUtilization = 100
		minScore       = 0
		maxScore       = int32(kubeschedulerconfig.MaxCustomPriorityScore)
	)

	var allErrs field.ErrorList

	if len(shape) == 0 {
		allErrs = append(allErrs, field.Required(path, "at least one point must be specified"))
		return allErrs
	}

	for i := 1; i < len(shape); i++ {
		if shape[i-1].Utilization >= shape[i].Utilization {
			allErrs = append(allErrs, field.Invalid(path.Index(i).Child("utilization"), shape[i].Utilization, "utilization values must be sorted in increasing order"))
			break
		}
	}

	for i, point := range shape {
		if point.Utilization < minUtilization || point.Utilization > maxUtilization {
			msg := fmt.Sprintf("not in valid range [%d, %d]", minUtilization, maxUtilization)
			allErrs = append(allErrs, field.Invalid(path.Index(i).Child("utilization"), point.Utilization, msg))
		}

		if point.Score < minScore || point.Score > maxScore {
			msg := fmt.Sprintf("not in valid range [%d, %d]", minScore, maxScore)
			allErrs = append(allErrs, field.Invalid(path.Index(i).Child("score"), point.Score, msg))
		}
	}

	return allErrs
}

var validScoringStrategy = sets.NewString(
	string(kubeschedulerconfig.MostAllocated),
	string(kubeschedulerconfig.LeastAllocated),
	string(consts.BalancedAllocation),
	string(consts.LeastNUMANodes),
)

// ValidateNodeResourceTopologyMatchArgs ...
func ValidateNodeResourceTopologyMatchArgs(path *field.Path, args *config.NodeResourceTopologyArgs) error {
	var allErrs field.ErrorList
	scoringStrategyTypePath := path.Child("scoringStrategy.type")
	if err := validateScoringStrategyType(args.ScoringStrategy.Type, scoringStrategyTypePath); err != nil {
		allErrs = append(allErrs, err)
	}

	resourcePluginPolicyPath := path.Child("resourcePluginPolicy")
	if err := validateResourcePolicy(args.ResourcePluginPolicy, resourcePluginPolicyPath); err != nil {
		allErrs = append(allErrs, err)
	}

	return allErrs.ToAggregate()
}

func validateScoringStrategyType(scoringStrategy kubeschedulerconfig.ScoringStrategyType, path *field.Path) *field.Error {
	if !validScoringStrategy.Has(string(scoringStrategy)) {
		return field.Invalid(path, scoringStrategy, "invalid ScoringStrategyType")
	}
	return nil
}

func validateResourcePolicy(resourcePolicy consts.ResourcePluginPolicyName, path *field.Path) *field.Error {
	if resourcePolicy != consts.ResourcePluginPolicyNameDynamic &&
		resourcePolicy != consts.ResourcePluginPolicyNameNative {
		return field.Invalid(path, resourcePolicy, "invalid ResourcePolicy")
	}
	return nil
}

// ValidateLoadAwareSchedulingArgs validates that LoadAwareArgs are correct.
func ValidateLoadAwareSchedulingArgs(args *config.LoadAwareArgs) error {
	if args.NodeMetricsExpiredSeconds != nil && *args.NodeMetricsExpiredSeconds <= 0 {
		return fmt.Errorf("NodeMonitorExpiredSeconds err, NodeMonitorExpiredSeconds should be a positive value")
	}
	if err := validateResourceWeights(args.ResourceToWeightMap); err != nil {
		return fmt.Errorf("ResourceWeights err, %v", err)
	}
	if err := validateResourceThresholds(args.ResourceToThresholdMap); err != nil {
		return fmt.Errorf("UsageThresholds err, %v", err)
	}
	if err := validateEstimatedResourceThresholds(args.ResourceToScalingFactorMap); err != nil {
		return fmt.Errorf("EstimatedScalingFactors err, %v", err)
	}
	for resourceName := range args.ResourceToWeightMap {
		if _, ok := args.ResourceToScalingFactorMap[resourceName]; !ok {
			return fmt.Errorf("LoadAwareValidating err, resourceName %v in ResourceWeights, but not find in EstimatedScalingFactors", resourceName)
		}
	}
	if err := validateCalculateIndicatorWeight(args.CalculateIndicatorWeight); err != nil {
		return fmt.Errorf("CalculateIndicatorWeight err, %v", err)
	}
	if err := validateResourceTargets(args.ResourceToTargetMap); err != nil {
		return fmt.Errorf("UsageTarget err, %v", err)
	}
	return nil
}
func validateResourceWeights(resources map[v1.ResourceName]int64) error {
	for resourceName, weight := range resources {
		if weight < 0 {
			return fmt.Errorf("resource Weight of %v should be a positive value, got %v", resourceName, weight)
		}
		if weight > 100 {
			return fmt.Errorf("resource Weight of %v should be less than 100, got %v", resourceName, weight)
		}
	}
	return nil
}

func validateResourceThresholds(thresholds map[v1.ResourceName]int64) error {
	for resourceName, thresholdPercent := range thresholds {
		if thresholdPercent < 0 {
			return fmt.Errorf("resource Threshold of %v should be a positive value, got %v", resourceName, thresholdPercent)
		}
		if thresholdPercent > 100 {
			return fmt.Errorf("resource Threshold of %v should be less than 100, got %v", resourceName, thresholdPercent)
		}
	}
	return nil
}

func validateEstimatedResourceThresholds(thresholds map[v1.ResourceName]int64) error {
	for resourceName, thresholdPercent := range thresholds {
		if thresholdPercent < 0 {
			return fmt.Errorf("estimated resource Threshold of %v should be a positive value, got %v", resourceName, thresholdPercent)
		}
		if thresholdPercent > 100 {
			return fmt.Errorf("estimated  resource Threshold of %v should be less than 100, got %v", resourceName, thresholdPercent)
		}
	}
	return nil
}

func validateCalculateIndicatorWeight(calculateIndicatorWeight map[config.IndicatorType]int64) error {
	for indicator, weight := range calculateIndicatorWeight {
		if weight < 0 {
			return fmt.Errorf("calculate Indicator weight of %v should be a positive value, got %v", indicator, weight)
		}
		if weight > 100 {
			return fmt.Errorf("calculate Indicator weight of %v should be less than 100, got %v", indicator, weight)
		}
	}
	return nil
}

func validateResourceTargets(targets map[v1.ResourceName]int64) error {
	for resourceName, target := range targets {
		if target < 0 {
			return fmt.Errorf("resource target of %v should be a positive value, got %v", resourceName, target)
		}

		if target > 100 {
			return fmt.Errorf("resource target of %v should be less than 100, got %v", resourceName, target)
		}
	}
	return nil
}
