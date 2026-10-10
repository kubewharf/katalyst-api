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
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/yaml"
)

func TestQRMPluginConfigRDTAndBulkheadRDTConfig(t *testing.T) {
	disableRDT := true
	enableCPUList := true
	enableCAT := true
	defaultCATWays := intstr.FromString("MaxCATWays")
	exclusiveClosIDs := []string{"clos-a", "peer-b"}
	config := QRMPluginConfig{
		RDTConfig: &RDTConfig{
			DisableRDT: &disableRDT,
		},
		CPUPluginConfig: &CPUPluginConfig{
			BulkheadConfig: &BulkheadConfig{
				BulkheadRDTConfig: &BulkheadRDTConfig{
					EnableCPUList:  &enableCPUList,
					EnableCAT:      &enableCAT,
					DefaultCATWays: &defaultCATWays,
					ClosCATWays: map[string]intstr.IntOrString{
						"share-00": intstr.FromString("MaxCATWays-MinCATWays"),
						"share-01": intstr.FromInt(2),
					},
					CATPolicy: &CATPolicy{
						DefaultPlacement: &CATPlacementPolicy{
							AllowedBitUsages: []CATBitUsage{CATBitUsageAll},
							Direction:        CATAllocationDirectionLow,
						},
						ClosPlacements: map[string]CATPlacementPolicy{
							"share-01": {
								AllowedBitUsages: []CATBitUsage{CATBitUsageSoftware, CATBitUsageExclusive},
								Direction:        CATAllocationDirectionHigh,
							},
						},
						ExclusiveClosIDs: &exclusiveClosIDs,
					},
				},
			},
		},
	}

	raw, err := json.Marshal(config)
	if err != nil {
		t.Fatalf("marshal config: %v", err)
	}

	var got QRMPluginConfig
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal config: %v", err)
	}
	if got.RDTConfig == nil || got.RDTConfig.DisableRDT == nil || !*got.RDTConfig.DisableRDT {
		t.Fatal("DisableRDT was not preserved")
	}
	bulkheadRDTConfig := got.CPUPluginConfig.BulkheadConfig.BulkheadRDTConfig
	if bulkheadRDTConfig == nil || bulkheadRDTConfig.EnableCPUList == nil || !*bulkheadRDTConfig.EnableCPUList {
		t.Fatal("EnableCPUList was not preserved")
	}
	if bulkheadRDTConfig.EnableCAT == nil || !*bulkheadRDTConfig.EnableCAT {
		t.Fatal("EnableCAT was not preserved")
	}
	if bulkheadRDTConfig.DefaultCATWays == nil || bulkheadRDTConfig.DefaultCATWays.String() != "MaxCATWays" {
		t.Fatalf("DefaultCATWays = %v, want MaxCATWays", bulkheadRDTConfig.DefaultCATWays)
	}
	share00CATWays := bulkheadRDTConfig.ClosCATWays["share-00"]
	if got := share00CATWays.String(); got != "MaxCATWays-MinCATWays" {
		t.Fatalf("ClosCATWays[share-00] = %s, want MaxCATWays-MinCATWays", got)
	}
	share01CATWays := bulkheadRDTConfig.ClosCATWays["share-01"]
	if got := share01CATWays.IntValue(); got != 2 {
		t.Fatalf("ClosCATWays[share-01] = %d, want 2", got)
	}
	if bulkheadRDTConfig.CATPolicy == nil || bulkheadRDTConfig.CATPolicy.ExclusiveClosIDs == nil {
		t.Fatalf("CATPolicy exclusiveClosIDs = %#v, want pointer", bulkheadRDTConfig.CATPolicy)
	}
	if got := *bulkheadRDTConfig.CATPolicy.ExclusiveClosIDs; len(got) != 2 || got[0] != "clos-a" || got[1] != "peer-b" {
		t.Fatalf("ExclusiveClosIDs = %#v, want clos-a and peer-b", got)
	}
	if got := bulkheadRDTConfig.CATPolicy.DefaultPlacement.Direction; got != CATAllocationDirectionLow {
		t.Fatalf("DefaultPlacement.Direction = %s, want %s", got, CATAllocationDirectionLow)
	}
	if got := bulkheadRDTConfig.CATPolicy.ClosPlacements["share-01"].Direction; got != CATAllocationDirectionHigh {
		t.Fatalf("ClosPlacements[share-01].Direction = %s, want %s", got, CATAllocationDirectionHigh)
	}

	copied := got.DeepCopy()
	copiedPolicy := copied.CPUPluginConfig.BulkheadConfig.BulkheadRDTConfig.CATPolicy
	copiedPolicy.DefaultPlacement.AllowedBitUsages[0] = CATBitUsageExclusive
	share01Placement := copiedPolicy.ClosPlacements["share-01"]
	share01Placement.AllowedBitUsages[0] = CATBitUsageExclusive
	copiedPolicy.ClosPlacements["share-01"] = share01Placement
	(*copiedPolicy.ExclusiveClosIDs)[0] = "mutated"
	require.Equal(t, CATBitUsageAll, bulkheadRDTConfig.CATPolicy.DefaultPlacement.AllowedBitUsages[0])
	require.Equal(t, CATBitUsageSoftware, bulkheadRDTConfig.CATPolicy.ClosPlacements["share-01"].AllowedBitUsages[0])
	require.Equal(t, "clos-a", (*bulkheadRDTConfig.CATPolicy.ExclusiveClosIDs)[0])
}

