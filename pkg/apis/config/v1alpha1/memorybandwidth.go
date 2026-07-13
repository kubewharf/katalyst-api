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
// +kubebuilder:resource:path=memorybandwidthconfigurations,shortName=mbc
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

// MemoryBandwidthConfiguration is the Schema for the configuration API used by the
// QRM memory bandwidth (MB) plugin. It exposes the per-node knobs that were
// previously only settable via katalyst-agent command line flags.
type MemoryBandwidthConfiguration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MemoryBandwidthConfigurationSpec `json:"spec,omitempty"`
	Status GenericConfigStatus              `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

// MemoryBandwidthConfigurationList contains a list of MemoryBandwidthConfiguration
type MemoryBandwidthConfigurationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MemoryBandwidthConfiguration `json:"items"`
}

type MemoryBandwidthConfigurationSpec struct {
	GenericConfigSpec `json:",inline"`

	// Config is custom field for memory bandwidth configuration
	Config MemoryBandwidthConfig `json:"config"`
}

// MemoryBandwidthConfig holds the tunables for the memory bandwidth suppression
// plugin. Every field is optional; a nil/empty field means "keep the agent's
// current (flag or built-in) default".
type MemoryBandwidthConfig struct {
	// CCDMinMB is the minimum memory bandwidth (in MB/s) that any control group
	// is guaranteed per CCD; it is the throttling floor so a CCD is never
	// starved to a stall. Maps to --mb-ccd-min.
	// +kubebuilder:validation:Minimum=0
	// +optional
	CCDMinMB *int `json:"ccdMinMB,omitempty"`

	// CCDMaxMB is the maximum memory bandwidth (in MB/s) allowed per CCD; it is
	// the saturation point of bandwidth control and doubles as the "fully open"
	// value for no-throttle groups. Maps to --mb-ccd-max.
	// +kubebuilder:validation:Minimum=0
	// +optional
	CCDMaxMB *int `json:"ccdMaxMB,omitempty"`

	// CapLimitPercent is the coefficient (percent) applied when translating a
	// desired outgoing target into the value written to the resctrl schemata,
	// since the set value tends to yield slightly more traffic than requested.
	// Maps to --mb-cap-limit-percent.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	// +optional
	CapLimitPercent *int `json:"capLimitPercent,omitempty"`

	// GroupAwareCapacityPercentage customizes the effective domain capacity per
	// control group, expressed as a percentage of the measured domain capacity.
	// Only values in the open interval (0,100) take effect. Keyed by resctrl
	// group name (e.g. "share-50"). Maps to --mb-group-aware-capacity-percentage.
	// +optional
	GroupAwareCapacityPercentage map[string]int `json:"groupAwareCapacityPercentage,omitempty"`

	// NoThrottleGroups is the list of resctrl control groups that must never be
	// throttled regardless of pressure (the built-in root "/" is always exempt).
	// Maps to --mb-no-throttle-groups.
	// +optional
	NoThrottleGroups []string `json:"noThrottleGroups,omitempty"`

	// CCDCapGroups sets, per control group, the target actual memory bandwidth
	// (in MB/s) per CCD that the P-Controller drives toward. Keyed by resctrl
	// group name (e.g. "dedicated=20000"). Requires CCDCapKp > 0 to take effect.
	// Maps to --mb-ccd-cap-groups.
	// +optional
	CCDCapGroups map[string]int `json:"ccdCapGroups,omitempty"`

	// CCDCapKp is the proportional gain of the P-Controller that governs the
	// per-CCD cap for CCDCapGroups. The P-Controller layer is only enabled when
	// this is greater than 0. Maps to --mb-ccd-cap-kp.
	// +optional
	CCDCapKp *float64 `json:"ccdCapKp,omitempty"`

	// ResetResctrlOnly, when true, makes the plugin only reset the resctrl
	// schemata files back to their default (unrestricted) state and then quit,
	// without actually managing bandwidth. Maps to --mb-reset-resctrl-only.
	// +optional
	ResetResctrlOnly *bool `json:"resetResctrlOnly,omitempty"`

	// LocalIsVictim selects the special AMD resctrl counter semantics where
	// mbm_total_bytes holds all reads and mbm_local_bytes holds victim writes,
	// so the reader reconstructs total bandwidth accordingly. Maps to
	// --mb-local-is-victim.
	// +optional
	LocalIsVictim *bool `json:"localIsVictim,omitempty"`

	// ExtraGroupPriorities registers priority weights for resctrl control groups
	// beyond the built-in ones, controlling throttling order (higher weight is
	// protected longer). Keyed by group name. Maps to --mb-extra-group-priorities.
	// +optional
	ExtraGroupPriorities map[string]int `json:"extraGroupPriorities,omitempty"`
}
