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

package fake_test

import (
	"context"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stesting "k8s.io/client-go/testing"

	resourcev1beta1 "github.com/kubewharf/katalyst-api/pkg/apis/resource/v1beta1"
	clientfake "github.com/kubewharf/katalyst-api/pkg/client/clientset/versioned/fake"
)

func TestResourceV1beta1FakeClient(t *testing.T) {
	const namespace = "test"
	client := clientfake.NewSimpleClientset(
		&resourcev1beta1.ComputeDomain{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "selected",
				Namespace: namespace,
				Labels:    map[string]string{"role": "selected"},
			},
		},
		&resourcev1beta1.ComputeDomain{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "other",
				Namespace: namespace,
				Labels:    map[string]string{"role": "other"},
			},
		},
	)
	domains := client.ResourceV1beta1().ComputeDomains(namespace)

	domain, err := domains.Get(context.Background(), "selected", metav1.GetOptions{})
	if err != nil {
		t.Fatalf("ComputeDomains().Get() error = %v", err)
	}
	if domain.Name != "selected" {
		t.Errorf("ComputeDomains().Get().Name = %q, want selected", domain.Name)
	}

	list, err := domains.List(context.Background(), metav1.ListOptions{LabelSelector: "role=selected"})
	if err != nil {
		t.Fatalf("ComputeDomains().List() error = %v", err)
	}
	if len(list.Items) != 1 || list.Items[0].Name != "selected" {
		t.Fatalf("ComputeDomains().List().Items = %#v, want selected only", list.Items)
	}

	domain.Status.Status = resourcev1beta1.ComputeDomainStatusReady
	if _, err := domains.UpdateStatus(context.Background(), domain, metav1.UpdateOptions{}); err != nil {
		t.Fatalf("ComputeDomains().UpdateStatus() error = %v", err)
	}
	if !hasStatusUpdateAction(client.Actions()) {
		t.Fatal("ComputeDomains().UpdateStatus() did not record an update action for status")
	}

	clique := &resourcev1beta1.ComputeDomainClique{
		ObjectMeta: metav1.ObjectMeta{Name: "clique", Namespace: namespace},
	}
	created, err := client.ResourceV1beta1().ComputeDomainCliques(namespace).
		Create(context.Background(), clique, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("ComputeDomainCliques().Create() error = %v", err)
	}
	if created.Name != clique.Name {
		t.Errorf("ComputeDomainCliques().Create().Name = %q, want %q", created.Name, clique.Name)
	}
}

func hasStatusUpdateAction(actions []k8stesting.Action) bool {
	for _, action := range actions {
		if action.GetVerb() == "update" && action.GetSubresource() == "status" {
			return true
		}
	}
	return false
}
