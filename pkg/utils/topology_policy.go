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

package utils

import (
	"k8s.io/kubelet/config/v1beta1"
	"k8s.io/kubernetes/pkg/kubelet/apis/config"

	nodev1alpha1 "github.com/kubewharf/katalyst-api/pkg/apis/node/v1alpha1"
)

// GenerateTopologyPolicy generates TopologyPolicy which presents both Topology manager policy and scope.
func GenerateTopologyPolicy(policy string, scope string) nodev1alpha1.TopologyPolicy {
	switch scope {
	default:
		return nodev1alpha1.TopologyPolicyNone
	}
}

func generateTopologyPolicyPodScope(policy string) nodev1alpha1.TopologyPolicy {
	switch policy {
	case config.RestrictedTopologyManagerPolicy:
		return nodev1alpha1.TopologyPolicyRestrictedPodLevel
	case config.BestEffortTopologyManagerPolicy:
		return nodev1alpha1.TopologyPolicyBestEffortPodLevel
	case config.NoneTopologyManagerPolicy:
		return nodev1alpha1.TopologyPolicyNone
	case v1beta1.NumericTopologyManagerPolicy:
		return nodev1alpha1.TopologyPolicyNumericPodLevel
	default:
		return nodev1alpha1.TopologyPolicyNone
	}
}

func generateTopologyPolicyContainerScope(policy string) nodev1alpha1.TopologyPolicy {
	switch policy {
	case config.RestrictedTopologyManagerPolicy:
		return nodev1alpha1.TopologyPolicyRestrictedContainerLevel
	case config.BestEffortTopologyManagerPolicy:
		return nodev1alpha1.TopologyPolicyBestEffortContainerLevel
	case config.NoneTopologyManagerPolicy:
		return nodev1alpha1.TopologyPolicyNone
	case v1beta1.NumericTopologyManagerPolicy:
		return nodev1alpha1.TopologyPolicyNumericContainerLevel
	default:
		return nodev1alpha1.TopologyPolicyNone
	}
}
