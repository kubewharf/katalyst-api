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

// NodeOvercommitConfigSpec is a description of a NodeOvercommitConfig
type NodeOvercommitConfigSpec struct {
	// NodeOvercommitSelectorVal is the value of node label selector with key consts.NodeOvercommitSelectorKey,
	// it decides whether to update Nodes if the Node matches the selector 'consts.NodeOvercommitSelectorKey=NodeOvercommitSelectorVal'
	// +optional
	NodeOvercommitSelectorVal string `json:"nodeOvercommitSelectorVal,omitempty"`

	// ResourceOvercommitRatio describes the resource overcommit ratio that needs to inject into Node.Annotations
	// cpu,memory are supported.
	// +optional
	ResourceOvercommitRatio map[v1.ResourceName]string `json:"resourceOvercommitRatio,omitempty"`

	// TimeBounds supports dynamically adjusting ResourceOvercommitRatio
	// +optional
	TimeBounds []TimeBound `json:"timeBounds,omitempty"`

	// If CronConsistency is true, controller will ensure ResourceOvercommitRatio is consistent with the configuration of the TimeBound that should have been executed last time
	CronConsistency bool `json:"CronConsistency"`

	// controller is not going to schedule anything below now - StartingDeadlineSeconds
	StartingDeadlineSeconds *int64 `json:"StartingDeadlineSeconds,omitempty"`
}

type NodeOvercommitConfigStatus struct {
	// NodeList which the nodeOvercommitConfig rules matched
	MatchedNodeList []string `json:"matchedNodeList,omitempty"`

	// LastScheduleTime recorded the time when the cronTab was triggered last time
	LastScheduleTime *metav1.Time `json:"lastScheduleTime,omitempty"`
}

type TimeBound struct {
	// TimeBound is valid between Start and End
	// +optional
	Start *metav1.Time `json:"start,omitempty"`

	// +optional
	End *metav1.Time `json:"end,omitempty"`

	// +optional
	Bounds []Bound `json:"bounds,omitempty"`
}

type Bound struct {
	// CronTab support standardSpec, requires 5 entries
	// representing: minute, hour, day of month, month and day of week
	CronTab string `json:"cronTab"`

	// Node overcommit ratio will be set as ResourceOvercommitRatio when CronTab is triggered
	// cpu,memory are supported.
	// +optional
	ResourceOvercommitRatio map[v1.ResourceName]string `json:"resourceOvercommitRatio,omitempty"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,shortName=noc
// +kubebuilder:printcolumn:name="OVERCOMMITRATIO",type=string,JSONPath=".spec.resourceOvercommitRatio"
// +kubebuilder:printcolumn:name="SELECTOR",type=string,JSONPath=".spec.nodeOvercommitSelectorVal"

// NodeOvercommitConfig is the Schema for the nodeovercommitconfigs API
type NodeOvercommitConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NodeOvercommitConfigSpec   `json:"spec,omitempty"`
	Status NodeOvercommitConfigStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
//+kubebuilder:object:root=true

// NodeOvercommitConfigList contains a list of NodeOvercommitConfig
type NodeOvercommitConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NodeOvercommitConfig `json:"items"`
}
