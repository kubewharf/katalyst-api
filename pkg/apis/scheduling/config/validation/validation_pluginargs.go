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
