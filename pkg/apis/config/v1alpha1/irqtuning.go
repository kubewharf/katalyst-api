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
// +kubebuilder:resource:path=irqtuningconfigurations,shortName=irtc
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
	EnableTuner *bool `json:"enableTuner"`
	// TuningPolicy represents the interrupt tuning strategy. One of Balance, Exclusive, Auto.
	// +kubebuilder:default=balance
	TuningPolicy TuningPolicy `json:"tuningPolicy"`
	// TuningInterval is the interval of interrupt tuning.
	// +kubebuilder:default=5
	TuningInterval *int `json:"tuningInterval"`

	// EnableRPS indicates whether to enable the RPS function.
	// Only balance policy support enable rps.
	EnableRPS *bool `json:"enableRPS"`
	// EnableRPSCPUVSNicsQueue enable rps when (cpus count)/(nics queue count) greater than this config.
	EnableRPSCPUVSNicsQueue *float64 `json:"enableRPSCPUVSNicsQueue"`
	// NICAffinityPolicy represents the NICs's irqs affinity sockets policy.
	// One of CompleteMap, OverallBalance, PhysicalTopo.
	NICAffinityPolicy NICAffinityPolicy `json:"nicAffinityPolicy"`

	// ReniceKsoftirqd indicates whether to renice ksoftirqd process.
	ReniceKsoftirqd *bool `json:"reniceKsoftirqd"`
	// KsoftirqdNice is the nice value of ksoftirqd process.
	KsoftirqdNice *int `json:"ksoftirqdNice"`

	// CoresExpectedCPUUtil is the expected CPU utilization of cores.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	CoresExpectedCPUUtil *int `json:"coresExpectedCPUUtil"`

	// ThrouputClassSwitch describes the switch configuration for a throughput class.
	ThrouputClassSwitch *ThroughputClassSwitchConfig `json:"throughputClassSwitch"`
	// Threshold description for interrupting core network overLoad.
	CoreNetOverLoadThresh *IRQCoreNetOverloadThresholds `json:"coreNetOverLoadThresh"`
	// Describes the constraints of the balanced configuration.
	LoadBalance *IRQLoadBalanceConfig `json:"loadBalance"`
	// Configuration that requires interrupt core adjustment.
	CoresAdjust *IRQCoresAdjustConfig `json:"coresAdjust"`
	// Need to adjust to interrupt exclusive core requirements.
	CoresExclusion *IRQCoresExclusionConfig `json:"coresExclusion"`
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
	LowThroughputThresholds    *LowThroughputThresholds    `json:"lowThroughputThresholds"`
	NormalThroughputThresholds *NormalThroughputThresholds `json:"normalThroughputThresholds"`
}

type LowThroughputThresholds struct {
	RxPPSThresh     *uint64 `json:"rxPPSThresh"`
	SuccessiveCount *int    `json:"successiveCount"`
}

type NormalThroughputThresholds struct {
	RxPPSThresh     *uint64 `json:"rxPPSThresh"`
	SuccessiveCount *int    `json:"successiveCount"`
}

// IRQCoreNetOverloadThresholds represents the threshold of interrupt core network overload.
// When there are one or more irq core's ratio of softnet_stat 3rd col time_squeeze packets / 1st col processed packets
// greater-equal IrqCoreSoftNetTimeSqueezeRatio,
// then tring to tune irq load balance first, if failed to tune irq load balance, then increase irq cores.
type IRQCoreNetOverloadThresholds struct {
	// Ratio of softnet_stat 3rd col time_squeeze packets / softnet_stat 1st col processed packets
	SoftNetTimeSqueezeRatio *float64 `json:"softNetTimeSqueezeRatio"`
}

// IRQLoadBalanceConfig represents the configuration of interrupt load balance.
// When there are one or more irq core's cpu util greater-equal IrqCoreCpuUtilThresh or irq core's net load greater-equal IrqCoreNetOverloadThresholds,
// then try to tune irq load balance, that need to find at least one other irq core with relatively low cpu util, their cpu util gap MUST greater-equal IrqCoreCpuUtilGapThresh,
// if succeed to find irq cores with eligible cpu util, then start to tuning load balance,
// or increase irq cores immediately.
type IRQLoadBalanceConfig struct {
	// Interval of two successive irq load balance MUST greater-equal this interval
	SuccessiveTuningInterval *int                            `json:"successiveTuningInterval"`
	Thresholds               *IRQLoadBalanceTuningThresholds `json:"thresholds"`
	// Two successive tunes whose interval is less-equal this threshold will be considered as pingpong tunings
	PingPongIntervalThresh *int `json:"pingPongIntervalThresh"`
	// Ping pong count greater-equal this threshold will trigger increasing irq cores
	PingPongCountThresh *int `json:"pingPongCountThresh"`
	// Max number of irqs are permitted to be tuned from some irq cores to other cores in each time, allowed value {1, 2}
	// +kubebuilder:validation:Enum={1,2}
	IRQTunedNumMaxEachTime *int `json:"irqTunedNumMaxEachTime"`
	// Max number of irq cores whose affinitied irqs are permitted to tuned to other cores in each time, allowed value {1,2}
	// +kubebuilder:validation:Enum={1,2}
	IRQCoresTunedNumMaxEachTime *int `json:"irqCoresTunedNumMaxEachTime"`
}

