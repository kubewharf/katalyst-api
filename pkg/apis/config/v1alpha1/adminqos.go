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

	"github.com/kubewharf/katalyst-api/pkg/apis/workload/v1alpha1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=adminqosconfigurations,shortName=aqc
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="AGE",type=date,JSONPath=.metadata.creationTimestamp
// +kubebuilder:printcolumn:name="PAUSED",type=boolean,JSONPath=".spec.paused"
// +kubebuilder:printcolumn:name="SELECTOR",type=string,JSONPath=".spec.nodeLabelSelector"
// +kubebuilder:printcolumn:name="PRIORITY",type=string,JSONPath=".spec.priority"
// +kubebuilder:printcolumn:name="NODES",type=string,JSONPath=".spec.ephemeralSelector.nodeNames"
// +kubebuilder:printcolumn:name="DURATION",type=string,JSONPath=".spec.ephemeralSelector.lastDuration"
// +kubebuilder:printcolumn:name="TARGET",type=integer,JSONPath=".status.targetNodes"
// +kubebuilder:printcolumn:name="CANARY",type=integer,JSONPath=".status.canaryNodes"
// +kubebuilder:printcolumn:name="UPDATED-TARGET",type=integer,JSONPath=".status.updatedTargetNodes"
// +kubebuilder:printcolumn:name="HASH",type=string,JSONPath=".status.currentHash"
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
	Config AdminQoSConfig `json:"config"`
}

type AdminQoSConfig struct {
	// ReclaimedResourceConfig is a configuration for reclaim resource
	// +optional
	ReclaimedResourceConfig *ReclaimedResourceConfig `json:"reclaimedResourceConfig,omitempty"`

	// QRMPluginConfig is a configuration for qrm plugin
	// +optional
	QRMPluginConfig *QRMPluginConfig `json:"qrmPluginConfig,omitempty"`

	// EvictionConfig is a configuration for eviction
	// +optional
	EvictionConfig *EvictionConfig `json:"evictionConfig,omitempty"`

	// +optional
	AdvisorConfig *AdvisorConfig `json:"advisorConfig,omitempty"`

	// FineGrainedResourceConfig is a configuration for more fine-grained resource control
	// +optional
	FineGrainedResourceConfig *FineGrainedResourceConfig `json:"fineGrainedResourceConfig,omitempty"`
}

type ReclaimedResourceConfig struct {
	// EnableReclaim is a flag to enable reclaim resource, if true, reclaim resource will be enabled,
	// which means reclaim resource will be reported to custom node resource and support colocation between
	// reclaimed_cores pod and other pods, otherwise, reclaim resource will be disabled.
	// +optional
	EnableReclaim *bool `json:"enableReclaim,omitempty"`

	// DisableReclaimSharePools is a list of share cpuset_pool that reclaim resource will be disabled.
	// default is empty, which means all share cpuset_pool will be enabled.
	// +optional
	DisableReclaimSharePools []string `json:"disableReclaimSharePools,omitempty"`

	// ReservedResourceForReport is a reserved resource for report to custom node resource, which is used to
	// prevent reclaim resource from being requested by reclaimed_cores pods.
	// For example, {"cpu": 0, "memory": 0Gi}.
	// +optional
	ReservedResourceForReport *v1.ResourceList `json:"reservedResourceForReport,omitempty"`

	// MinReclaimedResourceForReport is a minimum reclaimed resource for report to custom node resource, which means
	// if reclaimed resource is less than MinReclaimedResourceForReport, then reclaimed resource will be reported as
	// MinReclaimedResourceForReport.
	// For example, {"cpu": 4, "memory": 5Gi}.
	// +optional
	MinReclaimedResourceForReport *v1.ResourceList `json:"minReclaimedResourceForReport,omitempty"`

	// MinIgnoredReclaimedResourceForReport defines per-resource minimum thresholds. If ANY resource's current reclaimed amount
	// falls below its respective threshold, ALL reclaimed resources will be ignored and reported as zero. This prevents resource
	// fragmentation in quota calculations by avoiding reporting insignificant reclaimed quantities.
	// For example, {"cpu": 0.1, "memory": 100Mi}.
	// +optional
	MinIgnoredReclaimedResourceForReport *v1.ResourceList `json:"minIgnoredReclaimedResourceForReport,omitempty"`

	// ReservedResourceForAllocate is a resource reserved for non-reclaimed_cores pods that are not allocated to
	// reclaimed_cores pods. It is used to set aside some buffer resources to avoid sudden increase in resource
	// requirements.
	// For example, {"cpu": 4, "memory": 5Gi}.
	// +optional
	ReservedResourceForAllocate *v1.ResourceList `json:"reservedResourceForAllocate,omitempty"`

	// MinReclaimedResourceForAllocate is a resource reserved for reclaimed_cores pods，these resources will not be used
	// by shared_cores pods.
	// For example, {"cpu": 4, "memory": 0Gi}.
	// +optional
	MinReclaimedResourceForAllocate *v1.ResourceList `json:"minReclaimedResourceForAllocate,omitempty"`

	// NumaMinReclaimedResourceRatioForAllocate is a resource reserved ratio at numa level for reclaimed_cores pods，these resources will not be used
	// by shared_cores pods.
	// For example, {"cpu": 0.1, "memory": 0.2}.
	// +optional
	NumaMinReclaimedResourceRatioForAllocate *v1.ResourceList `json:"numaMinReclaimedResourceRatioForAllocate,omitempty"`

	// NumaMinReclaimedResourceForAllocate is a resource reserved at numa level for reclaimed_cores pods，these resources will not be used
	// by shared_cores pods.
	// For example, {"cpu": 2, "memory": 0Gi}.
	// +optional
	NumaMinReclaimedResourceForAllocate *v1.ResourceList `json:"numaMinReclaimedResourceForAllocate,omitempty"`

	// CPUHeadroomConfig is a configuration for cpu headroom
	// +optional
	CPUHeadroomConfig *CPUHeadroomConfig `json:"cpuHeadroomConfig,omitempty"`

	// MemoryHeadroomConfig is a configuration for memory headroom
	// +optional
	MemoryHeadroomConfig *MemoryHeadroomConfig `json:"memoryHeadroomConfig,omitempty"`
}

