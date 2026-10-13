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

package v1beta1

import (
	"encoding/json"
	"reflect"
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
)

func TestSchemeGroupVersion(t *testing.T) {
	if got, want := SchemeGroupVersion.String(), "resource.nvidia.com/v1beta1"; got != want {
		t.Fatalf("SchemeGroupVersion.String() = %q, want %q", got, want)
	}
}

func TestAddToSchemeRegistersComputeDomainTypes(t *testing.T) {
	scheme := runtime.NewScheme()
	if err := AddToScheme(scheme); err != nil {
		t.Fatalf("AddToScheme() error = %v", err)
	}

	for _, kind := range []string{
		ComputeDomainKind,
		ComputeDomainKind + "List",
		ComputeDomainCliqueKind,
		ComputeDomainCliqueKind + "List",
	} {
		gvk := SchemeGroupVersion.WithKind(kind)
		if _, err := scheme.New(gvk); err != nil {
			t.Errorf("scheme.New(%s) error = %v", gvk, err)
		}
	}
}

func TestComputeDomainJSONRoundTrip(t *testing.T) {
	input := []byte(`{
		"apiVersion":"resource.nvidia.com/v1beta1",
		"kind":"ComputeDomain",
		"spec":{
			"numNodes":2,
			"channel":{
				"resourceClaimTemplate":{"name":"imex-channel"},
				"allocationMode":"All"
			}
		},
		"status":{
			"status":"Ready",
			"nodes":[{
				"name":"node-a",
				"ipAddress":"10.0.0.1",
				"cliqueID":"clique-a",
				"index":1,
				"status":"Ready"
			}]
		}
	}`)

	var domain ComputeDomain
	if err := json.Unmarshal(input, &domain); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	output, err := json.Marshal(&domain)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	var roundTripped map[string]interface{}
	if err := json.Unmarshal(output, &roundTripped); err != nil {
		t.Fatalf("json.Unmarshal(round trip) error = %v", err)
	}
	var original map[string]interface{}
	if err := json.Unmarshal(input, &original); err != nil {
		t.Fatalf("json.Unmarshal(original) error = %v", err)
	}
	if !reflect.DeepEqual(roundTripped["spec"], original["spec"]) {
		t.Errorf("round-tripped spec = %#v, want %#v", roundTripped["spec"], original["spec"])
	}
	if !reflect.DeepEqual(roundTripped["status"], original["status"]) {
		t.Errorf("round-tripped status = %#v, want %#v", roundTripped["status"], original["status"])
	}
}

func TestComputeDomainChannelConfigJSONKeys(t *testing.T) {
	config := ComputeDomainChannelConfig{
		DomainID:       "domain-id",
		DomainName:     "domain-name",
		AllocationMode: ComputeDomainChannelAllocationModeAll,
	}
	data, err := json.Marshal(config)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	var fields map[string]interface{}
	if err := json.Unmarshal(data, &fields); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	for key, want := range map[string]string{
		"domainID":       "domain-id",
		"domainName":     "domain-name",
		"allocationMode": ComputeDomainChannelAllocationModeAll,
	} {
		if got := fields[key]; got != want {
			t.Errorf("field %q = %#v, want %q", key, got, want)
		}
	}
}

func TestDefaultComputeDomainConfigs(t *testing.T) {
	tests := []struct {
		name       string
		apiVersion string
		kind       string
	}{
		{
			name:       "channel",
			apiVersion: DefaultComputeDomainChannelConfig().APIVersion,
			kind:       DefaultComputeDomainChannelConfig().Kind,
		},
		{
			name:       "daemon",
			apiVersion: DefaultComputeDomainDaemonConfig().APIVersion,
			kind:       DefaultComputeDomainDaemonConfig().Kind,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, want := tt.apiVersion, "resource.nvidia.com/v1beta1"; got != want {
				t.Errorf("APIVersion = %q, want %q", got, want)
			}
			wantKind := ComputeDomainChannelConfigKind
			if tt.name == "daemon" {
				wantKind = ComputeDomainDaemonConfigKind
			}
			if tt.kind != wantKind {
				t.Errorf("Kind = %q, want %q", tt.kind, wantKind)
			}
		})
	}
}
