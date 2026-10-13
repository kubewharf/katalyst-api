/*
 * Copyright (c) 2024, NVIDIA CORPORATION.  All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
)

// Interface defines the set of common APIs for all configs.
// +k8s:deepcopy-gen=false
type Interface interface {
	Normalize() error
	Validate() error
}

// StrictDecoder rejects unknown fields in this API group.
var StrictDecoder runtime.Decoder

// NonstrictDecoder drops unknown fields in this API group.
var NonstrictDecoder runtime.Decoder

func init() {
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(SchemeGroupVersion,
		&ComputeDomainChannelConfig{},
		&ComputeDomainDaemonConfig{},
		&ComputeDomain{},
		&ComputeDomainClique{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)

	NonstrictDecoder = json.NewSerializerWithOptions(
		json.DefaultMetaFactory,
		scheme,
		scheme,
		json.SerializerOptions{Strict: false},
	)
	StrictDecoder = json.NewSerializerWithOptions(
		json.DefaultMetaFactory,
		scheme,
		scheme,
		json.SerializerOptions{Strict: true},
	)
}
