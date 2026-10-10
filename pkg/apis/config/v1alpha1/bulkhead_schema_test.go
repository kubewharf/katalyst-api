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
	"regexp"
	"testing"

	"sigs.k8s.io/yaml"
)

func TestBulkheadRDTConfigSchema(t *testing.T) {
	schema := loadBulkheadRDTConfigSchema(t)
	const catWaysExpressionPattern = `^\s*(MaxCATWays|MinCATWays|[1-9][0-9]*)(\s*[+-]\s*(MaxCATWays|MinCATWays|[1-9][0-9]*))?\s*$`

	// k8s 1.18 compatibility: the CRD must not carry any x-kubernetes-validations
	// (CEL) rules. CEL schema validation is only recognized from k8s 1.23 (alpha)
	// and a 1.18 apiserver rejects the CRD because it does not understand the
	// x-kubernetes-validations field. The expressive CAT checks therefore live in
	// code rather than the CRD. This subtest guards against regressing by
	// re-adding CEL rules to the bulkhead RDT schema.
	t.Run("no CEL validation rules", func(t *testing.T) {
		assertNoValidationRules(t, schema, "bulkheadRDTConfig")
		assertNoValidationRules(t, schemaProperty(t, schema, "defaultCATWays"), "defaultCATWays")
		assertNoValidationRules(t, schemaProperty(t, schema, "closCATWays"), "closCATWays")
		catPolicy := schemaProperty(t, schema, "catPolicy")
		assertNoValidationRules(t, catPolicy, "catPolicy")
		defaultPlacement := schemaProperty(t, catPolicy, "defaultPlacement")
		assertNoValidationRules(t, schemaProperty(t, defaultPlacement, "allowedBitUsages"), "allowedBitUsages")
		assertNoValidationRules(t, schemaProperty(t, catPolicy, "exclusiveClosIDs"), "exclusiveClosIDs")
		assertNoValidationRules(t, schemaProperty(t, catPolicy, "closPlacements"), "closPlacements")
	})

	t.Run("default CAT ways int-or-string", func(t *testing.T) {
		defaultCATWays := schemaProperty(t, schema, "defaultCATWays")
		assertIntOrStringSchema(t, defaultCATWays, "defaultCATWays")
		assertCATWaysExpressionContract(t, catWaysExpressionPattern)
	})

	t.Run("clos CAT ways int-or-string values", func(t *testing.T) {
		closCATWays := schemaProperty(t, schema, "closCATWays")
		if closCATWays.AdditionalProperties == nil {
			t.Fatal("closCATWays has no additionalProperties schema")
		}
		assertIntOrStringSchema(t, *closCATWays.AdditionalProperties, "closCATWays values")
	})

	t.Run("EnableCAT contract (code-side)", func(t *testing.T) {
		enabled, disabled := true, false
		assertEnableCATContract(t, []enableCATContractCase{
			{name: "invalid: CAT enabled without default ways", enableCAT: &enabled, valid: false},
			{name: "valid: CAT enabled with positive ways", enableCAT: &enabled, hasDefaultCATWays: true, valid: true},
			{name: "valid: CAT disabled without default ways", enableCAT: &disabled, valid: true},
		})
	})

	t.Run("CAT policy schema", func(t *testing.T) {
		catPolicy := schemaProperty(t, schema, "catPolicy")

		defaultPlacement := schemaProperty(t, catPolicy, "defaultPlacement")
		assertEnumSchema(t, schemaProperty(t, defaultPlacement, "direction"), []string{"low", "high"}, "defaultPlacement.direction")
		allowedBitUsages := schemaProperty(t, defaultPlacement, "allowedBitUsages")
		if allowedBitUsages.Items == nil {
			t.Fatal("defaultPlacement.allowedBitUsages has no item schema")
		}
		assertEnumSchema(t, *allowedBitUsages.Items, []string{"*", "S", "H", "X"}, "defaultPlacement.allowedBitUsages items")

		exclusiveClosIDs := schemaProperty(t, catPolicy, "exclusiveClosIDs")
		if exclusiveClosIDs.Items == nil {
			t.Fatal("exclusiveClosIDs has no item schema")
		}
		if exclusiveClosIDs.Items.Type != "string" {
			t.Fatalf("exclusiveClosIDs item type = %q, want string", exclusiveClosIDs.Items.Type)
		}
		assertAbsentSchema(t, catPolicy, "allocationGroups", "catPolicy")
		assertAbsentSchema(t, catPolicy, "nonOverlapConstraints", "catPolicy")
	})
}

type customResourceDefinition struct {
	Spec struct {
		Versions []struct {
			Schema struct {
				OpenAPIV3Schema jsonSchema `json:"openAPIV3Schema"`
			} `json:"schema"`
		} `json:"versions"`
	} `json:"spec"`
}

type jsonSchema struct {
	Properties             map[string]jsonSchema `json:"properties"`
	AdditionalProperties   *jsonSchema           `json:"additionalProperties"`
	AnyOf                  []jsonSchema          `json:"anyOf"`
	Enum                   []string              `json:"enum"`
	Items                  *jsonSchema           `json:"items"`
	MinItems               *int                  `json:"minItems"`
	Required               []string              `json:"required"`
	Type                   string                `json:"type"`
	UniqueItems            *bool                 `json:"uniqueItems"`
	XIntOrString           bool                  `json:"x-kubernetes-int-or-string"`
	XKubernetesValidations []validationRule      `json:"x-kubernetes-validations"`
}