type MemoryHeadroomConfig struct {
	// MemoryHeadroomUtilBasedConfig is a config for utilization based memory headroom policy
	// +optional
	UtilBasedConfig *MemoryHeadroomUtilBasedConfig `json:"utilBasedConfig,omitempty"`
}

type AdvisorConfig struct {
	// +optional
	CPUAdvisorConfig *CPUAdvisorConfig `json:"cpuAdvisorConfig,omitempty"`
	// +optional
	MemoryAdvisorConfig *MemoryAdvisorConfig `json:"memoryAdvisorConfig,omitempty"`
}

type CPUAdvisorConfig struct {
	// AllowSharedCoresOverlapReclaimedCores is a flag, when enabled,
	// we will rely on kernel features to ensure that shared_cores pods can suppress and preempt reclaimed_cores pods.
	// +optional
	AllowSharedCoresOverlapReclaimedCores *bool `json:"allowSharedCoresOverlapReclaimedCores,omitempty"`

	// optional
	CPUProvisionConfig *CPUProvisionConfig `json:"cpuProvisionConfig"`
}

// QoSRegionType declares pre-defined region types
type QoSRegionType string

const (
	// QoSRegionTypeShare for each share pool
	QoSRegionTypeShare QoSRegionType = "share"

	QoSRegionTypeDedicated QoSRegionType = "dedicated"

	// QoSRegionTypeIsolation for each isolation pool
	QoSRegionTypeIsolation QoSRegionType = "isolation"

	// QoSRegionTypeDedicatedNumaExclusive for each dedicated core with numa binding
	// and numa exclusive container
	// deprecated, will be removed later, use QoSRegionTypeDedicated instead
	QoSRegionTypeDedicatedNumaExclusive QoSRegionType = "dedicated-numa-exclusive"
)

// ControlKnobName defines available control knob key for provision policy
type ControlKnobName string

const (
	// ControlKnobNonReclaimedCPURequirement refers to cpu requirement of non-reclaimed workloads, like shared_cores and dedicated_cores
	// deprecated, will be removed later
	ControlKnobNonReclaimedCPURequirement ControlKnobName = "non-reclaimed-cpu-requirement"

	// ControlKnobNonIsolatedUpperCPUSize refers to the upper cpu size, for isolated pods now
	ControlKnobNonIsolatedUpperCPUSize ControlKnobName = "isolated-upper-cpu-size"

	// ControlKnobNonIsolatedLowerCPUSize refers to the lower cpu size, for isolated pods now
	ControlKnobNonIsolatedLowerCPUSize ControlKnobName = "isolated-lower-cpu-size"

	// ControlKnobReclaimedCoresCPUQuota is cpu limit for reclaimed-cores root cgroup
	ControlKnobReclaimedCoresCPUQuota ControlKnobName = "reclaimed-cores-cpu-quota"

	// ControlKnobReclaimedCoresCPUSize is the length of cpuset.cpus for reclaimed-cores
	ControlKnobReclaimedCoresCPUSize ControlKnobName = "reclaimed-cores-cpu-size"
)

type RegionIndicators struct {
	RegionType QoSRegionType                  `json:"regionType"`
	Targets    []IndicatorTargetConfiguration `json:"targets"`
}

type IndicatorTargetConfiguration struct {
	Name   v1alpha1.ServiceSystemIndicatorName `json:"name"`
	Target float64                             `json:"target"`
}

type RestrictConstraints struct {
	// MaxUpperGap is the maximum upward offset value from the baseline
	MaxUpperGap *float64 `json:"maxUpperGap,omitempty"`
	// MaxLowerGap is the maximum downward offset value from the baseline
	MaxLowerGap *float64 `json:"maxLowerGap,omitempty"`
	// MaxUpperGapRatio is the maximum upward offset ratio from the baseline
	MaxUpperGapRatio *float64 `json:"maxUpperGapRatio,omitempty"`
	// MaxLowerGapRatio is the maximum downward offset ratio from the baseline
	MaxLowerGapRatio *float64 `json:"maxLowerGapRatio,omitempty"`
}

type ControlKnobConstraints struct {
	Name                ControlKnobName `json:"name"`
	RestrictConstraints `json:",inline"`
}

type CPUProvisionConfig struct {
	RegionIndicators []RegionIndicators       `json:"regionIndicators,omitempty"`
	Constraints      []ControlKnobConstraints `json:"constraints,omitempty"`
}

type MemoryAdvisorConfig struct {
	// MemoryGuardConfig is a config for memory guard plugin, which is used to avoid high priority workload from being
	// affected by memory bursting caused by low priority workload.
	// +optional
	MemoryGuardConfig *MemoryGuardConfig `json:"memoryGuardConfig,omitempty"`
}

type MemoryGuardConfig struct {
	// Enable is a flag to enable memory guard plugin
	// +optional
	Enable *bool `json:"enable,omitempty"`

	// +kubebuilder:validation:Minimum=0
	// +optional
	CriticalWatermarkScaleFactor *float64 `json:"criticalWatermarkScaleFactor,omitempty"`
}

