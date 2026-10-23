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

package v1alpha1_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kubewharf/katalyst-api/pkg/apis/config/v1alpha1"
)

// TestReclaimedConsumerPercentageSourceCompatibility locks the source type of
// ReclaimedConsumerToReclaimedResourcePercentage at *map[string]int. Existing
// downstream callers assign plain map[string]int values; a distinct named value
// type would break that direct assignment at compile time, so the field must
// keep its released v1alpha1 type.
func TestReclaimedConsumerPercentageSourceCompatibility(t *testing.T) {
	legacy := map[string]int{"batch": 80, "online": 20}
	cfg := &v1alpha1.ReclaimedResourceConfig{
		ReclaimedConsumerToReclaimedResourcePercentage: &legacy,
	}

	require.NotNil(t, cfg.ReclaimedConsumerToReclaimedResourcePercentage)
	require.Equal(t, 80, (*cfg.ReclaimedConsumerToReclaimedResourcePercentage)["batch"])
	require.Equal(t, 20, (*cfg.ReclaimedConsumerToReclaimedResourcePercentage)["online"])
}
