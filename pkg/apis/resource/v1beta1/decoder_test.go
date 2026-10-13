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
	"strings"
	"testing"
)

func TestStrictDecoderDecodesComputeDomainChannelConfig(t *testing.T) {
	data := []byte(`{
		"apiVersion":"resource.nvidia.com/v1beta1",
		"kind":"ComputeDomainChannelConfig",
		"domainID":"domain-id",
		"domainName":"domain-name",
		"allocationMode":"All"
	}`)

	object, _, err := StrictDecoder.Decode(data, nil, nil)
	if err != nil {
		t.Fatalf("StrictDecoder.Decode() error = %v", err)
	}
	config, ok := object.(*ComputeDomainChannelConfig)
	if !ok {
		t.Fatalf("decoded object type = %T, want *ComputeDomainChannelConfig", object)
	}
	if config.DomainID != "domain-id" || config.DomainName != "domain-name" ||
		config.AllocationMode != ComputeDomainChannelAllocationModeAll {
		t.Fatalf("decoded config = %#v", config)
	}
}

func TestStrictDecoderRejectsUnknownField(t *testing.T) {
	data := []byte(`{
		"apiVersion":"resource.nvidia.com/v1beta1",
		"kind":"ComputeDomainChannelConfig",
		"domainID":"domain-id",
		"domainName":"domain-name",
		"unknownField":"value"
	}`)

	_, _, err := StrictDecoder.Decode(data, nil, nil)
	if err == nil {
		t.Fatal("StrictDecoder.Decode() error = nil, want unknown-field error")
	}
	if !strings.Contains(err.Error(), "unknownField") {
		t.Fatalf("StrictDecoder.Decode() error = %v, want unknownField", err)
	}
}

func TestNonstrictDecoderAcceptsUnknownField(t *testing.T) {
	data := []byte(`{
		"apiVersion":"resource.nvidia.com/v1beta1",
		"kind":"ComputeDomainChannelConfig",
		"domainID":"domain-id",
		"domainName":"domain-name",
		"unknownField":"value"
	}`)

	object, _, err := NonstrictDecoder.Decode(data, nil, nil)
	if err != nil {
		t.Fatalf("NonstrictDecoder.Decode() error = %v", err)
	}
	config, ok := object.(*ComputeDomainChannelConfig)
	if !ok {
		t.Fatalf("decoded object type = %T, want *ComputeDomainChannelConfig", object)
	}
	if config.DomainID != "domain-id" || config.DomainName != "domain-name" {
		t.Fatalf("decoded config = %#v", config)
	}
}

func TestComputeDomainChannelConfigValidateRejectsEmptyDomainID(t *testing.T) {
	err := (&ComputeDomainChannelConfig{}).Validate()
	if err == nil {
		t.Fatal("Validate() error = nil, want domainID validation error")
	}
	if got, want := err.Error(), "domainID cannot be empty"; got != want {
		t.Fatalf("Validate() error = %q, want %q", got, want)
	}
}