type MemoryHeadroomUtilBasedConfig struct {
	// Enable is a flag to enable utilization based memory headroom policy
	// +optional
	Enable *bool `json:"enable,omitempty"`

	// FreeBasedRatio is the estimation of free memory utilization, which can
	// be used as system buffer to oversold memory.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=1
	// +optional
	FreeBasedRatio *float64 `json:"freeBasedRatio,omitempty"`

	// StaticBasedCapacity is the static oversold memory size by bytes
	// +kubebuilder:validation:Minimum=0
	// +optional
	StaticBasedCapacity *float64 `json:"staticBasedCapacity,omitempty"`

	// CacheBasedRatio is the rate of cache oversold, 0 means disable cache oversold
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=1
	// +optional
	CacheBasedRatio *float64 `json:"cacheBasedRatio,omitempty"`

	// MaxOversoldRate is the max oversold rate of memory headroom to the memory limit of
	// reclaimed_cores cgroup
	// +kubebuilder:validation:Minimum=0
	// +optional
	MaxOversoldRate *float64 `json:"maxOversoldRate,omitempty"`
}

type CPUHeadroomConfig struct {
	// UtilBasedConfig is a config for utilization based cpu headroom policy
	// +optional
	UtilBasedConfig *CPUHeadroomUtilBasedConfig `json:"utilBasedConfig,omitempty"`
}

type CPUHeadroomUtilBasedConfig struct {
	// Enable is a flag to enable utilization based cpu headroom policy
	// +optional
	Enable *bool `json:"enable,omitempty"`

	// TargetReclaimedCoreUtilization is the target reclaimed core utilization to be used for
	// calculating the oversold cpu headroom
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=1
	// +optional
	TargetReclaimedCoreUtilization *float64 `json:"targetReclaimedCoreUtilization,omitempty"`

	// MaxReclaimedCoreUtilization is the max reclaimed core utilization of reclaimed_cores pool,
	// which is used to calculate the oversold cpu headroom, if zero means no limit
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=1
	// +optional
	MaxReclaimedCoreUtilization *float64 `json:"maxReclaimedCoreUtilization,omitempty"`

	// MaxOversoldRate is the max oversold rate of cpu headroom to the actual size of
	// reclaimed_cores pool
	// +kubebuilder:validation:Minimum=0
	// +optional
	MaxOversoldRate *float64 `json:"maxOversoldRate,omitempty"`

	// MaxHeadroomCapacityRate is the max headroom capacity rate of cpu headroom to the total
	// cpu capacity of node
	// +kubebuilder:validation:Minimum=0
	// +optional
	MaxHeadroomCapacityRate *float64 `json:"maxHeadroomCapacityRate,omitempty"`

	// NonReclaimUtilizationHigh is the high CPU utilization threshold
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=1
	// +optional
	NonReclaimUtilizationHigh *float64 `json:"nonReclaimUtilizationHigh,omitempty"`

	// NonReclaimUtilizationLow is the low CPU utilization threshold
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=1
	// +optional
	NonReclaimUtilizationLow *float64 `json:"nonReclaimUtilizationLow,omitempty"`
}

type QRMPluginConfig struct {
	// CPUPluginConfig is the config for cpu plugin
	// +optional
	CPUPluginConfig *CPUPluginConfig `json:"cpuPluginConfig,omitempty"`
}

type CPUPluginConfig struct {
	// PreferUseExistNUMAHintResult prefer to use existing numa hint results
	// The calculation results may originate from upstream components and be recorded in the pod annotation
	// +optional
	PreferUseExistNUMAHintResult *bool `json:"preferUseExistNUMAHintResult,omitempty"`
}

type EvictionConfig struct {
	// DryRun is the list of eviction plugins to dryRun
	// '*' means "all dry-run by default"
	// 'foo' means "dry-run 'foo'"
	// first item for a particular name wins
	// +optional
	DryRun []string `json:"dryRun"`

	// CPUPressureEvictionConfig is the config for cpu pressure eviction
	// +optional
	CPUPressureEvictionConfig *CPUPressureEvictionConfig `json:"cpuPressureEvictionConfig,omitempty"`

	// SystemLoadPressureEvictionConfig is the config for system load eviction
	// +optional
	//
	// Deprecated: Please use CPUSystemPressureEvictionConfig instead to configure params for CPU eviction plugin at node level
	SystemLoadPressureEvictionConfig *SystemLoadPressureEvictionConfig `json:"systemLoadPressureEvictionConfig,omitempty"`

	// CPUSystemPressureEvictionConfig is the config for cpu system pressure eviction at node level
	// +optional
	CPUSystemPressureEvictionConfig *CPUSystemPressureEvictionConfig `json:"cpuSystemPressureEvictionConfig,omitempty"`

	// MemoryPressureEvictionConfig is the config for memory pressure eviction
	// +optional
	MemoryPressureEvictionConfig *MemoryPressureEvictionConfig `json:"memoryPressureEvictionConfig,omitempty"`

	// RootfsPressureEvictionConfig is the config for rootfs pressure eviction
	// +optional
	RootfsPressureEvictionConfig *RootfsPressureEvictionConfig `json:"rootfsPressureEvictionConfig,omitempty"`

	// ReclaimedResourcesEvictionConfig is the config for reclaimed resources' eviction
	// +optional
	ReclaimedResourcesEvictionConfig *ReclaimedResourcesEvictionConfig `json:"reclaimedResourcesEvictionConfig,omitempty"`

	// NetworkEvictionConfig is the config for network eviction
	// +optional
	NetworkEvictionConfig *NetworkEvictionConfig `json:"networkEvictionConfig,omitempty"`
}

