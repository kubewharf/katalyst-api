package v1alpha2

import (
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kubewharf/katalyst-api/pkg/apis/config/v1alpha1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:resource:shortName=ihpa

// IntelligentHorizontalPodAutoscaler captures information about a IHPA object
type IntelligentHorizontalPodAutoscaler struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the behavior of a IntelligentHorizontalPodAutoscaler.
	// +optional
	Spec IntelligentHorizontalPodAutoscalerSpec `json:"spec,omitempty"`

	// Status represents the current information about a IntelligentHorizontalPodAutoscaler.
	// +optional
	Status IntelligentHorizontalPodAutoscalerStatus `json:"status,omitempty"`
}

// IntelligentHorizontalPodAutoscalerSpec is the specification of the behavior of the autoscaler.
type IntelligentHorizontalPodAutoscalerSpec struct {
	// Autoscaler defines the overall scaling configuration.
	Autoscaler AutoscalerSpec `json:"autoscaler"`

	// ScaleStrategy defines whether to enable IHPA to scale workloads.
	// IHPA provides a preview mode. In preview mode, IHPA does not modify the number of replicas of the workload.
	// The default is 'Preview'.
	// +kubebuilder:default:=Preview
	ScaleStrategy ScaleStrategyType `json:"scaleStrategy,omitempty"`

	// AlgorithmConfig defines the algorithm configuration. If there is no configuration,
	// the default configuration will be used.
	AlgorithmConfig v1alpha1.AlgorithmConfig `json:"algorithmConfig,omitempty"`

	// TimeBounds supports dynamically adjusting HPA Replicas in different time periods,
	// providing capability similar to CronHPA.
	TimeBounds []TimeBound `json:"timeBounds,omitempty"`
}

// IntelligentHorizontalPodAutoscalerStatus describes the runtime state of the autoscaler.
type IntelligentHorizontalPodAutoscalerStatus struct {
	// lastScaleTime is the last time the HorizontalPodAutoscaler scaled the number of pods,
	// used by the autoscaler to control how often the number of pods is changed.
	// +optional
	LastScaleTime *metav1.Time `json:"lastScaleTime,omitempty"`

	// currentMetrics is the last read state of the metrics used by this autoscaler.
	// +listType=atomic
	// +optional
	CurrentMetrics []autoscalingv2.MetricStatus `json:"currentMetrics,omitempty"`

	// CurrentReplicas indicates the current number of replicas of the workload corresponding to
	// the HPA associated with the IHPA.
	CurrentReplicas int32 `json:"currentReplicas,omitempty"`

	// DesiredReplicas represents the number of replicas of the workload expected by IHPA.
	DesiredReplicas int32 `json:"desiredReplicas,omitempty"`
}

// AutoscalerSpec defines the associated workload, the metrics used, and the scaling behavior.
type AutoscalerSpec struct {
	// ScaleTargetRef defines the associated workload.
	ScaleTargetRef autoscalingv2.CrossVersionObjectReference `json:"scaleTargetRef"`

	// Behavior defines scaling actions, which are transparently transmitted to HPA and implemented by HPA Controller.
	// +optional
	Behavior *autoscalingv2.HorizontalPodAutoscalerBehavior `json:"behavior,omitempty"`

	// Metrics define the metrics used for workload scaling.
	// +optional
	Metrics []MetricSpec `json:"metrics,omitempty"`

	// MinReplicas and MaxReplicas are consistent with those in HPA.
	// +optional
	MinReplicas *int32 `json:"minReplicas,omitempty"`
	MaxReplicas int32  `json:"maxReplicas"`
}

// MetricSpec is used to define a single metric of HPA.
// Metric is preferred. Only one of Metric and CustomMetric will take effect.
// The metrics source of Metric is API Server.
// The metrics source of CustomMetric is the resource portrait.
type MetricSpec struct {
	Metric       *autoscalingv2.MetricSpec `json:"metric,omitempty"`
	CustomMetric *CustomMetricSpec         `json:"customMetric,omitempty"`
}

// CustomMetricSpec configures metrics derived from resource portraits.
type CustomMetricSpec struct {
	// Identify defines the name of the resource metric
	Identify corev1.ResourceName `json:"identify"`

	// Query defines the metrics query statement.
	// There are preset templates for resource portraits, so the query statement can be empty.
	Query string `json:"query,omitempty"`

	// Value represents the threshold, corresponding to the AverageValue in HPA.
	Value *resource.Quantity `json:"value"`
}

// TimeBound supports adjusting HPA max/min replicas based on the period and the time within the period.
type TimeBound struct {
	// Start and End are the time period.
	Start  metav1.Time `json:"start,omitempty"`
	End    metav1.Time `json:"end,omitempty"`
	Bounds []Bound     `json:"bounds,omitempty"`
}

// Bound defines the max/min replicas of HPA configured under the specified CronTab.
type Bound struct {
	CronTab string `json:"cronTab"`

	MinReplicas *int32 `json:"minReplicas,omitempty"`
	MaxReplicas *int32 `json:"maxReplicas,omitempty"`
}

// ScaleStrategyType is the strategy of IHPA to scale workloads.
type ScaleStrategyType string

const (
	Preview ScaleStrategyType = "Preview"
	Auto    ScaleStrategyType = "Auto"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IntelligentHorizontalPodAutoscalerList is a collection of IntelligentHorizontalPodAutoscaler objects.
type IntelligentHorizontalPodAutoscalerList struct {
	metav1.TypeMeta `json:",inline"`

	// Standard list metadata
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	// items is the list of IntelligentHorizontalPodAutoscaler
	Items []IntelligentHorizontalPodAutoscaler `json:"items"`
}

// VirtualWorkloadSpec defines the desired state of VirtualWorkload
type VirtualWorkloadSpec struct {
	Replicas int32 `json:"replicas"`
}

// VirtualWorkloadStatus defines the observed state of VirtualWorkload
type VirtualWorkloadStatus struct {
	Replicas int32  `json:"replicas"`
	Selector string `json:"selector"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:subresource:scale:specpath=.spec.replicas,statuspath=.status.replicas,selectorpath=.status.selector

// VirtualWorkload is the Schema for the virtualworkloads API
// VirtualWorkload is used to support IHPA's Preview mode, that is, by providing
// a virtual workload reference so that scaling will not affect the real workload.
type VirtualWorkload struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VirtualWorkloadSpec   `json:"spec,omitempty"`
	Status VirtualWorkloadStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VirtualWorkloadList contains a list of VirtualWorkload
type VirtualWorkloadList struct {
	metav1.TypeMeta `json:",inline"`

	// Standard list metadata
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	// items is the list of VirtualWorkload
	Items []VirtualWorkload `json:"items"`
}
