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
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// GroupName is the group name used in this package
	GroupName string = "node.katalyst.kubewharf.io"
)

// SchemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1alpha1"}

// ResourceNameKatalystCNR is used to construct standard gvr
const (
	ResourceNameKatalystCNR = "customnoderesources"
	ResourceNameKatalystNPD = "nodeprofiledescriptors"
)

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

var (
	// SchemeBuilder collects schemas to build.
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	// AddToScheme is used by generated client to add this scheme to the generated client.
	AddToScheme = SchemeBuilder.AddToScheme
)

// Adds the list of known types to the given scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&CustomNodeResource{},
		&CustomNodeResourceList{},
		&NodeProfileDescriptor{},
		&NodeProfileDescriptorList{},
	)

	v1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
