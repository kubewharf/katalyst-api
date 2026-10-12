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

package v1alpha1

import (
	"testing"

	"github.com/gogo/protobuf/proto"
)

func TestConditionZoneTaintsWireRoundTrip(t *testing.T) {
	t.Parallel()

	condition := &Condition{ZoneTaints: []*TopologyZoneTaint{
		{
			Path:          []*TopologyZoneRef{{Type: "Socket", Name: "0"}, {Type: "Numa", Name: "0"}},
			ConditionName: "MemorySuppression",
			QosLevel:      "reclaimed_cores",
			Effect:        "NoSchedule",
			MetCondition:  true,
		},
		{
			Path:          []*TopologyZoneRef{{Type: "Socket", Name: "0"}, {Type: "Numa", Name: "1"}},
			ConditionName: "MemorySuppression",
			QosLevel:      "shared_cores",
			Effect:        "PreferNoSchedule",
			MetCondition:  true,
		},
	}}
	for _, message := range []proto.Message{
		&GetEvictPodsResponse{Condition: condition},
		&ThresholdMetResponse{MetType: ThresholdMetType_NOT_MET, Condition: condition},
	} {
		data, err := proto.Marshal(message)
		if err != nil {
			t.Fatal(err)
		}
		copy := proto.Clone(message)
		copy.Reset()
		if err := proto.Unmarshal(data, copy); err != nil {
			t.Fatal(err)
		}
		if !proto.Equal(message, copy) {
			t.Fatalf("zone taints changed during wire round trip: got %v, want %v", copy, message)
		}
	}
}