type ReclaimedResourcesEvictionConfig struct {
	// EvictionThreshold eviction threshold rate for reclaimed resources
	// +optional
	EvictionThreshold map[v1.ResourceName]float64 `json:"evictionThreshold"`

	// SoftEvictionThreshold soft eviction threshold rate for reclaimed resources
	// +optional
	SoftEvictionThreshold map[v1.ResourceName]float64 `json:"softEvictionThreshold"`

	// GracePeriod is the grace period of reclaimed resources' eviction
	// +kubebuilder:validation:Minimum=0
	// +optional
	GracePeriod *int64 `json:"gracePeriod,omitempty"`

	// ThresholdMetToleranceDuration is the tolerance duration before eviction.
	// +kubebuilder:validation:Minimum=0
	// +optional
	ThresholdMetToleranceDuration *int64 `json:"thresholdMetToleranceDuration,omitempty"`
}

type CPUPressureEvictionConfig struct {
	// EnableLoadEviction is whether to enable cpu load eviction
	// +optional
	EnableLoadEviction *bool `json:"enableLoadEviction,omitempty"`

	// LoadUpperBoundRatio is the upper bound ratio of cpuset pool load, if the load
	// of the target cpuset pool is greater than the load upper bound repeatedly, the
	// eviction will be triggered
	// +kubebuilder:validation:Minimum=1
	// +optional
	LoadUpperBoundRatio *float64 `json:"loadUpperBoundRatio,omitempty"`

	// LoadLowerBoundRatio is the lower bound ratio of cpuset pool load, if the load
	// of the target cpuset pool is greater than the load lower bound repeatedly, the
	// node taint will be triggered
	// +kubebuilder:validation:Minimum=0
	// +optional
	LoadLowerBoundRatio *float64 `json:"loadLowerBoundRatio,omitempty"`

	// LoadThresholdMetPercentage is the percentage of the number of times the load
	// over the upper bound to the total number of times the load is measured, if the
	// percentage is greater than the load threshold met percentage, the eviction or
	// node tainted will be triggered
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=1
	// +optional
	LoadThresholdMetPercentage *float64 `json:"loadThresholdMetPercentage,omitempty"`

	// LoadMetricRingSize is the size of the load metric ring, which is used to calculate the
	// load of the target cpuset pool
	// +kubebuilder:validation:Minimum=1
	// +optional
	LoadMetricRingSize *int `json:"loadMetricRingSize,omitempty"`

	// LoadEvictionCoolDownTime is the cool-down time of cpu load eviction,
	// if the cpu load eviction is triggered, the cpu load eviction will be
	// disabled for the cool-down time
	// +optional
	LoadEvictionCoolDownTime *metav1.Duration `json:"loadEvictionCoolDownTime,omitempty"`

	// EnableSuppressionEviction is whether to enable pod-level cpu suppression eviction
	// +optional
	EnableSuppressionEviction *bool `json:"enableSuppressionEviction,omitempty"`

	// MaxSuppressionToleranceRate is the maximum cpu suppression tolerance rate that
	// can be set by the pod, if the cpu suppression tolerance rate of the pod is greater
	// than the maximum cpu suppression tolerance rate, the cpu suppression tolerance rate
	// of the pod will be set to the maximum cpu suppression tolerance rate
	// +kubebuilder:validation:Minimum=0
	// +optional
	MaxSuppressionToleranceRate *float64 `json:"maxSuppressionToleranceRate,omitempty"`

	// MinSuppressionToleranceDuration is the minimum duration a pod can tolerate cpu
	// suppression, only if the cpu suppression duration of the pod is greater than the
	// minimum cpu suppression duration, the eviction will be triggered
	// +optional
	MinSuppressionToleranceDuration *metav1.Duration `json:"minSuppressionToleranceDuration,omitempty"`

	// GracePeriod is the grace period of cpu pressure eviction
	// +kubebuilder:validation:Minimum=0
	// +optional
	GracePeriod *int64 `json:"gracePeriod,omitempty"`

	// NumaCPUPressureEvictionConfig holds configurations for NUMA-level CPU pressure eviction.
	NumaCPUPressureEvictionConfig NumaCPUPressureEvictionConfig `json:"numaCPUPressureEvictionConfig,omitempty"`

	// NumaSysCPUPressureEvictionConfig holds configurations for NUMA-level system CPU pressure eviction.
	NumaSysCPUPressureEvictionConfig NumaSysCPUPressureEvictionConfig `json:"numaSysCPUPressureEvictionConfig,omitempty"`
}

