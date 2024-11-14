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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type GenericConfigSpec struct {
	// RevisionHistoryLimit is the maximum number of revisions that will
	// be maintained in the resource's revision history. The revision history
	// consists of all revisions not represented by a currently applied
	// Spec version. The default value is 3.
	// +kubebuilder:default:=3
	RevisionHistoryLimit int64 `json:"revisionHistoryLimit,omitempty"`

	// NodeLabelSelector select nodes to apply these configurations,
	// the priority and node label selector must be matched according
	// to KatalystCustomConfig.spec.nodeLabelSelectorAllowedKeyList,
	// otherwise it will not be synced.
	// +optional
	NodeLabelSelector string `json:"nodeLabelSelector,omitempty"`

	// Priority is used by one node matched by NodeLabelSelector of more
	// than one configuration, and the higher priority will be considered.
	// The priority only be supported when NodeLabelSelector set
	// +optional
	Priority int32 `json:"priority,omitempty"`

	// EphemeralSelector is a selector for temporary use only
	// +optional
	EphemeralSelector EphemeralSelector `json:"ephemeralSelector,omitempty"`

	// Indicates that the config is paused.
	// +optional
	Paused bool `json:"paused,omitempty"`

	// An update strategy to replace existing CustomNodeConfig configurations with new ones.
	// +optional
	UpdateStrategy ConfigUpdateStrategy `json:"updateStrategy,omitempty"`
}

// ConfigUpdateStrategy is a struct used to control the update strategy for a KatalystCustomConfig target.
type ConfigUpdateStrategy struct {
	// Type of config update. Only `RollingUpdate` is supported.
	// +optional
	Type ConfigUpdateStrategyType `json:"type,omitempty"`

	// Rolling update config params. Present only if type = "RollingUpdate".
	//---
	// TODO: Update this to follow our convention for oneOf, whatever we decide it
	// to be. Same as Deployment `strategy.rollingUpdate`.
	// See https://github.com/kubernetes/kubernetes/issues/35345
	// +optional
	RollingUpdate *RollingUpdateConfig `json:"rollingUpdate,omitempty"`
}

// +kubebuilder:validation:Enum=RollingUpdate
type ConfigUpdateStrategyType string

const (
	// Replace the old configs by new ones using rolling update i.e replace them on each node one after the other.
	RollingUpdateConfigStrategyType ConfigUpdateStrategyType = "RollingUpdate"
)

// Spec to control the desired behavior of config rolling update.
type RollingUpdateConfig struct {
	// The number or percentage of target CustomNodeConfigs to update to the current configuration.
	// For example: `100`` and `20%` are valid values.
	// +kubebuilder:validation:XIntOrString
	// +kubebuilder:validation:Pattern=`^(100|[1-9][0-9]?|0)%$`
	// +optional
	Canary *intstr.IntOrString `json:"canary,omitempty"`
}

type EphemeralSelector struct {
	// Specific nodes' name the configurations will be effected.
	// +optional
	NodeNames []string `json:"nodeNames,omitempty"`

	// define the duration this configuration will last from creationTimestamp.
	// must and only set when NodeNames already set
	LastDuration *metav1.Duration `json:"lastDuration,omitempty"`
}

type GenericConfigStatus struct {
	// The number of nodes that this config is targeting.
	// +optional
	TargetNodes int32 `json:"targetNodes"`

	// The number of nodes that this config is targeting and should be updated given the current strategy.
	// +optional
	CanaryNodes int32 `json:"canaryNodes"`

	// The number of target nodes that have been updated by this config.
	// +optional
	UpdatedTargetNodes int32 `json:"updatedTargetNodes"`

	// The number of nodes (including non-target nodes) that have been updated by this config.
	// +optional
	UpdatedNodes int32 `json:"updatedNodes"`

	// The hash of the current config observed by the kcc controller.
	// +optional
	CurrentHash string `json:"currentHash,omitempty"`

	// Count of hash collisions for this cr. The kcc controller
	// uses this field as a collision avoidance mechanism when it needs to
	// create the name for the newest ControllerRevision.
	// +optional
	CollisionCount *int32 `json:"collisionCount,omitempty"`

	// The most recent generation observed by the kcc controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Represents the latest available observations of a config's current state.
	// +optional
	Conditions []GenericConfigCondition `json:"conditions,omitempty"`
}

type GenericConfigCondition struct {
	// Type of config condition
	Type ConfigConditionType `json:"type"`
	// Status of the condition, one of True, False, Unknown.
	Status v1.ConditionStatus `json:"status"`
	// Last time the condition transit from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// reason is the reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty"`
	// message is a human-readable explanation containing details about the transition
	// +optional
	Message string `json:"message,omitempty"`
}

type ConfigConditionType string

const (
	// ConfigConditionTypeValid means this config whether is valid
	ConfigConditionTypeValid ConfigConditionType = "Valid"
)