func TestCPUProvisionConfigFillDefaultSharePoolRoundTrip(t *testing.T) {
	enabled := true
	in := &CPUProvisionConfig{
		FillDefaultSharePoolWithNonReclaimCPUs: &enabled,
	}
	data, err := json.Marshal(in)
	require.NoError(t, err)
	require.JSONEq(t, `{"fillDefaultSharePoolWithNonReclaimCPUs":true}`, string(data))
	out := &CPUProvisionConfig{}
	require.NoError(t, json.Unmarshal(data, out))
	require.NotNil(t, out.FillDefaultSharePoolWithNonReclaimCPUs)
	require.True(t, *out.FillDefaultSharePoolWithNonReclaimCPUs)
	copied := in.DeepCopy()
	require.NotSame(t, in.FillDefaultSharePoolWithNonReclaimCPUs, copied.FillDefaultSharePoolWithNonReclaimCPUs)
	require.True(t, *copied.FillDefaultSharePoolWithNonReclaimCPUs)
}

func TestMemoryGuardConfigCriticalWatermarkSource(t *testing.T) {
	t.Run("JSON round trip and optional omission", func(t *testing.T) {
		source := CriticalWatermarkSourceHigh
		in := &MemoryGuardConfig{CriticalWatermarkSource: &source}

		data, err := json.Marshal(in)
		require.NoError(t, err)
		require.JSONEq(t, `{"criticalWatermarkSource":"high"}`, string(data))

		out := &MemoryGuardConfig{}
		require.NoError(t, json.Unmarshal(data, out))
		require.NotNil(t, out.CriticalWatermarkSource)
		require.Equal(t, CriticalWatermarkSourceHigh, *out.CriticalWatermarkSource)

		data, err = json.Marshal(&MemoryGuardConfig{})
		require.NoError(t, err)
		require.JSONEq(t, `{}`, string(data))
	})

	t.Run("deep copy owns an independent pointer", func(t *testing.T) {
		source := CriticalWatermarkSourceLow
		in := &MemoryGuardConfig{CriticalWatermarkSource: &source}

		copied := in.DeepCopy()
		require.NotNil(t, copied.CriticalWatermarkSource)
		require.NotSame(t, in.CriticalWatermarkSource, copied.CriticalWatermarkSource)
		require.Equal(t, CriticalWatermarkSourceLow, *copied.CriticalWatermarkSource)

		*copied.CriticalWatermarkSource = CriticalWatermarkSourceHigh
		require.Equal(t, CriticalWatermarkSourceLow, *in.CriticalWatermarkSource)
	})

	t.Run("CRD schema exposes the low and high enum", func(t *testing.T) {
		data, err := os.ReadFile("../../../../config/crd/bases/config.katalyst.kubewharf.io_adminqosconfigurations.yaml")
		require.NoError(t, err)

		jsonData, err := yaml.YAMLToJSON(data)
		require.NoError(t, err)

		var crd struct {
			Spec struct {
				Versions []struct {
					Schema struct {
						OpenAPIV3Schema adminQoSSchema `json:"openAPIV3Schema"`
					} `json:"schema"`
				} `json:"versions"`
			} `json:"spec"`
		}
		require.NoError(t, json.Unmarshal(jsonData, &crd))
		require.NotEmpty(t, crd.Spec.Versions)

		schema := crd.Spec.Versions[0].Schema.OpenAPIV3Schema
		schema = adminQoSSchemaProperty(t, schema, "spec")
		schema = adminQoSSchemaProperty(t, schema, "config")
		schema = adminQoSSchemaProperty(t, schema, "advisorConfig")
		schema = adminQoSSchemaProperty(t, schema, "memoryAdvisorConfig")
		schema = adminQoSSchemaProperty(t, schema, "memoryGuardConfig")
		criticalWatermarkSource := adminQoSSchemaProperty(t, schema, "criticalWatermarkSource")

		require.Equal(t, "string", criticalWatermarkSource.Type)
		require.Equal(t, []string{"low", "high"}, criticalWatermarkSource.Enum)
	})
}

type adminQoSSchema struct {
	Properties map[string]adminQoSSchema `json:"properties"`
	Type       string                    `json:"type"`
	Enum       []string                  `json:"enum"`
}

func adminQoSSchemaProperty(t *testing.T, schema adminQoSSchema, name string) adminQoSSchema {
	t.Helper()
	property, ok := schema.Properties[name]
	require.Truef(t, ok, "schema property %q not found", name)
	return property
}
