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

package externalversions_test

import (
	"strings"
	"testing"
	"time"

	resourcev1beta1 "github.com/kubewharf/katalyst-api/pkg/apis/resource/v1beta1"
	clientfake "github.com/kubewharf/katalyst-api/pkg/client/clientset/versioned/fake"
	externalversions "github.com/kubewharf/katalyst-api/pkg/client/informers/externalversions"
)

func TestForResourceRegistersComputeDomainInformers(t *testing.T) {
	factory := externalversions.NewSharedInformerFactory(clientfake.NewSimpleClientset(), time.Minute)

	for _, resource := range []string{
		resourcev1beta1.ResourceNameComputeDomains,
		resourcev1beta1.ResourceNameComputeDomainCliques,
	} {
		gvr := resourcev1beta1.SchemeGroupVersion.WithResource(resource)
		informer, err := factory.ForResource(gvr)
		if err != nil {
			t.Errorf("ForResource(%s) error = %v", gvr, err)
			continue
		}
		if informer.Informer() == nil {
			t.Errorf("ForResource(%s).Informer() = nil", gvr)
		}
	}

	unknown := resourcev1beta1.SchemeGroupVersion.WithResource("unknown")
	if _, err := factory.ForResource(unknown); err == nil {
		t.Fatalf("ForResource(%s) error = nil, want no informer found", unknown)
	} else if !strings.Contains(err.Error(), "no informer found") {
		t.Fatalf("ForResource(%s) error = %v, want no informer found", unknown, err)
	}
}
