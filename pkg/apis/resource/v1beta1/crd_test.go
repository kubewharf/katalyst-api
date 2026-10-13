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
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"sigs.k8s.io/yaml"
)

func TestComputeDomainCRD(t *testing.T) {
	crd := loadCRD(t, "resource.nvidia.com_computedomains.yaml")
	assertCRDIdentity(t, crd, "computedomains.resource.nvidia.com", ComputeDomainKind, ComputeDomainKind+"List")

	version := crd.Spec.Versions[0]
	if version.Subresources == nil || version.Subresources.Status == nil {
		t.Fatal("ComputeDomain status subresource is not configured")
	}

	spec := requiredProperty(t, version.Schema.OpenAPIV3Schema, "spec")
	numNodes := requiredProperty(t, &spec, "numNodes")
	if got, want := numNodes.Type, "integer"; got != want {
		t.Errorf("spec.numNodes type = %q, want %q", got, want)
	}

	channel := requiredProperty(t, &spec, "channel")
	template := requiredProperty(t, &channel, "resourceClaimTemplate")
	if !contains(template.Required, "name") {
		t.Errorf("spec.channel.resourceClaimTemplate required = %v, want name", template.Required)
	}
	allocationMode := requiredProperty(t, &channel, "allocationMode")
	if got, want := allocationMode.Enum, []string{"All", "Single"}; !reflect.DeepEqual(got, want) {
		t.Errorf("spec.channel.allocationMode enum = %v, want %v", got, want)
	}

	status := requiredProperty(t, version.Schema.OpenAPIV3Schema, "status")
	nodes := requiredProperty(t, &status, "nodes")
	if got, want := nodes.XListType, ptr("map"); !reflect.DeepEqual(got, want) {
		t.Errorf("status.nodes list type = %v, want %v", got, want)
	}
	if got, want := nodes.XListMapKeys, []string{"name"}; !reflect.DeepEqual(got, want) {
		t.Errorf("status.nodes list map keys = %v, want %v", got, want)
	}
}

func TestComputeDomainCliqueCRD(t *testing.T) {
	crd := loadCRD(t, "resource.nvidia.com_computedomaincliques.yaml")
	assertCRDIdentity(t, crd, "computedomaincliques.resource.nvidia.com", ComputeDomainCliqueKind, ComputeDomainCliqueKind+"List")

	daemons := requiredProperty(t, crd.Spec.Versions[0].Schema.OpenAPIV3Schema, "daemons")
	if got, want := daemons.XListType, ptr("map"); !reflect.DeepEqual(got, want) {
		t.Errorf("daemons list type = %v, want %v", got, want)
	}
	if got, want := daemons.XListMapKeys, []string{"nodeName"}; !reflect.DeepEqual(got, want) {
		t.Errorf("daemons list map keys = %v, want %v", got, want)
	}
}

func loadCRD(t *testing.T, name string) *crdManifest {
	t.Helper()

	path := filepath.Join("..", "..", "..", "..", "config", "crd", "bases", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("os.ReadFile(%q) error = %v", path, err)
	}

	var crd crdManifest
	if err := yaml.Unmarshal(data, &crd); err != nil {
		t.Fatalf("yaml.Unmarshal(%q) error = %v", path, err)
	}
	return &crd
}

func assertCRDIdentity(t *testing.T, crd *crdManifest, name, kind, listKind string) {
	t.Helper()

	if crd.Metadata.Name != name {
		t.Errorf("metadata.name = %q, want %q", crd.Metadata.Name, name)
	}
	if crd.Spec.Group != GroupName {
		t.Errorf("spec.group = %q, want %q", crd.Spec.Group, GroupName)
	}
	if crd.Spec.Scope != "Namespaced" {
		t.Errorf("spec.scope = %q, want Namespaced", crd.Spec.Scope)
	}
	if crd.Spec.Names.Kind != kind {
		t.Errorf("spec.names.kind = %q, want %q", crd.Spec.Names.Kind, kind)
	}
	if crd.Spec.Names.ListKind != listKind {
		t.Errorf("spec.names.listKind = %q, want %q", crd.Spec.Names.ListKind, listKind)
	}
	if len(crd.Spec.Versions) != 1 {
		t.Fatalf("len(spec.versions) = %d, want 1", len(crd.Spec.Versions))
	}
	version := crd.Spec.Versions[0]
	if version.Name != Version || !version.Served || !version.Storage {
		t.Errorf("spec.versions[0] = %#v, want name=%q served=true storage=true", version, Version)
	}
	if version.Schema == nil || version.Schema.OpenAPIV3Schema == nil {
		t.Fatal("spec.versions[0].schema.openAPIV3Schema is nil")
	}
}

func requiredProperty(t *testing.T, schema *jsonSchema, name string) jsonSchema {
	t.Helper()

	property, ok := schema.Properties[name]
	if !ok {
		t.Fatalf("schema property %q is missing", name)
	}
	return property
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func ptr[T any](value T) *T {
	return &value
}

type crdManifest struct {
	Metadata crdMetadata `json:"metadata"`
	Spec     crdSpec     `json:"spec"`
}

type crdMetadata struct {
	Name string `json:"name"`
}

type crdSpec struct {
	Group    string       `json:"group"`
	Names    crdNames     `json:"names"`
	Scope    string       `json:"scope"`
	Versions []crdVersion `json:"versions"`
}

type crdNames struct {
	Kind     string `json:"kind"`
	ListKind string `json:"listKind"`
}

type crdVersion struct {
	Name         string           `json:"name"`
	Served       bool             `json:"served"`
	Storage      bool             `json:"storage"`
	Schema       *crdSchema       `json:"schema"`
	Subresources *crdSubresources `json:"subresources"`
}

type crdSchema struct {
	OpenAPIV3Schema *jsonSchema `json:"openAPIV3Schema"`
}

type crdSubresources struct {
	Status *struct{} `json:"status"`
}

type jsonSchema struct {
	Type         string                `json:"type"`
	Properties   map[string]jsonSchema `json:"properties"`
	Required     []string              `json:"required"`
	Enum         []string              `json:"enum"`
	XListType    *string               `json:"x-kubernetes-list-type"`
	XListMapKeys []string              `json:"x-kubernetes-list-map-keys"`
}
