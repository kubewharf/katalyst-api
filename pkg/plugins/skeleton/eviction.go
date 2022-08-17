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

package skeleton

import (
	"github.com/kubewharf/katalyst-api/pkg/protocol/evictionplugin/v1alpha1"
	pluginapi "github.com/kubewharf/katalyst-api/pkg/protocol/evictionplugin/v1alpha1"
)

const (
	fakeEvictionPluginName = "fake-eviction-plugin"
)

// EvictionPlugin defines the interface that eviction plugins defined out of katalyst should follow.
type EvictionPlugin interface {
	GenericPlugin
	pluginapi.EvictionPluginServer
}

// DummyEvictionPlugin defines dummy eviction plugin
type DummyEvictionPlugin struct {
	v1alpha1.UnimplementedEvictionPluginServer
}

// Name of the dummy eviction plugin
func (DummyEvictionPlugin) Name() string { return fakeEvictionPluginName }

// Start the dummy eviction plugin
func (DummyEvictionPlugin) Start() error { return nil }

// Stop the dummy eviction plugin
func (DummyEvictionPlugin) Stop() error { return nil }

// ResourceName returns name of the dummy eviction plugin
func (DummyEvictionPlugin) ResourceName() string { return "" }

var _ EvictionPlugin = &DummyEvictionPlugin{
	v1alpha1.UnimplementedEvictionPluginServer{},
}
