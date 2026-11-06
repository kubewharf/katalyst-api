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

func TestReclaimedMemorySafeZoneGateConfig(t *testing.T) {
	t.Run("JSON round trip and optional omission", func(t *testing.T) {
		enabled := true
		threshold := 0.5
		in := &MemoryPressureEvictionConfig{
			EnableReclaimedMemorySafeZoneGate:  &enabled,
			ReclaimedMemoryUsageRatioThreshold: &threshold,
		}

		data, err := json.Marshal(in)
		require.NoError(t, err)
		require.Contains(t, string(data), `"enableReclaimedMemorySafeZoneGate":true`)
		require.Contains(t, string(data), `"reclaimedMemoryUsageRatioThreshold":0.5`)

		out := &MemoryPressureEvictionConfig{}
		require.NoError(t, json.Unmarshal(data, out))
		require.NotNil(t, out.EnableReclaimedMemorySafeZoneGate)
		require.True(t, *out.EnableReclaimedMemorySafeZoneGate)
		require.NotNil(t, out.ReclaimedMemoryUsageRatioThreshold)
		require.Equal(t, 0.5, *out.ReclaimedMemoryUsageRatioThreshold)

		// Explicit false round-trips as non-nil false and is preserved.
		disabled := false
		falseIn := &MemoryPressureEvictionConfig{EnableReclaimedMemorySafeZoneGate: &disabled}
		falseData, err := json.Marshal(falseIn)
		require.NoError(t, err)
		require.Contains(t, string(falseData), `"enableReclaimedMemorySafeZoneGate":false`)
		falseOut := &MemoryPressureEvictionConfig{}
		require.NoError(t, json.Unmarshal(falseData, falseOut))
		require.NotNil(t, falseOut.EnableReclaimedMemorySafeZoneGate)
		require.False(t, *falseOut.EnableReclaimedMemorySafeZoneGate)

		data, err = json.Marshal(&MemoryPressureEvictionConfig{})
		require.NoError(t, err)
		require.NotContains(t, string(data), "enableReclaimedMemorySafeZoneGate")
		require.NotContains(t, string(data), "reclaimedMemoryUsageRatioThreshold")

		// Omitted objects keep both pointers nil after round-trip, so omission
		// stays distinguishable from explicit false / zero values.
		legacy := `{"enableNumaLevelEviction":true}`
		legacyOut := &MemoryPressureEvictionConfig{}
		require.NoError(t, json.Unmarshal([]byte(legacy), legacyOut))
		require.Nil(t, legacyOut.EnableReclaimedMemorySafeZoneGate)
		require.Nil(t, legacyOut.ReclaimedMemoryUsageRatioThreshold)
	})

	t.Run("deep copy owns independent pointers", func(t *testing.T) {
		enabled := true
		threshold := 0.25
		in := &MemoryPressureEvictionConfig{
			EnableReclaimedMemorySafeZoneGate:  &enabled,
			ReclaimedMemoryUsageRatioThreshold: &threshold,
		}

		copied := in.DeepCopy()
		require.NotSame(t, in.EnableReclaimedMemorySafeZoneGate, copied.EnableReclaimedMemorySafeZoneGate)
		require.NotSame(t, in.ReclaimedMemoryUsageRatioThreshold, copied.ReclaimedMemoryUsageRatioThreshold)

		*copied.EnableReclaimedMemorySafeZoneGate = false
		*copied.ReclaimedMemoryUsageRatioThreshold = 0.75
		require.True(t, *in.EnableReclaimedMemorySafeZoneGate)
		require.Equal(t, 0.25, *in.ReclaimedMemoryUsageRatioThreshold)
	})

	t.Run("CRD schema exposes bool flag and 0..1 threshold", func(t *testing.T) {
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
		schema = adminQoSSchemaProperty(t, schema, "evictionConfig")
		schema = adminQoSSchemaProperty(t, schema, "memoryPressureEvictionConfig")

		flag := adminQoSSchemaProperty(t, schema, "enableReclaimedMemorySafeZoneGate")
		require.Equal(t, "boolean", flag.Type)
		require.Empty(t, flag.Enum)

		threshold := adminQoSSchemaProperty(t, schema, "reclaimedMemoryUsageRatioThreshold")
		require.Equal(t, "number", threshold.Type)
		require.NotNil(t, threshold.Minimum)
		require.Equal(t, 0.0, *threshold.Minimum)
		require.NotNil(t, threshold.Maximum)
		require.Equal(t, 1.0, *threshold.Maximum)
	})
}

