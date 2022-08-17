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
	pluginapi "k8s.io/kubelet/pkg/apis/resourceplugin/v1alpha1"
)

const (
	fakeQRMPluginName = "fake-qrm-plugin"
)

// QRMPlugin defines the interface that QoS-Resource plugins defined out of katalyst should follow.
type QRMPlugin interface {
	GenericPlugin
	pluginapi.ResourcePluginServer

	ResourceName() string
}

// DummyQRMPlugin defines dummy QoS-Resource plugin
type DummyQRMPlugin struct {
	pluginapi.UnimplementedResourcePluginServer
}

// Name of the dummy QoS-Resource plugin
func (DummyQRMPlugin) Name() string { return fakeQRMPluginName }

// Start the dummy QoS-Resource plugin
func (DummyQRMPlugin) Start() error { return nil }

// Stop the dummy QoS-Resource plugin
func (DummyQRMPlugin) Stop() error { return nil }

// ResourceName returns name of the dummy QoS-Resource plugin
func (DummyQRMPlugin) ResourceName() string { return "test_resource" }

var _ QRMPlugin = &DummyQRMPlugin{
	pluginapi.UnimplementedResourcePluginServer{},
}
