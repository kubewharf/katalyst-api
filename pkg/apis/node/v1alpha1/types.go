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

	"github.com/kubewharf/katalyst-api/pkg/consts"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:scope=Cluster,shortName=kcnr
// +kubebuilder:subresource:status

// CustomNodeResource captures information about a custom defined node resource, mainly focus on static attributes and resources
// CustomNodeResource objects are non-namespaced.
type CustomNodeResource struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the behavior of a CustomNodeResource.
	// +optional
	Spec CustomNodeResourceSpec `json:"spec,omitempty"`

	// Status represents the current information about a CustomNodeResource.
	// This data may not be up-to-date.
	// +optional
	Status CustomNodeResourceStatus `json:"status,omitempty"`
}

type CustomNodeResourceSpec struct {
	// +optional
	NodeResourceProperties []*Property `json:"nodeResourceProperties,omitempty"`

	// Taints customized taint for katalyst, which affect pod scheduling
	// based on their QoS levels and the specified taint's QoS level.
	// +optional
	Taints []Taint `json:"taints,omitempty"`
}

// Taint wraps standard Kubernetes Taint with QoSLevel.
type Taint struct {
	// Taint is standard Kubernetes Taint
	v1.Taint `json:",inline"`

	// QoSLevel specifies the QoS level of pods that this taint applies to.
	// +kubebuilder:validation:Enum=reclaimed_cores;shared_cores;dedicated_cores;system_cores
	QoSLevel consts.QoSLevel `json:"qosLevel"`
}

type Property struct {
	// property name
	PropertyName string `json:"propertyName"`

	// values of the specific property
	// +optional
	PropertyValues []string `json:"propertyValues,omitempty"`

	// values of the quantity-types property
	// +optional
	PropertyQuantity *resource.Quantity `json:"propertyQuantity,omitempty"`
}

type CustomNodeResourceStatus struct {
	// Resources defines the numeric quantities in this node; for instance reclaimed resources for this node
	// +optional
	Resources Resources `json:"resources"`

	// +optional
	TopologyZone []*TopologyZone `json:"topologyZone,omitempty"`

	// TopologyPolicy indicates placement policy for scheduler or other centralized components to follow.
	// this policy (including topology scope) is defined in topology-manager, katalyst is
	// responsible to parse the policy, and transform to TopologyPolicy here.
	// +kubebuilder:default:=none
	TopologyPolicy TopologyPolicy `json:"topologyPolicy"`

	// Conditions is an array of current observed cnr conditions.
	// +optional
	Conditions []CNRCondition `json:"conditions,omitempty"`

	// NodeMetricStatus report node real-time metrics
	// +optional
	NodeMetricStatus *NodeMetricStatus `json:"nodeMetricStatus,omitempty"`
}

type TopologyPolicy string

const (
	// TopologyPolicyNone policy is the default policy and does not perform any topology alignment.
	TopologyPolicyNone TopologyPolicy = "None"

	// TopologyPolicySingleNUMANodeContainerLevel represents single-numa-node policy and container level.
	TopologyPolicySingleNUMANodeContainerLevel TopologyPolicy = "SingleNUMANodeContainerLevel"

	// TopologyPolicySingleNUMANodePodLevel represents single-numa-node policy and pod level.
	TopologyPolicySingleNUMANodePodLevel TopologyPolicy = "SingleNUMANodePodLevel"

	// TopologyPolicyRestrictedContainerLevel represents restricted policy and container level.
	TopologyPolicyRestrictedContainerLevel TopologyPolicy = "RestrictedContainerLevel"

	// TopologyPolicyRestrictedPodLevel represents restricted policy and pod level.
	TopologyPolicyRestrictedPodLevel TopologyPolicy = "RestrictedPodLevel"

	// TopologyPolicyBestEffortContainerLevel represents best-effort policy and container level.
	TopologyPolicyBestEffortContainerLevel TopologyPolicy = "BestEffortContainerLevel"

	// TopologyPolicyBestEffortPodLevel represents best-effort policy and pod level.
	TopologyPolicyBestEffortPodLevel TopologyPolicy = "BestEffortPodLevel"

	// TopologyPolicyNumericContainerLevel represents numeric policy and container level.
	TopologyPolicyNumericContainerLevel TopologyPolicy = "NumericContainerLevel"

	// TopologyPolicyNumericPodLevel represents numeric policy and pod level.
	TopologyPolicyNumericPodLevel TopologyPolicy = "NumericPodLevel"
)