// NumaCPUPressureEvictionConfig holds the configurations for NUMA-level CPU pressure eviction.
type NumaCPUPressureEvictionConfig struct {
	// EnableEviction indicates whether to enable NUMA-level CPU pressure eviction.
	// +optional
	EnableEviction *bool `json:"enableEviction,omitempty"`

	// ThresholdMetPercentage is the percentage of time the NUMA's CPU pressure
	// must be above the threshold for an eviction to be triggered.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=1
	// +optional
	ThresholdMetPercentage *float64 `json:"thresholdMetPercentage,omitempty"`

	// MetricRingSize is the size of the metric ring buffer for calculating NUMA CPU pressure.
	// +kubebuilder:validation:Minimum=1
	// +optional
	MetricRingSize *int `json:"metricRingSize,omitempty"`

	// GracePeriod is the grace period (in seconds) after a pod starts before it can be considered for eviction
	// due to NUMA CPU pressure. 0 means no grace period.
	// +kubebuilder:validation:Minimum=0
	// +optional
	GracePeriod *int64 `json:"gracePeriod,omitempty"`

	// ThresholdExpandFactor expands the metric threshold from a specific machine to set the eviction threshold.
	// E.g., 1.1 means a 10% increase.
	// +kubebuilder:validation:Minimum=0
	// +optional
	ThresholdExpandFactor *float64 `json:"thresholdExpandFactor,omitempty"`

	// CandidateCount is the candidate pod count when selecting pods to be evicted.
	// +kubebuilder:validation:Minimum=0
	// +optional
	CandidateCount *int `json:"candidateCount,omitempty"`

	// SkippedPodKinds is the pod kind that will be skipped when selecting pods to be evicted.
	// +optional
	SkippedPodKinds []string `json:"skippedPodKinds,omitempty"`
}

// NumaSysCPUPressureEvictionConfig holds the configurations for NUMA-level system CPU pressure eviction.
type NumaSysCPUPressureEvictionConfig struct {
	// EnableEviction indicates whether to enable NUMA-level system CPU pressure eviction.
	// +optional
	EnableEviction *bool `json:"enableEviction,omitempty"`
	// MetricRingSize is the size of the metric ring buffer for calculating NUMA system CPU pressure.
	// +optional
	MetricRingSize *int `json:"metricRingSize,omitempty"`

	// GracePeriod is the grace period (in seconds) after a pod starts before it can be considered for eviction
	// due to NUMA system CPU pressure. 0 means no grace period.
	// +optional
	GracePeriod *int64 `json:"gracePeriod,omitempty"`
	// SyncPeriod is the sync period (in seconds) for updating NUMA system CPU pressure metrics.
	// +optional
	SyncPeriod *int64 `json:"syncPeriod,omitempty"`

	// ThresholdMetPercentage is the percentage of time the NUMA's system CPU pressure
	// must be above the threshold for an eviction to be triggered.
	// +optional
	ThresholdMetPercentage *float64 `json:"thresholdMetPercentage,omitempty"`
	// NumaCPUUsageSoftThreshold is the soft threshold for NUMA system CPU pressure.
	// +optional
	NumaCPUUsageSoftThreshold *float64 `json:"numaCPUUsageSoftThreshold,omitempty"`
	// NumaCPUUsageHardThreshold is the hard threshold for NUMA system CPU pressure.
	// +optional
	NumaCPUUsageHardThreshold *float64 `json:"numaCPUUsageHardThreshold,omitempty"`
	// NUMASysOverTotalUsageSoftThreshold is the soft threshold for NUMA system CPU pressure over total system CPU pressure.
	// +optional
	NUMASysOverTotalUsageSoftThreshold *float64 `json:"numaSysOverTotalUsageSoftThreshold,omitempty"`
	// NUMASysOverTotalUsageHardThreshold is the hard threshold for NUMA system CPU pressure over total system CPU pressure.
	// +optional
	NUMASysOverTotalUsageHardThreshold *float64 `json:"numaSysOverTotalUsageHardThreshold,omitempty"`
	// NUMASysOverTotalUsageEvictionThreshold is the eviction threshold for NUMA system CPU pressure over total system CPU pressure.
	// +optional
	NUMASysOverTotalUsageEvictionThreshold *float64 `json:"numaSysOverTotalUsageEvictionThreshold,omitempty"`
}

