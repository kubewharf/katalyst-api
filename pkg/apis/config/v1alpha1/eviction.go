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
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=evictionconfigurations,shortName=ec
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="AGE",type=date,JSONPath=.metadata.creationTimestamp
// +kubebuilder:printcolumn:name="SELECTOR",type=string,JSONPath=".spec.nodeLabelSelector"
// +kubebuilder:printcolumn:name="NODES",type=string,JSONPath=".spec.ephemeralSelector.nodeNames"
// +kubebuilder:printcolumn:name="DURATION",type=string,JSONPath=".spec.ephemeralSelector.lastDuration"
// +kubebuilder:printcolumn:name="VALID",type=string,JSONPath=".status.conditions[?(@.type==\"Valid\")].status"
// +kubebuilder:printcolumn:name="REASON",type=string,JSONPath=".status.conditions[?(@.type==\"Valid\")].reason"
// +kubebuilder:printcolumn:name="MESSAGE",type=string,JSONPath=".status.conditions[?(@.type==\"Valid\")].message"

// EvictionConfiguration is the Schema for the configuration API used by eviction
type EvictionConfiguration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EvictionConfigurationSpec `json:"spec,omitempty"`
	Status GenericConfigStatus       `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

// EvictionConfigurationList contains a list of EvictionConfiguration
type EvictionConfigurationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EvictionConfiguration `json:"items"`
}

// EvictionConfigurationSpec defines the desired state of EvictionConfiguration
type EvictionConfigurationSpec struct {
	GenericConfigSpec `json:",inline"`

	// Config is custom field for eviction configuration
	// all configuration crd should contain a field named with `config`
	Config EvictionConfig `json:"config,omitempty"`
}

type EvictionConfig struct {
	// EvictionPluginsConfig is the config for all eviction plugins
	EvictionPluginsConfig EvictionPluginsConfig `json:"evictionPluginsConfig,omitempty"`
}

type EvictionPluginsConfig struct {
	// ReclaimedResourcesEvictionPluginConfig is the config for reclaimed resources eviction plugin
	ReclaimedResourcesEvictionPluginConfig ReclaimedResourcesEvictionPluginConfig `json:"reclaimedResourcesEvictionPluginConfig,omitempty" json:"reclaimedResourcesEvictionPluginConfig"`

	// MemoryEvictionPluginConfig is the config for memory eviction plugin
	MemoryEvictionPluginConfig MemoryEvictionPluginConfig `json:"memoryEvictionPluginConfig,omitempty" json:"memoryEvictionPluginConfig"`
}

type ReclaimedResourcesEvictionPluginConfig struct {
	// EvictionThreshold eviction threshold rate for reclaimed resources
	EvictionThreshold map[v1.ResourceName]float64 `json:"evictionThreshold,omitempty"`
}

type MemoryEvictionPluginConfig struct {
	// EnableNumaLevelDetection is whether to enable numa-level detection
	EnableNumaLevelDetection *bool `json:"enableNumaLevelDetection,omitempty"`

	// EnableSystemLevelDetection is whether to enable system-level detection
	EnableSystemLevelDetection *bool `json:"enableSystemLevelDetection,omitempty"`

	// NumaFreeBelowWatermarkTimesThreshold is the threshold for the number of times NUMA's free memory falls below the watermark
	NumaFreeBelowWatermarkTimesThreshold *int `json:"numaFreeBelowWatermarkTimesThreshold,omitempty"`

	// NumaFreeBelowWatermarkTimesThreshold is the threshold for the rate of kswapd reclaiming rate
	SystemKswapdRateThreshold *int `json:"systemKswapdRateThreshold,omitempty"`

	// SystemKswapdRateExceedCountThreshold is the threshold for the number of times the kswapd reclaiming rate exceeds the threshold
	SystemKswapdRateExceedTimesThreshold *int `json:"systemKswapdRateExceedTimesThreshold,omitempty"`

	// NumaEvictionRankingMetrics is the metrics used to rank pods for eviction at the NUMA level
	NumaEvictionRankingMetrics []string `json:"numaEvictionRankingMetrics,omitempty"`

	// SystemEvictionRankingMetrics is the metrics used to rank pods for eviction at the system level
	SystemEvictionRankingMetrics []string `json:"systemEvictionRankingMetrics,omitempty"`

	// GracePeriod is the grace period of memory eviction
	GracePeriod *int64 `json:"gracePeriod,omitempty"`

	// EnableRSSOveruseDetection is whether to enable pod-level rss overuse detection
	EnableRssOveruseDetection *bool `json:"enableRssOveruseDetection,omitempty"`

	// RssOveruseRateThreshold is the threshold for the rate of rss
	RssOveruseRateThreshold *float64 `json:"rssOveruseRateThreshold,omitempty"`
}