// CNRCondition contains condition information for a cnr.
type CNRCondition struct {
	// Type is the type of the condition.
	Type CNRConditionType `json:"type"`
	// Status of the condition, one of True, False, Unknown.
	Status v1.ConditionStatus `json:"status" `
	// Last time we got an update on a given condition.
	// +optional
	LastHeartbeatTime metav1.Time `json:"lastHeartbeatTime,omitempty"`
	// (brief) reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty"`
	// Human-readable message indicating details about last transition.
	// +optional
	Message string `json:"message,omitempty"`
}

type CNRConditionType string

const (
	CNRAgentReady    CNRConditionType = "AgentReady"
	CNRAgentNotFound CNRConditionType = "AgentNotFound"
)

type TopologyZone struct {
	// Type represents which kind of resource this TopologyZone is for;
	// for instance, Socket, Numa, GPU, NIC, Disk and so on.
	Type TopologyType `json:"type"`

	// Name represents the name for the given type for resource; for instance,
	// - disk-for-log, disk-for-storage may have different usage or attributes, so we
	//   need separate structure to distinguish them.
	Name string `json:"name"`

	// Resources defines the numeric quantities in this TopologyZone; for instance,
	// - a TopologyZone with type TopologyTypeGPU may have both gpu and gpu-memory
	// - a TopologyZone with type TopologyTypeNIC may have both ingress and egress bandwidth
	// +optional
	Resources Resources `json:"resources"`

	// +optional
	// +patchMergeKey=name
	// +patchStrategy=merge
	Attributes []Attribute `json:"attributes,omitempty" patchStrategy:"merge" patchMergeKey:"name"`

	// +optional
	// +patchMergeKey=consumer
	// +patchStrategy=merge
	Allocations []*Allocation `json:"allocations,omitempty" patchStrategy:"merge" patchMergeKey:"consumer"`

	// Children represents the ownerships between multiple TopologyZone; for instance,
	// - a TopologyZone with type TopologyTypeSocket may have multiple childed TopologyZone
	//   with type TopologyTypeNuma to reflect the physical connections for a node
	// - a TopologyZone with type `nic` may have multiple childed TopologyZone with type `vf`
	//   to reflect the `physical and virtual` relations between devices
	// todo: in order to bypass the lacked functionality of recursive structure definition,
	//  we need to skip validation of this field for now; will re-add this validation logic
	//  if the community supports $ref, for more information, please
	//  refer to https://github.com/kubernetes/kubernetes/issues/62872
	// +optional
	// +kubebuilder:validation:Schemaless
	// +kubebuilder:pruning:PreserveUnknownFields
	Children []*TopologyZone `json:"children,omitempty"`

	// Siblings represents the relationship between TopologyZones at the same level; for instance,
	// the distance between NUMA nodes.
	// +optional
	// +patchMergeKey=name
	// +patchStrategy=merge
	Siblings []Sibling `json:"siblings,omitempty" patchStrategy:"merge" patchMergeKey:"name"`
}

type TopologyType string

const (
	// TopologyTypeSocket indicates socket-level topology
	TopologyTypeSocket TopologyType = "Socket"

	// TopologyTypeNuma indicates numa-level topology
	TopologyTypeNuma TopologyType = "Numa"

	// TopologyTypeCacheGroup indicates cache-group-level topology
	TopologyTypeCacheGroup TopologyType = "CacheGroup"

	// TopologyTypeGPU indicates a zone for gpu device
	TopologyTypeGPU TopologyType = "GPU"

	// TopologyTypeNIC indicates a zone for network device
	TopologyTypeNIC TopologyType = "NIC"
)

type Resources struct {
	// +optional
	Allocatable *v1.ResourceList `json:"allocatable,omitempty"`

	// +optional
	Capacity *v1.ResourceList `json:"capacity,omitempty"`

	// ResourcePackages defines compute packages available on the node/numa.
	// Concept:
	//   - ResourcePackages are node/numa-oriented.
	//   - A ResourcePackage represents a subdivision of the total node/numa resources into
	//     standardized units. Each unit may define one or more resource dimensions
	//     (e.g., CPU, memory, disk, network).
	//   - Pods associated with a package must consume resources following the
	//     same ratio. If a Pod is not bound to a package, it falls back to
	//     the "default" package (or its variants: "default-1", "default-2").
	// Difference vs ResourcePools:
	//   - ResourcePackages: split all node/numa resources into standard shapes (abstracting
	//     physical resources into units).
	//   - ResourcePools: reserve/limit resources for a particular workload type.
	// +optional
	// +listMapKey=packageName
	// +listType=map
	ResourcePackages []ResourcePackage `json:"resourcePackages,omitempty"`

	// ResourcePools defines pools of resources reserved for specific workloads.
	// Concept:
	//   - ResourcePools are workload-oriented.
	//   - They allow a workload type (e.g. GPU jobs, latency-sensitive tasks) to
	//     reserve a guaranteed amount of resources (via MinAllocatable) while
	//     also optionally borrowing up to a maximum limit (via MaxAllocatable).
	// Difference vs ResourcePackages:
	//   - ResourcePools are for *workload reservations* (claiming resources for
	//     certain job types).
	//   - ResourcePackages are for *node subdivisions* (splitting all node resources into
	//     standard allocation units).
	// +optional
	// +listMapKey=poolName
	// +listType=map
	ResourcePools []ResourcePool `json:"resourcePools,omitempty"`
}