type MemoryPressureEvictionConfig struct {
	// EnableNumaLevelEviction is whether to enable numa-level eviction
	// +optional
	EnableNumaLevelEviction *bool `json:"enableNumaLevelEviction,omitempty"`

	// EnableSystemLevelEviction is whether to enable system-level eviction
	// +optional
	EnableSystemLevelEviction *bool `json:"enableSystemLevelEviction,omitempty"`

	// NumaVictimMinimumUtilizationThreshold is the victim's minimum memory usage on a NUMA node, if a pod
	// uses less memory on a NUMA node than this threshold,it won't be evicted by this NUMA's memory pressure.
	// +optional
	NumaVictimMinimumUtilizationThreshold *float64 `json:"numaVictimMinimumUtilizationThreshold,omitempty"`

	// NumaFreeBelowWatermarkTimesThreshold is the threshold for the number of
	// times NUMA's free memory falls below the watermark
	// +kubebuilder:validation:Minimum=0
	// +optional
	NumaFreeBelowWatermarkTimesThreshold *int `json:"numaFreeBelowWatermarkTimesThreshold,omitempty"`

	// NumaFreeBelowWatermarkTimesReclaimedThreshold is the threshold for the number of
	// times NUMA's free memory of the reclaimed instance falls below the watermark
	// +kubebuilder:validation:Minimum=0
	// +optional
	NumaFreeBelowWatermarkTimesReclaimedThreshold *int `json:"numaFreeBelowWatermarkTimesReclaimedThreshold,omitempty"`

	// NumaFreeConstraintFastEvictionWaitCycle is the wait cycle for fast eviction when numa memory is extremely tight
	// +kubebuilder:validation:Minimum=0
	// +optional
	NumaFreeConstraintFastEvictionWaitCycle *int `json:"numaFreeConstraintFastEvictionWaitCycle,omitempty"`

	// NumaFreeBelowWatermarkTimesThreshold is the threshold for the rate of
	// kswapd reclaiming rate
	// +kubebuilder:validation:Minimum=0
	// +optional
	SystemKswapdRateThreshold *int `json:"systemKswapdRateThreshold,omitempty"`

	// SystemKswapdRateExceedDurationThreshold is the threshold for the duration the kswapd reclaiming rate
	// exceeds the threshold
	// +kubebuilder:validation:Minimum=0
	// +optional
	SystemKswapdRateExceedDurationThreshold *int `json:"systemKswapdRateExceedDurationThreshold,omitempty"`

	// SystemFreeMemoryThresholdMinimum is the system free memory threshold minimum.
	// +optional
	SystemFreeMemoryThresholdMinimum *resource.Quantity `json:"systemFreeMemoryThresholdMinimum,omitempty"`

	// NumaEvictionRankingMetrics is the metrics used to rank pods for eviction
	// at the NUMA level
	// +kubebuilder:validation:MinItems=1
	// +optional
	NumaEvictionRankingMetrics []NumaEvictionRankingMetric `json:"numaEvictionRankingMetrics,omitempty"`

	// SystemEvictionRankingMetrics is the metrics used to rank pods for eviction
	// at the system level
	// +kubebuilder:validation:MinItems=1
	// +optional
	SystemEvictionRankingMetrics []SystemEvictionRankingMetric `json:"systemEvictionRankingMetrics,omitempty"`

	// EnableRSSOveruseEviction is whether to enable pod-level rss overuse eviction
	// +optional
	EnableRSSOveruseEviction *bool `json:"enableRSSOveruseEviction,omitempty"`

	// RSSOveruseRateThreshold is the threshold for the rate of rss
	// +kubebuilder:validation:Minimum=0
	// +optional
	RSSOveruseRateThreshold *float64 `json:"rssOveruseRateThreshold,omitempty"`

	// GracePeriod is the grace period of memory pressure eviction
	// +kubebuilder:validation:Minimum=0
	// +optional
	GracePeriod *int64 `json:"gracePeriod,omitempty"`

	// ReclaimedGracePeriod is the grace period of memory pressure reclaimed eviction
	// +kubebuilder:validation:Minimum=0
	// +optional
	ReclaimedGracePeriod *int64 `json:"reclaimedGracePeriod,omitempty"`

	// EvictNonReclaimedAnnotationSelector is a non-reclaimed pod eviction anno selector
	// +optional
	EvictNonReclaimedAnnotationSelector string `json:"evictNonReclaimedAnnotationSelector,omitempty"`

	// EvictNonReclaimedLabelSelector is a non-reclaimed pod eviction label selector
	// +optional
	EvictNonReclaimedLabelSelector string `json:"evictNonReclaimedLabelSelector,omitempty"`
}

type SystemLoadPressureEvictionConfig struct {
	// SoftThreshold is the soft threshold of system load pressure, it should be an integral multiple of 100, which means
	// the real threshold is (SoftThreshold / 100) * CoreNumber
	// +optional
	SoftThreshold *int64 `json:"softThreshold,omitempty"`

	// HardThreshold is the hard threshold of system load pressure, it should be an integral multiple of 100, which means
	// the real threshold is (SoftThreshold / 100) * CoreNumber
	// +optional
	HardThreshold *int64 `json:"hardThreshold,omitempty"`

	// HistorySize is the size of the load metric ring, which is used to calculate the system load
	// +kubebuilder:validation:Minimum=1
	// +optional
	HistorySize *int64 `json:"historySize,omitempty"`

	// SyncPeriod is the interval in seconds of the plugin fetch the load information
	// +kubebuilder:validation:Minimum=1
	// +optional
	SyncPeriod *int64 `json:"syncPeriod,omitempty"`

	// CoolDownTime is the cool-down time of the plugin evict pods
	// +kubebuilder:validation:Minimum=1
	// +optional
	CoolDownTime *int64 `json:"coolDownTime,omitempty"`

	// GracePeriod is the grace period of pod deletion
	// +optional
	GracePeriod *int64 `json:"gracePeriod,omitempty"`

	// the plugin considers the node is facing load pressure only when the ratio of load history which is greater than
	// threshold is greater than this percentage
	// +optional
	ThresholdMetPercentage *float64 `json:"thresholdMetPercentage,omitempty"`
}