type IRQLoadBalanceTuningThresholds struct {
	// IRQ core cpu util threshold, which will trigger irq cores load balance, generally this value should greater-equal IRQCoresExpectedCpuUtil
	CPUUtilThresh *int `json:"cpuUtilThresh"`
	// Threshold of cpu util gap between source core and dest core of irq affinity changing
	CPUUtilGapThresh *int `json:"cpuUtilGapThresh"`
}

type IRQCoresAdjustConfig struct {
	// Minimum percent of (100 * irq cores/total(or socket) cores), valid value [0,100], default 2
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	PercentMin *int `json:"percentMin"`

	// Maximum percent of (100 * irq cores/total(or socket) cores), valid value [0,100], default 30
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	PercentMax *int `json:"percentMax"`

	IncConf *IRQCoresIncConfig `json:"incConf"`
	DecConf *IRQCoresDecConfig `json:"decConf"`
}

// IRQCoresIncConfig represents the configuration of increasing irq cores.
// When irq cores cpu util nearly full(e.g., greater-equal 85%), in order to reduce the impact time on the applications, it is necessary to immediately
// fall back to the balance-fair policy first, and later irq tuning manager will auto switch back to IrqCoresExclusive policy based on policies and conditions.
type IRQCoresIncConfig struct {
	// Interval of two successive irq cores increase MUST greater-equal this interval
	// +kubebuilder:validation:Minimum=1
	SuccessiveIncInterval *int `json:"successiveIncInterval"`
	// When irq cores cpu util hit this thresh, then fallback to balance-fair policy
	FullThresh *int                   `json:"fullThresh"`
	Thresholds *IRQCoresIncThresholds `json:"thresholds"`
}

// IRQCoresIncThresholds represents the threshold of increasing irq cores.
// When irq cores average cpu util greater-equal IrqCoresAvgCpuUtilThresh, then increase irq cores,
// when there are one or more irq core's net load greater-equal IrqCoreNetOverloadThresholds, and failed to tune to irq load balance,
// then increase irq cores.
type IRQCoresIncThresholds struct {
	// Threshold of increasing irq cores, generally this thresh equal to or a litter greater-than IrqCoresExpectedCpuUtil
	AvgCPUUtilThresh *int `json:"avgCPUUtilThresh"`
}

type IRQCoresDecConfig struct {
	// Interval of two successive irq cores decrease MUST greater-equal this interval
	SuccessiveDecInterval    *int                   `json:"successiveDecInterval"`
	PingPongAdjustInterval   *int                   `json:"pingPongAdjustInterval"`
	SinceLastBalanceInterval *int                   `json:"sinceLastBalanceInterval"`
	DecCoresMaxEachTime      *int                   `json:"decCoresMaxEachTime"`
	Thresholds               *IRQCoresDecThresholds `json:"thresholds"`
}

// IRQCoresDecThresholds represents the threshold of decreasing irq cores.
// when irq cores average cpu util less-equal IrqCoresAvgCpuUtilThresh, then decrease irq cores.
type IRQCoresDecThresholds struct {
	// Threshold of decreasing irq cores, generally this thresh should be less-than IrqCoresExpectedCpuUtil
	AvgCPUUtilThresh *int `json:"avgCPUUtilThresh"`
}

type IRQCoresExclusionConfig struct {
	Thresholds *IRQCoresExclusionThresholds `json:"thresholds"`
	// Interval of successive enable/disable irq cores exclusion MUST >= SuccessiveSwitchInterval
	SuccessiveSwitchInterval *float64 `json:"successiveSwitchInterval"`
}

type IRQCoresExclusionThresholds struct {
	EnableThresholds  *EnableIRQCoresExclusionThresholds  `json:"enableThresholds"`
	DisableThresholds *DisableIRQCoresExclusionThresholds `json:"disableThresholds"`
}

// EnableIRQCoresExclusionThresholds represents the threshold of enabling irq cores exclusion.
// When successive count of nic's total PPS >= RxPPSThresh is greater-equal SuccessiveCount, then enable exclusion of this nic's irq cores.
type EnableIRQCoresExclusionThresholds struct {
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Format="uint64"
	RxPPSThresh     *uint64 `json:"rxPPSThresh"`
	SuccessiveCount *int    `json:"successiveCount"`
}

// DisableIRQCoresExclusionThresholds represents the threshold of disabling irq cores exclusion.
// When successive count of nic's total PPS <= RxPPSThresh is greater-equal SuccessiveCount, then disable exclusion of this nic's irq cores.
type DisableIRQCoresExclusionThresholds struct {
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Format="uint64"
	RxPPSThresh     *uint64 `json:"rxPPSThresh"`
	SuccessiveCount *int    `json:"successiveCount"`
}