// Attribute records the resource-specified info with name-value pairs
type Attribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Allocation struct {
	Consumer string `json:"consumer"`
	// +optional
	Requests *v1.ResourceList `json:"requests,omitempty"`
}

// Sibling describes the relationship between two Zones.
type Sibling struct {
	// Type represents the type of this Sibling.
	// For instance, Socket, Numa, GPU, NIC, Disk and so on.
	Type TopologyType `json:"type"`

	// Name represents the name of this Sibling.
	Name string `json:"name"`

	// Attributes are the attributes of the relationship between two Zones.
	// For instance, the distance between tow NUMA nodes, the connection type between two GPUs, etc.
	// +patchMergeKey=name
	// +patchStrategy=merge
	Attributes []Attribute `json:"attributes,omitempty" patchStrategy:"merge" patchMergeKey:"name"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CustomNodeResourceList is a collection of CustomNodeResource objects.
type CustomNodeResourceList struct {
	metav1.TypeMeta `json:",inline"`

	// Standard list metadata
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	// items is the list of CNRs
	Items []CustomNodeResource `json:"items"`
}

// NodeMetricStatus defines the observed state of NodeMetric
type NodeMetricStatus struct {
	// UpdateTime is the last time this NodeMetricStatus was updated.
	UpdateTime metav1.Time `json:"updateTime"`

	// NodeMetric contains the metrics for this node.
	NodeMetric *NodeMetricInfo `json:"nodeMetric,omitempty"`

	// GroupMetric contains the metrics aggregated by QoS level groups
	GroupMetric []GroupMetricInfo `json:"groupMetric,omitempty"`
}

type ResourceMetric struct {
	CPU    *resource.Quantity `json:"cpu,omitempty"`
	Memory *resource.Quantity `json:"memory,omitempty"`
}

type ResourceUsage struct {
	// NUMAUsage contains the real-time resource usage for each NUMA
	NUMAUsage []NUMAMetricInfo `json:"numaUsage,omitempty"`

	// GenericUsage contains the real-time resource usage
	GenericUsage *ResourceMetric `json:"genericUsage,omitempty"`
}

type GroupMetricInfo struct {
	// +kubebuilder:validation:Enum=reclaimed_cores;shared_cores;dedicated_cores;system_cores
	QoSLevel      string `json:"QoSLevel"`
	ResourceUsage `json:",inline"`
	// PodList indicates the pods belongs to this qos group, in format of {namespace}/{name}.
	// Pods that have been scheduled but are not listed in the PodList need to be estimated by the scheduler.
	PodList []string `json:"podList,omitempty"`
}

type NodeMetricInfo struct {
	ResourceUsage `json:",inline"`
}

type NUMAMetricInfo struct {
	NUMAId int `json:"numaId"`
	// Usage contains the real-time resource usage for this NUMA node
	Usage *ResourceMetric `json:"usage"`
}

// ResourcePackage represents a single compute package definition.
// Concept:
//   - A ResourcePackage subdivides the node/numa’s total resources into standardized units.
//   - The most common naming convention is based on CPU:Memory ratio.
//   - Example: "x2" means 1 core : 2 GiB memory, "x8" means 1 core : 8 GiB memory
//   - In the future, packages may also define additional resource dimensions
//     (e.g., local SSDs, network bandwidth, GPUs).
//
// Behavior:
//   - Pods associated with this package must consume resources following the
//     shape defined in Allocatable.
//   - A special "default" resource package (or "default-N" variants) represents all
//     remaining resources not explicitly assigned to a named resource package.
//
// Data source:
//   - ResourcePackages are derived from metrics reported in a NodeProfileDescriptor CRD.
//   - These metrics are aggregated (e.g., across NUMA nodes) to compute the
//     Allocatable resources for each package.
type ResourcePackage struct {
	// PackageName is the identifier for this package, e.g. "x2", "x8".
	// Rules:
	// - Names like "default" or "default-N" (N = integer) are reserved identifiers
	//   for the default package(s), which represent resources not explicitly
	//   subdivided into other named packages.
	// - Other names are user-defined and may follow any convention
	//   (e.g., "x2", "x8").
	PackageName string `json:"packageName"`

	// Allocatable defines the total resources available for this package.
	// Keys usually include "cpu" and "memory" (e.g. cpu: "64", memory: "128Gi").
	Allocatable *v1.ResourceList `json:"allocatable,omitempty"`
}

// ResourcePool represents a pool of resources reserved for a specific workload type.
// Concept:
//   - A ResourcePool is workload-oriented. It defines a reserved or guaranteed
//     set of resources and the possible upper bound for a workload (e.g., GPU workloads, latency-sensitive services).
//   - Unlike ResourcePackages (which divide total node/numa resources into fixed CPU:Memory bundles),
//     ResourcePools allow flexible reservation and borrowing of resources.
//
// Data source:
//   - ResourcePools are derived from metrics reported in a NodeProfileDescriptor CRD.
//   - These metrics are mapped into MinAllocatable/MaxAllocatable values.
type ResourcePool struct {
	// PoolName is the unique identifier of the resource pool.
	PoolName string `json:"poolName"`

	// MinAllocatable defines the minimum amount of resources *guaranteed* or reserved
	// for this pool. This is the lower bound that the pool can always access.
	// Interpretation:
	//   - Acts like a reservation or guaranteed quota.
	//   - Ensures workloads in this pool always get at least these resources.
	MinAllocatable *v1.ResourceList `json:"minAllocatable,omitempty"`

	// MaxAllocatable defines the maximum amount of resources this pool
	// can consume, including any resources it may opportunistically
	// borrowed from other pools.
	// Interpretation:
	//   - Acts like an upper bound / quota ceiling.
	//   - Workloads in this pool cannot exceed this allocation, even if more
	//     resources are available on the node.
	//   - Supports resource sharing between pools, but with clear limits.
	MaxAllocatable *v1.ResourceList `json:"maxAllocatable,omitempty"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:scope=Cluster,shortName=npd
// +kubebuilder:subresource:status

// NodeProfileDescriptor captures information about node, such as node-related metrics
// NodeProfileDescriptor objects are non-namespaced.
type NodeProfileDescriptor struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the behavior of a NodeProfileDescriptor.
	// +optional
	Spec NodeProfileDescriptorSpec `json:"spec,omitempty"`

	// Status represents the current information about a NodeProfileDescriptor.
	// This data may not be up-to-date.
	// +optional
	Status NodeProfileDescriptorStatus `json:"status,omitempty"`
}

