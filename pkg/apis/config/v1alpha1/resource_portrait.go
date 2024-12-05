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
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ResourcePortraitConfiguration defines the global configuration and matching rules of resource portraits.
// When the workload corresponding to the SPD meets the matching rules, the corresponding resource portrait
// config will be written to the SPD.
type ResourcePortraitConfiguration struct {
	// Config defines the algorithm and indicator configuration used for resource portrait.
	Config ResourcePortraitConfig `json:"config"`

	// Filter filters SPDs matching filter rules by namespaces and labels.
	// Namespaces and labels are derived from the workload referenced by the SPD.
	Filter ResourcePortraitFilter `json:"filter,omitempty"`
}

// ResourcePortraitConfig defines configuration detail of ResourcePortraitConfiguration
type ResourcePortraitConfig struct {
	// Source indicates the source of the configuration.
	// When this field is empty, the configuration is ignored.
	Source string `json:"source"`

	// Algorithm is the algorithm config used to calculate the resource portrait.
	AlgorithmConfig AlgorithmConfig `json:"algorithmConfig"`

	// Metrics contains preset metric query templates, such as CPU and memory metric queries.
	// +optional
	Metrics []string `json:"metrics,omitempty"`

	// CustomMetrics contains custom metric query, the key is the metric name, the value is the metric query.
	// +optional
	CustomMetrics map[string]string `json:"customMetrics,omitempty"`
}

type ResourcePortraitFilter struct {
	// Selector is used to filter workloads associated with SPD from the label dimension.
	// +optional
	Selector *metav1.LabelSelector `json:"selector,omitempty"`

	// Namespaces contains a list of namespaces used for filtering the target workload with spd.
	// The workload only needs to match a single namespace.
	// +optional
	Namespaces []string `json:"namespaces,omitempty"`
}

type AlgorithmConfig struct {
	// Method is the method used to calculate the resource portrait.
	Method string `json:"method"`

	// Params contains a list of parameters used by the algorithm.
	// +optional
	Params map[string]string `json:"params,omitempty"`

	// TimeWindow contains the time window used by the algorithm.
	TimeWindow TimeWindow `json:"timeWindow"`

	// ResyncPeriod contains the resync period used by the algorithm.
	ResyncPeriod time.Duration `json:"resyncPeriod"`
}

// +kubebuilder:validation:Enum=avg;max

// Aggregator is used in conjunction with TimeWindow to represent the conversion of time series data
// from input to output, which is the processing logic inside the resource portrait.
type Aggregator string

const (
	Avg Aggregator = "avg"
	Max Aggregator = "max"
)

type TimeWindow struct {
	// Input is the time interval used when querying the time series data source, with the unit being seconds.
	Input int `json:"input"`

	// HistorySteps is the number of steps to use as the input for the algorithm.
	HistorySteps int `json:"historySteps"`

	// Aggregator is the aggregator used to aggregate the time series data group.
	Aggregator Aggregator `json:"aggregator"`

	// Output is the grouping interval of portrait data, which groups portrait data at an integer multiple of Input.
	Output int `json:"output"`

	// PredictionSteps is the number of steps to predict the future value.
	PredictionSteps int `json:"predictionSteps"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GlobalResourcePortraitConfiguration is the global configuration of the resource portrait plug-in,
// including resource portrait configuration and workload matching rules.
type GlobalResourcePortraitConfiguration struct {
	metav1.TypeMeta `json:",inline"`

	// Configs are global resource portrait configuration
	Configs []ResourcePortraitConfiguration `json:"configs,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ResourcePortraitIndicators define multiple resource portrait configurations that a single SPD may contain.
type ResourcePortraitIndicators struct {
	metav1.TypeMeta `json:",inline"`

	// Configs are spd's resource portrait configuration
	Configs []ResourcePortraitConfig `json:"configs,omitempty"`
}
