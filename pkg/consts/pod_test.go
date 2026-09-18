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

package consts

import (
	"testing"
)

func TestPodSaleModeConstants(t *testing.T) {
	t.Parallel()

	tests := map[string]string{
		"PodAnnotationSaleModeKey": PodAnnotationSaleModeKey,
		"PodSaleModeSpot":          PodSaleModeSpot,
		"PodSaleModeScheduled":     PodSaleModeScheduled,
		"PodSaleModeReserved":      PodSaleModeReserved,
		"PodSaleModeUnknown":       PodSaleModeUnknown,
	}
	wants := map[string]string{
		"PodAnnotationSaleModeKey": "katalyst.kubewharf.io/sale_mode",
		"PodSaleModeSpot":          "spot",
		"PodSaleModeScheduled":     "scheduled",
		"PodSaleModeReserved":      "reserved",
		"PodSaleModeUnknown":       "unknown",
	}

	for name, got := range tests {
		if got != wants[name] {
			t.Errorf("%s = %q, want %q", name, got, wants[name])
		}
	}
}