type NodeProfileDescriptorSpec struct {
}

type NodeProfileDescriptorStatus struct {
	// NodeMetrics contains the node-related metrics
	// +optional
	NodeMetrics []ScopedNodeMetrics `json:"nodeMetrics,omitempty"`

	// PodMetrics contains the pod-related metrics
	// +optional
	PodMetrics []ScopedPodMetrics `json:"podMetrics,omitempty"`
}

type ScopedNodeMetrics struct {
	// +optional
	Scope string `json:"scope,omitempty"`

	// +optional
	Metrics []MetricValue `json:"metrics,omitempty"`
}

type ScopedPodMetrics struct {
	// +optional
	Scope string `json:"scope,omitempty"`

	// +optional
	PodMetrics []PodMetric `json:"podMetrics,omitempty"`
}

type PodMetric struct {
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// +optional
	Name string `json:"name,omitempty"`

	// +optional
	Metrics []MetricValue `json:"metrics,omitempty"`
}

type Aggregator string

const (
	AggregatorAvg   Aggregator = "avg"
	AggregatorMax   Aggregator = "max"
	AggregatorMin   Aggregator = "min"
	AggregatorCount Aggregator = "count"
	AggregatorP99   Aggregator = "p99"
	AggregatorP95   Aggregator = "p95"
	AggregatorP90   Aggregator = "p90"
)

type MetricValue struct {
	// the name of the metric
	// +optional
	MetricName string `json:"metricName,omitempty"`

	// a set of labels that identify a single time series for the metric
	// +optional
	MetricLabels map[string]string `json:"metricLabels,omitempty"`

	// indicates the time at which the metrics were produced
	// +optional
	Timestamp metav1.Time `json:"timestamp,omitempty"`

	// the aggregator of the metric
	// +optional
	Aggregator *Aggregator `json:"aggregator,omitempty"`

	// indicates the window ([Timestamp-Window, Timestamp]) from
	// which these metrics were calculated, when returning rate
	// metrics calculated from cumulative metrics (or zero for
	// non-calculated instantaneous metrics).
	// +optional
	Window *metav1.Duration `json:"window,omitempty"`

	// the value of the metric
	// +optional
	Value resource.Quantity `json:"value,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NodeProfileDescriptorList is a collection of NodeProfileDescriptor objects.
type NodeProfileDescriptorList struct {
	metav1.TypeMeta `json:",inline"`

	// Standard list metadata
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	// items is the list of CNRs
	Items []NodeProfileDescriptor `json:"items"`
}
