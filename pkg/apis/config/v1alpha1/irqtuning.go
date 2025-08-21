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

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=irqtuningconfigurations,shortName=irqtc
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

// IRQTuningConfiguration is the Schema for the configuration API used by IRQ Tuning
type IRQTuningConfiguration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IRQTuningConfigurationSpec `json:"spec,omitempty"`
	Status GenericConfigStatus        `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

// IRQTuningConfigurationList contains a list of IRQTuningConfiguration
type IRQTuningConfigurationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IRQTuningConfiguration `json:"items"`
}

type IRQTuningConfigurationSpec struct {
	GenericConfigSpec `json:",inline"`

	// Config is custom field for irq tuning configuration
	Config IRQTuningConfig `json:"config"`
}

type IRQTuningConfig struct {
	// EnableTuner indicates whether to enable the interrupt tuning function.
	// +optional
	EnableTuner *bool `json:"enableTuner,omitempty"`
	// TuningPolicy represents the interrupt tuning strategy. One of Balance, Exclusive, Auto.
	// +kubebuilder:default=balance
	// +optional
	TuningPolicy TuningPolicy `json:"tuningPolicy,omitempty"`
	// TuningInterval is the interval of interrupt tuning.
	// +kubebuilder:default=5
	// +optional
	TuningInterval *int `json:"tuningInterval,omitempty"`

	// EnableRPS indicates whether to enable the RPS function.
	// Only balance policy support enable rps.
	// +optional
	EnableRPS *bool `json:"enableRPS,omitempty"`
	// EnableRPSCPUVSNicsQueue enable rps when (cpus count)/(nics queue count) greater than this config.
	// +optional
	EnableRPSCPUVSNicsQueue *float64 `json:"enableRPSCPUVSNicsQueue,omitempty"`
	// NICAffinityPolicy represents the NICs's irqs affinity sockets policy.
	// One of CompleteMap, OverallBalance, PhysicalTopo.
	// +optional
	NICAffinityPolicy NICAffinityPolicy `json:"nicAffinityPolicy,omitempty"`

	// ReniceKsoftirqd indicates whether to renice ksoftirqd process.
	// +optional
	ReniceKsoftirqd *bool `json:"reniceKsoftirqd,omitempty"`
	// KsoftirqdNice is the nice value of ksoftirqd process.
	// +optional
	KsoftirqdNice *int `json:"ksoftirqdNice,omitempty"`

	// CoresExpectedCPUUtil is the expected CPU utilization of cores.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	// +optional
	CoresExpectedCPUUtil *int `json:"coresExpectedCPUUtil,omitempty"`

	// ThroughputClassSwitch describes the switch configuration for a throughput class.
	// +optional
	ThroughputClassSwitch *ThroughputClassSwitchConfig `json:"throughputClassSwitch,omitempty"`
	// Threshold description for interrupting core network overLoad.
	// +optional
	CoreNetOverLoadThreshold *IRQCoreNetOverloadThresholds `json:"coreNetOverLoadThreshold,omitempty"`
	// Describes the constraints of the balanced configuration.
	// +optional
	LoadBalance *IRQLoadBalanceConfig `json:"loadBalance,omitempty"`
	// Configuration that requires interrupt core adjustment.
	// +optional
	CoresAdjust *IRQCoresAdjustConfig `json:"coresAdjust,omitempty"`
	// Need to adjust to interrupt exclusive core requirements.
	// +optional
	CoresExclusion *IRQCoresExclusionConfig `json:"coresExclusion,omitempty"`
}

type TuningPolicy string

const (
	// TuningPolicyBalance means to adjust the interrupt core in a balanced way.
	TuningPolicyBalance TuningPolicy = "Balance"
	// TuningPolicyExclusive means that the interrupt core will be exclusive and specially used to handle interrupts.
	TuningPolicyExclusive TuningPolicy = "Exclusive"
	// TuningPolicyAuto will select the tuning policy based on the actual status.
	TuningPolicyAuto TuningPolicy = "Auto"
)

type NICAffinityPolicy string

const (
	// NICAffinityPolicyCompleteMap means that no matter how many network cards there are, the interrupt affinity of each network card is balanced across all sockets.
	NICAffinityPolicyCompleteMap NICAffinityPolicy = "CompleteMap"
	// NICAffinityPolicyOverallBalance represents according number of nics and nics's physical topo bound numa, decide nic's irqs affinity which socket(s).
	NICAffinityPolicyOverallBalance NICAffinityPolicy = "OverallBalance"
	// NICAffinityPolicyPhysicalTopo nic's irqs affitnied socket strictly follow whose physical topology bound socket.
	NICAffinityPolicyPhysicalTopo NICAffinityPolicy = "PhysicalTopo"
)

type ThroughputClassSwitchConfig struct {
	// +optional
	LowThroughputThresholds *LowThroughputThresholds `json:"lowThroughputThresholds,omitempty"`
	// +optional
	NormalThroughputThresholds *NormalThroughputThresholds `json:"normalThroughputThresholds,omitempty"`
}

type LowThroughputThresholds struct {
	// +optional
	RxPPSThreshold *uint64 `json:"rxPPSThreshold,omitempty"`
	// +optional
	SuccessiveCount *int `json:"successiveCount,omitempty"`
}

type NormalThroughputThresholds struct {
	// +optional
	RxPPSThreshold *uint64 `json:"rxPPSThreshold,omitempty"`
	// +optional
	SuccessiveCount *int `json:"successiveCount,omitempty"`
}

// IRQCoreNetOverloadThresholds represents the threshold of interrupt core network overload.
// When there are one or more irq core's ratio of softnet_stat 3rd col time_squeeze packets / 1st col processed packets
// greater-equal IrqCoreSoftNetTimeSqueezeRatio,
// then tring to tune irq load balance first, if failed to tune irq load balance, then increase irq cores.
type IRQCoreNetOverloadThresholds struct {
	// Ratio of softnet_stat 3rd col time_squeeze packets / softnet_stat 1st col processed packets
	// +optional
	SoftNetTimeSqueezeRatio *float64 `json:"softNetTimeSqueezeRatio,omitempty"`
}

// IRQLoadBalanceConfig represents the configuration of interrupt load balance.
// When there are one or more irq core's cpu util greater-equal IrqCoreCpuUtilThresh or irq core's net load greater-equal IrqCoreNetOverloadThresholds,
// then try to tune irq load balance, that need to find at least one other irq core with relatively low cpu util, their cpu util gap MUST greater-equal IrqCoreCpuUtilGapThresh,
// if succeed to find irq cores with eligible cpu util, then start to tuning load balance,
// or increase irq cores immediately.
type IRQLoadBalanceConfig struct {
	// Interval of two successive irq load balance MUST greater-equal this interval
	// +optional
	SuccessiveTuningInterval *int `json:"successiveTuningInterval,omitempty"`
	// +optional
	Thresholds *IRQLoadBalanceTuningThresholds `json:"thresholds,omitempty"`
	// Two successive tunes whose interval is less-equal this threshold will be considered as ping-pong tunings
	// +optional
	PingPongIntervalThreshold *int `json:"pingPongIntervalThreshold,omitempty"`
	// Ping pong count greater-equal this threshold will trigger increasing irq cores
	// +optional
	PingPongCountThreshold *int `json:"pingPongCountThreshold,omitempty"`
	// Max number of irqs are permitted to be tuned from some irq cores to other cores in each time, allowed value {1, 2}
	// +kubebuilder:validation:Enum={1,2}
	// +optional
	IRQTunedNumMaxEachTime *int `json:"irqTunedNumMaxEachTime,omitempty"`
	// Max number of irq cores whose affinity irqs are permitted to tuned to other cores in each time, allowed value {1,2}
	// +kubebuilder:validation:Enum={1,2}
	// +optional
	IRQCoresTunedNumMaxEachTime *int `json:"irqCoresTunedNumMaxEachTime,omitempty"`
}

type IRQLoadBalanceTuningThresholds struct {
	// IRQ core cpu util threshold, which will trigger irq cores load balance, generally this value should greater-equal IRQCoresExpectedCpuUtil
	// +optional
	CPUUtilThreshold *int `json:"cpuUtilThreshold,omitempty"`
	// Threshold of cpu util gap between source core and dest core of irq affinity changing
	// +optional
	CPUUtilGapThreshold *int `json:"cpuUtilGapThreshold,omitempty"`
}

type IRQCoresAdjustConfig struct {
	// Minimum percent of (100 * irq cores/total(or socket) cores), valid value [0,100], default 2
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	// +optional
	PercentMin *int `json:"percentMin,omitempty"`

	// Maximum percent of (100 * irq cores/total(or socket) cores), valid value [0,100], default 30
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	// +optional
	PercentMax *int `json:"percentMax,omitempty"`

	// +optional
	IncConf *IRQCoresIncConfig `json:"incConf,omitempty"`
	// +optional
	DecConf *IRQCoresDecConfig `json:"decConf,omitempty"`
}

// IRQCoresIncConfig represents the configuration of increasing irq cores.
// When irq cores cpu util nearly full(e.g., greater-equal 85%), in order to reduce the impact time on the applications, it is necessary to immediately
// fall back to the balance-fair policy first, and later irq tuning manager will auto switch back to IrqCoresExclusive policy based on policies and conditions.
type IRQCoresIncConfig struct {
	// Interval of two successive irq cores increase MUST greater-equal this interval
	// +kubebuilder:validation:Minimum=1
	// +optional
	SuccessiveIncInterval *int `json:"successiveIncInterval,omitempty"`
	// When irq cores cpu util hit this thresh, then fallback to balance-fair policy
	// +optional
	FullThreshold *int `json:"fullThreshold,omitempty"`
	// +optional
	Thresholds *IRQCoresIncThresholds `json:"thresholds,omitempty"`
}

// IRQCoresIncThresholds represents the threshold of increasing irq cores.
// When irq cores average cpu util greater-equal IrqCoresAvgCpuUtilThresh, then increase irq cores,
// when there are one or more irq core's net load greater-equal IrqCoreNetOverloadThresholds, and failed to tune to irq load balance,
// then increase irq cores.
type IRQCoresIncThresholds struct {
	// Threshold of increasing irq cores, generally this thresh equal to or a litter greater-than IrqCoresExpectedCpuUtil
	// +optional
	AvgCPUUtilThreshold *int `json:"avgCPUUtilThreshold,omitempty"`
}

type IRQCoresDecConfig struct {
	// Interval of two successive irq cores decrease MUST greater-equal this interval
	// +optional
	SuccessiveDecInterval *int `json:"successiveDecInterval,omitempty"`
	// +optional
	PingPongAdjustInterval *int `json:"pingPongAdjustInterval,omitempty"`
	// +optional
	SinceLastBalanceInterval *int `json:"sinceLastBalanceInterval,omitempty"`
	// +optional
	DecCoresMaxEachTime *int `json:"decCoresMaxEachTime,omitempty"`
	// +optional
	Thresholds *IRQCoresDecThresholds `json:"thresholds,omitempty"`
}

// IRQCoresDecThresholds represents the threshold of decreasing irq cores.
// when irq cores average cpu util less-equal IrqCoresAvgCpuUtilThresh, then decrease irq cores.
type IRQCoresDecThresholds struct {
	// Threshold of decreasing irq cores, generally this thresh should be less-than IrqCoresExpectedCpuUtil
	// +optional
	AvgCPUUtilThreshold *int `json:"avgCPUUtilThreshold,omitempty"`
}

type IRQCoresExclusionConfig struct {
	// +optional
	Thresholds *IRQCoresExclusionThresholds `json:"thresholds,omitempty"`
	// Interval of successive enable/disable irq cores exclusion MUST >= SuccessiveSwitchInterval
	// +optional
	SuccessiveSwitchInterval *float64 `json:"successiveSwitchInterval,omitempty"`
}

type IRQCoresExclusionThresholds struct {
	// +optional
	EnableThresholds *EnableIRQCoresExclusionThresholds `json:"enableThresholds,omitempty"`
	// +optional
	DisableThresholds *DisableIRQCoresExclusionThresholds `json:"disableThresholds"`
}

// EnableIRQCoresExclusionThresholds represents the threshold of enabling irq cores exclusion.
// When successive count of nic's total PPS >= RxPPSThresh is greater-equal SuccessiveCount, then enable exclusion of this nic's irq cores.
type EnableIRQCoresExclusionThresholds struct {
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Format="uint64"
	// +optional
	RxPPSThreshold *uint64 `json:"rxPPSThreshold,omitempty"`
	// +optional
	SuccessiveCount *int `json:"successiveCount,omitempty"`
}

// DisableIRQCoresExclusionThresholds represents the threshold of disabling irq cores exclusion.
// When successive count of nic's total PPS <= RxPPSThresh is greater-equal SuccessiveCount, then disable exclusion of this nic's irq cores.
type DisableIRQCoresExclusionThresholds struct {
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Format="uint64"
	// +optional
	RxPPSThreshold *uint64 `json:"rxPPSThreshold,omitempty"`
	// +optional
	SuccessiveCount *int `json:"successiveCount,omitempty"`
}
