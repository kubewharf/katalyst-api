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

package registration

import (
	"context"

	"k8s.io/klog/v2"

	watcherapi "k8s.io/kubelet/pkg/apis/pluginregistration/v1"
	plugincache "k8s.io/kubernetes/pkg/kubelet/pluginmanager/cache"
)

const (
	EvictionPlugin    = "EvictionPlugin"
	ReporterPlugin    = "ReporterPlugin"
	QoSResourcePlugin = "QoSResourcePlugin"

	BaseVersion = "v1alpha1"
)

type AgentPluginHandler interface {
	GetHandlerType() string
	plugincache.PluginHandler
}

type RegistrationHandler struct {
	pluginType        string
	pluginName        string
	supportedVersions []string
}

func NewRegistrationHandler(pluginType, pluginName string, supportedVersions []string) watcherapi.RegistrationServer {
	return &RegistrationHandler{
		pluginType:        pluginType,
		pluginName:        pluginName,
		supportedVersions: supportedVersions,
	}
}

// GetInfo is the RPC which return pluginInfo
func (handler *RegistrationHandler) GetInfo(ctx context.Context, req *watcherapi.InfoRequest) (*watcherapi.PluginInfo, error) {
	klog.Infof("%s plugin of type %s GetInfo called", handler.pluginName, handler.pluginType)
	return &watcherapi.PluginInfo{
		Type:              handler.pluginType,
		Name:              handler.pluginName,
		SupportedVersions: handler.supportedVersions}, nil
}

// NotifyRegistrationStatus receives the registration notification from watcher
func (handler *RegistrationHandler) NotifyRegistrationStatus(ctx context.Context, status *watcherapi.RegistrationStatus) (*watcherapi.RegistrationStatusResponse, error) {
	// TODO: consider strategy when register failed
	if !status.PluginRegistered {
		klog.Errorf("%s of type %s Registration failed: %v",
			handler.pluginName, handler.pluginType, status.Error)
	}
	return &watcherapi.RegistrationStatusResponse{}, nil
}
