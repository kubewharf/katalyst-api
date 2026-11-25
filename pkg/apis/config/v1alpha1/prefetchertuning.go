/*
Copyright 2026 The Katalyst Authors.

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

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=prefetchertunings,scope=Cluster
// +kubebuilder:subresource:status

// PrefetcherTuning selects nodes and defines hardware prefetcher rules for their workloads.
type PrefetcherTuning struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PrefetcherTuningSpec `json:"spec"`
	Status GenericConfigStatus  `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

// PrefetcherTuningList contains a list of PrefetcherTuning objects.
type PrefetcherTuningList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PrefetcherTuning `json:"items"`
}

type PrefetcherTuningSpec struct {
	GenericConfigSpec `json:",inline"`

	Config PrefetcherTuningConfig `json:"config"`
}

type PrefetcherTuningConfig struct {
	// EnableTuner enables tuning on selected nodes. If omitted or false, nodes restore their managed MSRs.
	// +optional
	EnableTuner bool `json:"enableTuner,omitempty"`

	// Rules defines default and strategy-specific rules. An empty list revokes all rules.
	// Matching rules must agree on Disable; conflicting rules restore All On and report an error.
	// +listType=map
	// +listMapKey=labelSelector
	// +listMapKey=cpu_codename
	// +listMapKey=strategyName
	Rules []PrefetcherTuningRule `json:"rules"`
}

type PrefetcherTuningRule struct {
	// StrategyName selects an SG tuning strategy. Empty means the default when SG has no prefetcher strategy.
	// +optional
	// +kubebuilder:default=""
	// +kubebuilder:validation:Pattern=`^$|^prefetcher_tuning_.+$`
	StrategyName string `json:"strategyName,omitempty"`

	// LabelSelector matches Pod labels using Kubernetes label selector syntax. Empty selectors are invalid.
	// +kubebuilder:validation:MinLength=1
	LabelSelector string `json:"labelSelector"`
	// CPUCodeName matches the cpu_codename reported by Malachite.
	// +kubebuilder:validation:Enum=AMD_K19Zen3;AMD_K19Zen4;AMD_K1AZen5;Intel_IceLake;Intel_SapphireRapids;Intel_GraniteRapids;Intel_EmeraldRapids
	CPUCodeName string `json:"cpu_codename"`
	// Disable lists the prefetchers to turn off. An empty list enables all managed prefetchers.
	// +listType=set
	Disable []Prefetcher `json:"disable"`
}

// Prefetcher identifies a managed hardware prefetcher.
// +kubebuilder:validation:Enum=l1_stream;l1_stride;l1_region;l2_stream;up_down
type Prefetcher string