func TestCPUHeadroomAdaptiveUtilBasedConfigRoundTrip(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name string
		data string
	}{
		{name: "legacy config", data: `{"enable":true,"targetReclaimedCoreUtilization":0.6}`},
		{name: "empty adaptive config", data: `{"adaptiveConfig":{}}`},
		{name: "enabled with omitted parameters", data: `{"adaptiveConfig":{"enable":true}}`},
		{name: "explicit false and zero", data: `{"adaptiveConfig":{"enable":false,"reclaimedCoreUtilizationTolerance":0}}`},
		{
			name: "custom parameters",
			data: `{"adaptiveConfig":{"enable":true,"targetReclaimedCoreUtilizationLowerBound":0.7,"targetReclaimedCoreUtilizationUpperBound":0.93,"reclaimedCoreUtilizationTolerance":0.02}}`,
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var config CPUHeadroomUtilBasedConfig
			require.NoError(t, json.Unmarshal([]byte(tt.data), &config))

			data, err := json.Marshal(config.DeepCopy())
			require.NoError(t, err)
			require.JSONEq(t, tt.data, string(data))
		})
	}
}

func TestAdminQoSConfigAdaptiveCPUHeadroomDeepCopy(t *testing.T) {
	t.Parallel()

	const input = `{
		"reclaimedResourceConfig": {
			"cpuHeadroomConfig": {
				"utilBasedConfig": {
					"enable": true,
					"adaptiveConfig": {
						"enable": true,
						"targetReclaimedCoreUtilizationLowerBound": 0.6,
						"targetReclaimedCoreUtilizationUpperBound": 0.9,
						"reclaimedCoreUtilizationTolerance": 0.05
					}
				}
			}
		}
	}`
	var config AdminQoSConfig
	require.NoError(t, json.Unmarshal([]byte(input), &config))

	copied := config.DeepCopy()
	original := config.ReclaimedResourceConfig.CPUHeadroomConfig.UtilBasedConfig.AdaptiveConfig
	adaptive := copied.ReclaimedResourceConfig.CPUHeadroomConfig.UtilBasedConfig.AdaptiveConfig
	require.NotSame(t, original, adaptive)
	require.NotSame(t, original.Enable, adaptive.Enable)
	require.NotSame(t, original.TargetReclaimedCoreUtilizationLowerBound, adaptive.TargetReclaimedCoreUtilizationLowerBound)
	require.NotSame(t, original.TargetReclaimedCoreUtilizationUpperBound, adaptive.TargetReclaimedCoreUtilizationUpperBound)
	require.NotSame(t, original.ReclaimedCoreUtilizationTolerance, adaptive.ReclaimedCoreUtilizationTolerance)

	*adaptive.Enable = false
	*adaptive.TargetReclaimedCoreUtilizationLowerBound = 0.7
	*adaptive.TargetReclaimedCoreUtilizationUpperBound = 0.95
	*adaptive.ReclaimedCoreUtilizationTolerance = 0

	data, err := json.Marshal(config)
	require.NoError(t, err)
	require.JSONEq(t, input, string(data))
	require.Nil(t, (*CPUHeadroomAdaptiveUtilBasedConfig)(nil).DeepCopy())
}

func TestCPUHeadroomAdaptiveUtilBasedConfigSchema(t *testing.T) {
	t.Parallel()

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
	for _, name := range []string{"spec", "config", "reclaimedResourceConfig", "cpuHeadroomConfig", "utilBasedConfig"} {
		schema = adminQoSSchemaProperty(t, schema, name)
	}
	require.NotContains(t, schema.Required, "adaptiveConfig")
	adaptive := adminQoSSchemaProperty(t, schema, "adaptiveConfig")
	require.Equal(t, "object", adaptive.Type)
	require.Empty(t, adaptive.Required)
	require.Nil(t, adaptive.Default)
	require.Len(t, adaptive.Properties, 4)

	enable := adminQoSSchemaProperty(t, adaptive, "enable")
	require.Equal(t, "boolean", enable.Type)
	require.Nil(t, enable.Default)

	for _, name := range []string{"targetReclaimedCoreUtilizationLowerBound", "targetReclaimedCoreUtilizationUpperBound", "reclaimedCoreUtilizationTolerance"} {
		field := adminQoSSchemaProperty(t, adaptive, name)
		require.Equal(t, "number", field.Type)
		require.NotNil(t, field.Minimum)
		require.Equal(t, 0.0, *field.Minimum)
		require.NotNil(t, field.Maximum)
		require.Equal(t, 1.0, *field.Maximum)
		require.Equal(t, name != "reclaimedCoreUtilizationTolerance", field.ExclusiveMinimum)
		require.True(t, field.ExclusiveMaximum)
		require.Nil(t, field.Default)
	}
}

type adminQoSSchema struct {
	Properties       map[string]adminQoSSchema `json:"properties"`
	Type             string                    `json:"type"`
	Enum             []string                  `json:"enum"`
	Minimum          *float64                  `json:"minimum"`
	Maximum          *float64                  `json:"maximum"`
	ExclusiveMinimum bool                      `json:"exclusiveMinimum"`
	ExclusiveMaximum bool                      `json:"exclusiveMaximum"`
	Required         []string                  `json:"required"`
	Default          *json.RawMessage          `json:"default"`
}

func adminQoSSchemaProperty(t *testing.T, schema adminQoSSchema, name string) adminQoSSchema {
	t.Helper()
	property, ok := schema.Properties[name]
	require.Truef(t, ok, "schema property %q not found", name)
	return property
}
