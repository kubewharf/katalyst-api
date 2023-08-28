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
	// Selector decides whether to update Nodes if the Node matches the selector when Type equals to LabelSelector
	// +optional
	Selector *metav1.LabelSelector `json:"nodeSelector,omitempty"`

	// NodeList decides whether to check and update nodes if node in NodeList when Type equals to NodeList
	// +optional
	NodeList []string `json:"nodeList,omitempty"`

	// ResourceOvercommitRatioConfig describes the resource overcommit ratio that needs to inject into Node.Annotations
	// cpu,memory are supported.
	// +optional
	ResourceOvercommitRatioConfig map[v1.ResourceName]string `json:"resourceOvercommitRatioConfig,omitempty"`
}

type NodeOvercommitConfigStatus struct {
	// NodeList which the nodeOvercommitConfig rules matched
	MatchedNodeList []string `json:"matchedNodeList,omitempty"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,shortName=noc
// +kubebuilder:printcolumn:name="OVERCOMMITRATIO",type=string,JSONPath=".spec.resourceOvercommitRatioConfig"

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