type RootfsPressureEvictionConfig struct {
	// EnableRootfsPressureEviction is whether to enable rootfs pressure eviction.
	// +optional
	EnableRootfsPressureEviction *bool `json:"enableRootfsPressureEviction,omitempty"`

	// MinimumImageFsFreeThreshold is a threshold for a node.
	// Once the image rootfs free space of current node is lower than this threshold, the eviction manager will try to evict some pods.
	// For example: "200Gi", "10%".
	// +optional
	// +kubebuilder:validation:Pattern=`^(0|[1-9][0-9]*)(\.[0-9]+)?%?$|^(0|[1-9][0-9]*)([kKmMGTPeE]i?)$`
	MinimumImageFsFreeThreshold *string `json:"minimumImageFsFreeThreshold,omitempty"`

	// MinimumImageFsInodesFreeThreshold is a threshold for a node.
	// Once the image rootfs free inodes of current node is lower than this threshold, the eviction manager will try to evict some pods.
	// For example: "100000", "10%".
	// +optional
	// +kubebuilder:validation:Pattern=`^(0|[1-9]\d*)(\.\d+)?%?$|^\d+$`
	MinimumImageFsInodesFreeThreshold *string `json:"minimumImageFsInodesFreeThreshold,omitempty"`

	// PodMinimumUsedThreshold is a threshold for all pods.
	// The eviction manager will ignore this pod if its rootfs used in bytes is lower than this threshold.
	// For example: "200Gi", "1%".
	// +optional
	// +kubebuilder:validation:Pattern=`^(0|[1-9][0-9]*)(\.[0-9]+)?%?$|^(0|[1-9][0-9]*)([kKmMGTPeE]i?)$`
	PodMinimumUsedThreshold *string `json:"podMinimumUsedThreshold,omitempty"`

	// PodMinimumInodesUsedThreshold is a threshold for all pods.
	// The eviction manager will ignore this pod if its rootfs inodes used is lower than this threshold.
	// For example: "1000", "1%".
	// +optional
	// +kubebuilder:validation:Pattern=`^(0|[1-9]\d*)(\.\d+)?%?$|^\d+$`
	PodMinimumInodesUsedThreshold *string `json:"podMinimumInodesUsedThreshold,omitempty"`

	// ReclaimedQoSPodUsedPriorityThreshold is a threshold for all offline pods.
	// The eviction manager will prioritize the eviction of offline pods that reach this threshold.
	// For example: "100Gi", "1%".
	// +optional
	// +kubebuilder:validation:Pattern=`^(0|[1-9][0-9]*)(\.[0-9]+)?%?$|^(0|[1-9][0-9]*)([kKmMGTPeE]i?)$`
	ReclaimedQoSPodUsedPriorityThreshold *string `json:"reclaimedQoSPodUsedPriorityThreshold,omitempty"`

	// ReclaimedQoSPodInodesUsedPriorityThreshold is a threshold for all offline pods.
	// The eviction manager will prioritize the eviction of reclaimed pods that reach this threshold.
	// For example: "500", "1%".
	// +optional
	// +kubebuilder:validation:Pattern=`^(0|[1-9]\d*)(\.\d+)?%?$|^\d+$`
	ReclaimedQoSPodInodesUsedPriorityThreshold *string `json:"reclaimedQoSPodInodesUsedPriorityThreshold,omitempty"`

	// MinimumImageFsDiskCapacityThreshold is a threshold for all nodes.
	// The eviction manager will ignore those nodes whose image fs disk capacity is less than this threshold.
	// Fox example: "100Gi".
	MinimumImageFsDiskCapacityThreshold *resource.Quantity `json:"minimumImageFsDiskCapacityThreshold,omitempty"`

	// GracePeriod is the grace period of pod deletion
	// +optional
	GracePeriod *int64 `json:"gracePeriod,omitempty"`

	// EnableRootfsOveruseEviction is whether to enable rootfs overuse eviction.
	// +optional
	EnableRootfsOveruseEviction *bool `json:"enableRootfsOveruseEviction,omitempty"`

	// RootfsOveruseEvictionSupportedQoSLevels is the supported qos levels for rootfs overuse eviction.
	// +optional
	RootfsOveruseEvictionSupportedQoSLevels []string `json:"rootfsOveruseEvictionSupportedQoSLevels,omitempty"`

	// SharedQoSRootfsOveruseThreshold is the threshold for rootfs overuse.
	// For example: "100Gi", "10%".
	// +optional
	SharedQoSRootfsOveruseThreshold *string `json:"sharedQoSRootfsOveruseThreshold,omitempty"`

	// ReclaimedQoSRootfsOveruseThreshold is the threshold for rootfs overuse.
	// For example: "100Gi", "10%".
	// +optional
	ReclaimedQoSRootfsOveruseThreshold *string `json:"reclaimedQoSRootfsOveruseThreshold,omitempty"`
}

type NetworkEvictionConfig struct {
	// EnableNICHealthEviction is whether to enable NIC health eviction.
	// +optional
	EnableNICHealthEviction *bool `json:"enableNICNetworkEviction,omitempty"`

	// NICUnhealthyToleranceDuration is the default duration a pod can tolerate nic
	// unhealthy
	// +optional
	NICUnhealthyToleranceDuration *metav1.Duration `json:"nicUnhealthyToleranceDuration,omitempty"`

	// GracePeriod is the grace period of nic health eviction
	// +kubebuilder:validation:Minimum=0
	// +optional
	GracePeriod *int64 `json:"gracePeriod,omitempty"`
}

type CPUSystemPressureEvictionConfig struct {
	// +optional
	EnableCPUSystemPressureEviction *bool `json:"enableCPUSystemPressureEviction,omitempty"`

	// LoadUpperBoundRatio is the upper bound ratio of node, if the load
	// of the node is greater than the load upper bound repeatedly, the
	// eviction will be triggered
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=1
	// +optional
	LoadUpperBoundRatio *float64 `json:"loadUpperBoundRatio,omitempty"`

	// LoadLowerBoundRatio is the lower bound ratio of node, if the load
	// of the node is greater than the load lower bound repeatedly, the
	// cordon will be triggered
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=1
	// +optional
	LoadLowerBoundRatio *float64 `json:"loadLowerBoundRatio,omitempty"`

	// UsageUpperBoundRatio is the upper bound ratio of node, if the cpu usage
	// of the node is greater than the usage upper bound repeatedly, the
	// eviction will be triggered
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=1
	// +optional
	UsageUpperBoundRatio *float64 `json:"usageUpperBoundRatio,omitempty"`

	// UsageLowerBoundRatio is the lower bound ratio of node, if the cpu usage
	// of the node is greater than the usage lower bound repeatedly, the
	// cordon will be triggered
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=1
	// +optional
	UsageLowerBoundRatio *float64 `json:"usageLowerBoundRatio,omitempty"`

	// ThresholdMetPercentage is the percentage of the number of times the metric
	// over the upper bound to the total number of times the metric is measured, if the
	// percentage is greater than the threshold met percentage, the eviction or
	// node tainted will be triggered
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=1
	// +optional
	ThresholdMetPercentage *float64 `json:"thresholdMetPercentage,omitempty"`

	// MetricRingSize is the size of the load metric ring
	// +kubebuilder:validation:Minimum=1
	// +optional
	MetricRingSize *int `json:"metricRingSize,omitempty"`

	// EvictionCoolDownTime is the cool down duration of pod eviction
	// +optional
	EvictionCoolDownTime *metav1.Duration `json:"evictionCoolDownTime,omitempty"`

	// EvictionRankingMetrics is the metric list for ranking eviction pods
	// +optional
	EvictionRankingMetrics []string `json:"evictionRankingMetrics,omitempty"`

	// +optional
	GracePeriod *int64 `json:"gracePeriod,omitempty"`

	// +optional
	CheckCPUManager *bool `json:"checkCPUManager,omitempty"`

	// +optional
	RankingLabels map[string][]string `json:"RankingLabels,omitempty"`
}

