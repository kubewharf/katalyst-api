package utils

import (
	"k8s.io/kubernetes/pkg/kubelet/apis/config"

	nodev1alpha1 "github.com/kubewharf/katalyst-api/pkg/apis/node/v1alpha1"
)

// GenerateTopologyPolicy generates TopologyPolicy which presents both Topology manager policy and scope.
func GenerateTopologyPolicy(policy string, scope string) nodev1alpha1.TopologyPolicy {
	switch scope {
	case config.PodTopologyManagerScope:
		return generateTopologyPolicyPodScope(policy)
	case config.ContainerTopologyManagerScope:
		return generateTopologyPolicyContainerScope(policy)
	default:
		return nodev1alpha1.TopologyPolicyNone
	}
}

func generateTopologyPolicyPodScope(policy string) nodev1alpha1.TopologyPolicy {
	switch policy {
	case config.SingleNumaNodeTopologyManagerPolicy:
		return nodev1alpha1.TopologyPolicySingleNUMANodePodLevel
	case config.RestrictedTopologyManagerPolicy:
		return nodev1alpha1.TopologyPolicyRestrictedPodLevel
	case config.BestEffortTopologyManagerPolicy:
		return nodev1alpha1.TopologyPolicyBestEffortPodLevel
	case config.NoneTopologyManagerPolicy:
		return nodev1alpha1.TopologyPolicyNone
	default:
		return nodev1alpha1.TopologyPolicyNone
	}
}

func generateTopologyPolicyContainerScope(policy string) nodev1alpha1.TopologyPolicy {
	switch policy {
	case config.SingleNumaNodeTopologyManagerPolicy:
		return nodev1alpha1.TopologyPolicySingleNUMANodeContainerLevel
	case config.RestrictedTopologyManagerPolicy:
		return nodev1alpha1.TopologyPolicyRestrictedContainerLevel
	case config.BestEffortTopologyManagerPolicy:
		return nodev1alpha1.TopologyPolicyBestEffortContainerLevel
	case config.NoneTopologyManagerPolicy:
		return nodev1alpha1.TopologyPolicyNone
	default:
		return nodev1alpha1.TopologyPolicyNone
	}
}
