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

import "k8s.io/apimachinery/pkg/util/intstr"

type BulkheadConfig struct {
	// Enable controls whether core bulkhead is enabled.
	// +optional
	Enable *bool `json:"enable,omitempty"`
	// EnableBulkheadCpusetTopology controls whether the core bulkhead cpuset
	// topology plugin is enabled when Enable is true.
	// +optional
	EnableBulkheadCpusetTopology *bool `json:"enableBulkheadCpusetTopology,omitempty"`
	// EnableBulkheadCpusetMems controls whether the core bulkhead cpuset_mems
	// plugin is enabled when Enable is true. The plugin writes cpuset.mems for
	// reclaim NUMA buckets independently from cpuset topology.
	// +optional
	EnableBulkheadCpusetMems *bool `json:"enableBulkheadCpusetMems,omitempty"`
	// EnableBulkheadWorkqueue controls whether the core bulkhead workqueue
	// plugin is enabled when Enable is true.
	// +optional
	EnableBulkheadWorkqueue *bool `json:"enableBulkheadWorkqueue,omitempty"`
	// EnableBulkheadSystemService controls whether the core bulkhead
	// system_service plugin is enabled when Enable is true.
	// +optional
	EnableBulkheadSystemService *bool `json:"enableBulkheadSystemService,omitempty"`
	// NonReclaimPoolMinSize is the minimum CPU count kept in the non-reclaim
	// pool for bulkhead cpuset topology. When the current non-reclaim pool is
	// smaller than this value, CPUs are padded from reclaim effective candidates.
	// The padded CPUs never overlap with reserve pool CPUs.
	// +optional
	// +kubebuilder:validation:Minimum=0
	NonReclaimPoolMinSize *int64 `json:"nonReclaimPoolMinSize,omitempty"`
	// BulkheadRDTConfig is the dynamic config for RDT bulkhead plugins.
	// +optional
	BulkheadRDTConfig *BulkheadRDTConfig `json:"bulkheadRDTConfig,omitempty"`
}

type BulkheadRDTConfig struct {
	// EnableCPUList controls whether the RDT CPU list bulkhead plugin is enabled.
	// +optional
	EnableCPUList *bool `json:"enableCPUList,omitempty"`
	// EnableCAT controls whether the RDT CAT bulkhead plugin is enabled.
	// +optional
	EnableCAT *bool `json:"enableCAT,omitempty"`
	// DefaultCATWays is the default CAT way count or expression for non-root CLOS.
	// Supported expression variables are MaxCATWays and MinCATWays.
	// +optional
	// +kubebuilder:validation:XIntOrString
	DefaultCATWays *intstr.IntOrString `json:"defaultCATWays,omitempty"`
	// ClosCATWays maps pool names or CLOS IDs to CAT way counts or expressions.
	// +optional
	ClosCATWays map[string]intstr.IntOrString `json:"closCATWays,omitempty"`
	// CATPolicy controls CAT bit placement and exclusive non-overlap policy.
	// +optional
	CATPolicy *CATPolicy `json:"catPolicy,omitempty"`
}

// CATWaysExpressionVariable is a supported symbolic operand in CAT way expressions.
type CATWaysExpressionVariable string

const (
	CATWaysExpressionVariableMaxCATWays CATWaysExpressionVariable = "MaxCATWays"
	CATWaysExpressionVariableMinCATWays CATWaysExpressionVariable = "MinCATWays"
)

// CATPolicy controls CAT bit placement and exclusive non-overlap policy.
type CATPolicy struct {
	// DefaultPlacement is used by CLOS IDs without a CLOS-specific placement.
	// +optional
	DefaultPlacement *CATPlacementPolicy `json:"defaultPlacement,omitempty"`
	// ClosPlacements maps pool names, CLOS IDs or restricted trailing-prefix selectors to per-CLOS placement policy.
	// +optional
	ClosPlacements map[string]CATPlacementPolicy `json:"closPlacements,omitempty"`
	// ExclusiveClosIDs lists exact canonical CLOS IDs that must not overlap with any other configured CLOS.
	// Empty means the exclusive policy is explicitly disabled.
	// +optional
	ExclusiveClosIDs *[]string `json:"exclusiveClosIDs,omitempty"`
}

// CATPlacementPolicy controls where CAT masks are selected for one CLOS.
type CATPlacementPolicy struct {
	// AllowedBitUsages restricts candidate CAT bits by resctrl bit_usage class.
	// Empty means all ways supported by the domain are available.
	// +optional
	AllowedBitUsages []CATBitUsage `json:"allowedBitUsages,omitempty"`
	// Direction controls deterministic contiguous mask selection.
	// Empty means low.
	// +optional
	Direction CATAllocationDirection `json:"direction,omitempty"`
}

// CATBitUsage is one resctrl info/L3/bit_usage class.
// +kubebuilder:validation:Enum=*;S;H;X
type CATBitUsage string

const (
	CATBitUsageAll       CATBitUsage = "*"
	CATBitUsageSoftware  CATBitUsage = "S"
	CATBitUsageHardware  CATBitUsage = "H"
	CATBitUsageExclusive CATBitUsage = "X"
)

// CATAllocationDirection controls deterministic contiguous mask selection.
// +kubebuilder:validation:Enum=low;high
type CATAllocationDirection string

const (
	CATAllocationDirectionLow  CATAllocationDirection = "low"
	CATAllocationDirectionHigh CATAllocationDirection = "high"
)
