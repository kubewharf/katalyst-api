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
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"sigs.k8s.io/yaml"
)

func TestPrefetcherRulesCopyAndRoundTrip(t *testing.T) {
	t.Parallel()
	conf := &PrefetcherTuningConfig{
		EnableTuner: true,
		Rules: []PrefetcherTuningRule{
			{LabelSelector: "psm=example.a", CPUCodeName: "AMD_K19Zen4", Disable: []Prefetcher{}},
			{StrategyName: "prefetcher_tuning_a", LabelSelector: "psm=example.a", CPUCodeName: "AMD_K19Zen4", Disable: []Prefetcher{"l2_stream"}},
		},
	}
	copied := conf.DeepCopy()
	copied.Rules[1].StrategyName = "prefetcher_tuning_b"
	copied.Rules[1].LabelSelector = "experiment=b"
	copied.Rules[1].Disable[0] = "up_down"
	require.Equal(t, "prefetcher_tuning_a", conf.Rules[1].StrategyName)
	require.Equal(t, "psm=example.a", conf.Rules[1].LabelSelector)
	require.Equal(t, Prefetcher("l2_stream"), conf.Rules[1].Disable[0])

	data, err := json.Marshal(conf)
	require.NoError(t, err)
	require.Contains(t, string(data), `"labelSelector":"psm=example.a"`)
	require.NotContains(t, string(data), `"psm":`)
	var restored PrefetcherTuningConfig
	require.NoError(t, json.Unmarshal(data, &restored))
	require.Equal(t, *conf, restored)
}

func TestPrefetcherRulesSchema(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile("../../../../config/crd/bases/config.katalyst.kubewharf.io_prefetchertunings.yaml")
	require.NoError(t, err)
	var crd customResourceDefinition
	require.NoError(t, yaml.Unmarshal(data, &crd))
	require.NotEmpty(t, crd.Spec.Versions)
	config := crd.Spec.Versions[0].Schema.OpenAPIV3Schema.Properties["spec"].Properties["config"]
	require.Contains(t, config.Required, "rules")
	rules := config.Properties["rules"]
	require.Equal(t, "array", rules.Type)
	require.ElementsMatch(t, []string{"labelSelector", "cpu_codename", "strategyName"}, rules.XKubernetesListMapKeys)
	require.NotNil(t, rules.Items)
	require.ElementsMatch(t, []string{"labelSelector", "cpu_codename", "disable"}, rules.Items.Required)
	require.NotContains(t, rules.Items.Properties, "psm")
	selector := rules.Items.Properties["labelSelector"]
	require.Equal(t, "string", selector.Type)
	require.NotNil(t, selector.MinLength)
	require.Equal(t, 1, *selector.MinLength)
	strategy := rules.Items.Properties["strategyName"]
	require.Equal(t, "string", strategy.Type)
	require.Equal(t, "", strategy.Default, "omitted strategyName must default for the list map key")
}

func TestPrefetcherDefaultRules(t *testing.T) {
	t.Parallel()
	var conf PrefetcherTuningConfig
	require.NoError(t, json.Unmarshal([]byte(`{"enableTuner":true,"rules":[{"labelSelector":"psm=example.a,experiment in (a,b)","cpu_codename":"AMD_K19Zen4","disable":[]}]}`), &conf))
	require.Len(t, conf.Rules, 1)
	require.Equal(t, "psm=example.a,experiment in (a,b)", conf.Rules[0].LabelSelector)
	require.Empty(t, conf.Rules[0].StrategyName)
	require.NotNil(t, conf.Rules[0].Disable)
}
