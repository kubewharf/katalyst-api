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

package skeleton

import (
	drapluginapi "k8s.io/kubelet/pkg/apis/dra/v1"
)

const (
	dummyDRAPluginName = "fake-dra-plugin"
)

// DRAPlugin defines the interface that DRA plugins outside katalyst-api must implement.
// Name returns the DRA driver name advertised to kubelet.
type DRAPlugin interface {
	GenericPlugin
	drapluginapi.DRAPluginServer
}

// DummyDRAPlugin defines a dummy DRA plugin.
type DummyDRAPlugin struct {
	drapluginapi.UnimplementedDRAPluginServer
}

// Name of the dummy DRA plugin.
func (DummyDRAPlugin) Name() string { return dummyDRAPluginName }

// Start the dummy DRA plugin.
func (DummyDRAPlugin) Start() error { return nil }

// Stop the dummy DRA plugin.
func (DummyDRAPlugin) Stop() error { return nil }

var _ DRAPlugin = &DummyDRAPlugin{
	drapluginapi.UnimplementedDRAPluginServer{},
}
