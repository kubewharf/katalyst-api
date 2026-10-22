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
	"encoding/json"
	"reflect"
	"testing"
)

func TestBECPUPSIEvictionConfigRoundTripAndDeepCopy(t *testing.T) {
	t.Parallel()
	// Explicit false and zero must survive JSON; nil must remain omitted.
	input := []byte(`{"spec":{"config":{"evictionConfig":{"beCPUPSIEvictionConfig":{"enable":false,"enableTaint":true,"enableEviction":false,"windowSeconds":30,"taintThreshold":0.2,"evictionThreshold":0.5,"recoveryThreshold":0.1,"taintDurationSeconds":60,"evictionDurationSeconds":120,"cooldownSeconds":120,"budgetWindowSeconds":600,"maxEvictionsPerWindow":3,"maxSampleGapSeconds":30,"minOutsideUsageRatio":0.2,"minRecoveryIdleRatio":0.2,"gracePeriodSeconds":0}}}}}`)
	var original AdminQoSConfiguration
	if err := json.Unmarshal(input, &original); err != nil {
		t.Fatal(err)
	}
	copy := original.DeepCopy()
	if !reflect.DeepEqual(&original, copy) {
		t.Fatal("deep copy changed configuration")
	}
	a := reflect.ValueOf(original.Spec.Config.EvictionConfig.BECPUPSIEvictionConfig).Elem()
	b := reflect.ValueOf(copy.Spec.Config.EvictionConfig.BECPUPSIEvictionConfig).Elem()
	for i := 0; i < a.NumField(); i++ {
		if a.Field(i).IsNil() || b.Field(i).IsNil() {
			t.Fatalf("field %s was lost", a.Type().Field(i).Name)
		}
		if a.Field(i).Pointer() == b.Field(i).Pointer() {
			t.Fatalf("field %s aliases original", a.Type().Field(i).Name)
		}
	}
	data, err := json.Marshal(copy)
	if err != nil {
		t.Fatal(err)
	}
	var decoded AdminQoSConfiguration
	if err = json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(&original, &decoded) {
		t.Fatal("JSON round trip changed configuration")
	}
	empty := &BECPUPSIEvictionConfig{}
	data, err = json.Marshal(empty.DeepCopy())
	if err != nil || string(data) != "{}" {
		t.Fatalf("nil fields should remain omitted: %s, %v", data, err)
	}
}