// NumaEvictionRankingMetric is the metrics used to rank pods for eviction at the
// NUMA level
// +kubebuilder:validation:Enum=qos.pod;priority.pod;mem.total.numa.container
type NumaEvictionRankingMetric string

// SystemEvictionRankingMetric is the metrics used to rank pods for eviction at the
// system level
// +kubebuilder:validation:Enum=qos.pod;priority.pod;mem.usage.container;native.qos.pod;owner.pod
type SystemEvictionRankingMetric string

// ReclaimResourceIndicators defines the workload configuration for reclaim resource indicators.
// It works in conjunction with the ServiceProfileDescriptor's (SPD) extendedIndicator to control
// resource reclaim behavior at different levels of granularity. This allows fine-tuned control
// over when and how resources are reclaimed based on system pressure and performance considerations.
//
// Example usage:
//
// To disable resource reclaiming at the NUMA level when system pressure is detected:
//
// ```yaml
// apiVersion: workload.katalyst.kubewharf.io/v1alpha1
// kind: ServiceProfileDescriptor
// metadata:
//
//	name: example-spd
//
// spec:
//
//	targetRef:
//	  apiVersion: apps/v1
//	  kind: Deployment
//	  name: example-deployment
//	extendedIndicator:
//	- name: ReclaimResource
//	  indicators:
//	    disableReclaimLevel: "NUMA"
//
// ```
//
// In this example, when system pressure is detected, resource reclaiming will be disabled at the NUMA level
// for the specified workload, preventing performance degradation while still allowing reclaiming at finer
// granularities (e.g., Pod level).
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ReclaimResourceIndicators struct {
	metav1.TypeMeta `json:",inline"`

	// DisableReclaimLevel specifies at which level to disable reclaim resources.
	// When system pressure is detected, resource reclaiming can be disabled at different
	// granularity levels (Node, Socket, NUMA or Pod) to prevent performance degradation.
	// If not set, the default value is DisableReclaimLevelPod.
	// +optional
	DisableReclaimLevel *DisableReclaimLevel `json:"disableReclaimLevel,omitempty"`
}

// DisableReclaimLevel defines the level at which reclaim resources are disabled.
// The levels are ordered from broadest (Node) to finest (Pod) granularity.
// +kubebuilder:validation:Enum=Node;Socket;NUMA;Pod
type DisableReclaimLevel string

const (
	// DisableReclaimLevelNode disables reclaim resources at the node level.
	// This is the broadest level where all reclaim activities are disabled for the entire node.
	// Use this when you want to completely disable resource reclaiming across the whole node.
	DisableReclaimLevelNode DisableReclaimLevel = "Node"

	// DisableReclaimLevelSocket disables reclaim resources at the socket level.
	// This disables reclaim activities for an entire CPU socket and all associated resources.
	// Use this when you want to disable resource reclaiming for a specific CPU socket.
	DisableReclaimLevelSocket DisableReclaimLevel = "Socket"

	// DisableReclaimLevelNUMA disables reclaim resources at the NUMA level.
	// This disables reclaim activities for a specific NUMA node and its associated resources.
	// Use this when you want to disable resource reclaiming for a specific NUMA node.
	DisableReclaimLevelNUMA DisableReclaimLevel = "NUMA"

	// DisableReclaimLevelPod disables reclaim resources at the pod level.
	// This is the finest granularity where reclaim activities are disabled only for specific pods.
	// This is the default level if not explicitly specified.
	// Use this when you want to disable resource reclaiming only for specific pods.
	DisableReclaimLevelPod DisableReclaimLevel = "Pod"
)

type FineGrainedResourceConfig struct {
	// CPUBurstConfig has cpu burst related configurations
	// +optional
	CPUBurstConfig *CPUBurstConfig `json:"cpuBurstConfig,omitempty"`
}

type CPUBurstConfig struct {
	// EnableDedicatedCoresDefaultCPUBurst indicates whether cpu burst should be enabled by default for pods with dedicated cores.
	// If set to true, it means that cpu burst should be enabled by default for pods with dedicated cores (cpu burst value is calculated and set).
	// If set to false, it means that cpu burst should be disabled for pods with dedicated cores (cpu burst value is set to 0).
	// If set to nil, it means that no operation is done on the cpu burst value.
	// +optional
	EnableDedicatedCoresDefaultCPUBurst *bool `json:"enableDedicatedCoresDefaultCPUBurst,omitempty"`

	// DefaultCPUBurstPercent is the default cpu burst percent to be set for pods with dedicated cores.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	// +optional
	DefaultCPUBurstPercent *int64 `json:"defaultCPUBurstPercent,omitempty"`
}