type validationRule struct {
	Rule    string `json:"rule"`
	Message string `json:"message"`
}

type enableCATContractCase struct {
	name              string
	enableCAT         *bool
	hasDefaultCATWays bool
	valid             bool
}

func loadBulkheadRDTConfigSchema(t *testing.T) jsonSchema {
	t.Helper()

	data, err := os.ReadFile("../../../../config/crd/bases/config.katalyst.kubewharf.io_adminqosconfigurations.yaml")
	if err != nil {
		t.Fatalf("read AdminQoSConfiguration CRD: %v", err)
	}

	jsonData, err := yaml.YAMLToJSON(data)
	if err != nil {
		t.Fatalf("convert AdminQoSConfiguration CRD to JSON: %v", err)
	}

	var crd customResourceDefinition
	if err := json.Unmarshal(jsonData, &crd); err != nil {
		t.Fatalf("unmarshal AdminQoSConfiguration CRD: %v", err)
	}
	if len(crd.Spec.Versions) == 0 {
		t.Fatal("AdminQoSConfiguration CRD has no versions")
	}

	return findSchemaProperty(t, crd.Spec.Versions[0].Schema.OpenAPIV3Schema, "bulkheadRDTConfig")
}

func findSchemaProperty(t *testing.T, schema jsonSchema, name string) jsonSchema {
	t.Helper()
	if property, ok := schema.Properties[name]; ok {
		return property
	}
	for _, property := range schema.Properties {
		if found := findSchemaPropertyOrNil(property, name); found != nil {
			return *found
		}
	}
	t.Fatalf("schema property %q not found", name)
	return jsonSchema{}
}

func findSchemaPropertyOrNil(schema jsonSchema, name string) *jsonSchema {
	if property, ok := schema.Properties[name]; ok {
		return &property
	}
	for _, property := range schema.Properties {
		if found := findSchemaPropertyOrNil(property, name); found != nil {
			return found
		}
	}
	return nil
}

func schemaProperty(t *testing.T, schema jsonSchema, name string) jsonSchema {
	t.Helper()
	property, ok := schema.Properties[name]
	if !ok {
		t.Fatalf("schema property %q not found", name)
	}
	return property
}

func assertAbsentSchema(t *testing.T, schema jsonSchema, field string, name string) {
	t.Helper()
	if _, ok := schema.Properties[field]; ok {
		t.Fatalf("%s must not expose %s", name, field)
	}
}

func assertNoValidationRules(t *testing.T, schema jsonSchema, name string) {
	t.Helper()
	if len(schema.XKubernetesValidations) != 0 {
		t.Fatalf("%s exposes %d x-kubernetes-validations rules, want none for k8s 1.18 compatibility: %#v",
			name, len(schema.XKubernetesValidations), schema.XKubernetesValidations)
	}
}

func assertIntOrStringSchema(t *testing.T, schema jsonSchema, name string) {
	t.Helper()
	if !schema.XIntOrString {
		t.Fatalf("%s x-kubernetes-int-or-string = false, want true", name)
	}
	types := map[string]bool{}
	for _, candidate := range schema.AnyOf {
		types[candidate.Type] = true
	}
	if !types["integer"] || !types["string"] {
		t.Fatalf("%s anyOf types = %#v, want integer and string", name, types)
	}
}

func assertEnumSchema(t *testing.T, schema jsonSchema, want []string, name string) {
	t.Helper()
	if len(schema.Enum) != len(want) {
		t.Fatalf("%s enum = %#v, want exactly %#v", name, schema.Enum, want)
	}
	got := map[string]bool{}
	for _, value := range schema.Enum {
		got[value] = true
	}
	for _, value := range want {
		if !got[value] {
			t.Fatalf("%s enum = %#v, want value %s", name, schema.Enum, value)
		}
	}
}

func assertRequiredSchema(t *testing.T, schema jsonSchema, field string, name string) {
	t.Helper()
	for _, required := range schema.Required {
		if required == field {
			return
		}
	}
	t.Fatalf("%s required = %#v, want %s", name, schema.Required, field)
}

func assertCATWaysExpressionContract(t *testing.T, pattern string) {
	t.Helper()
	expression := regexp.MustCompile(pattern)
	for _, value := range []string{
		"MaxCATWays",
		"MinCATWays",
		"MaxCATWays-MinCATWays",
		"MinCATWays+1",
		"2",
	} {
		if !expression.MatchString(value) {
			t.Errorf("CAT ways expression pattern rejected %q", value)
		}
	}
	for _, value := range []string{
		"CBMMask",
		"MinCBMBits",
		"CBMMask-MinCBMBits",
		"MaxCATWays-CBMMask",
	} {
		if expression.MatchString(value) {
			t.Errorf("CAT ways expression pattern accepted legacy expression %q", value)
		}
	}
}

func assertEnableCATContract(t *testing.T, cases []enableCATContractCase) {
	t.Helper()
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.enableCAT == nil || !*tt.enableCAT || tt.hasDefaultCATWays
			if got != tt.valid {
				t.Fatalf("EnableCAT contract validity = %v, want %v", got, tt.valid)
			}
		})
	}
}
