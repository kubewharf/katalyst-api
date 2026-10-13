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
	"context"
	"net"
	"reflect"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	drapluginapi "k8s.io/kubelet/pkg/apis/dra/v1"
	drapluginapiv1beta1 "k8s.io/kubelet/pkg/apis/dra/v1beta1"
	watcherapi "k8s.io/kubelet/pkg/apis/pluginregistration/v1"
	qrmpluginapi "k8s.io/kubelet/pkg/apis/resourceplugin/v1alpha1"
)

const fakeDRAPluginName = "fake-dra-plugin"

type fakeDRAPlugin struct {
	drapluginapi.UnimplementedDRAPluginServer

	prepareCalls int
}

func (*fakeDRAPlugin) Name() string { return fakeDRAPluginName }

func (*fakeDRAPlugin) Start() error { return nil }

func (*fakeDRAPlugin) Stop() error { return nil }

func (p *fakeDRAPlugin) NodePrepareResources(
	_ context.Context,
	_ *drapluginapi.NodePrepareResourcesRequest,
) (*drapluginapi.NodePrepareResourcesResponse, error) {
	p.prepareCalls++
	return &drapluginapi.NodePrepareResourcesResponse{}, nil
}

func (*fakeDRAPlugin) NodeUnprepareResources(
	_ context.Context,
	_ *drapluginapi.NodeUnprepareResourcesRequest,
) (*drapluginapi.NodeUnprepareResourcesResponse, error) {
	return &drapluginapi.NodeUnprepareResourcesResponse{}, nil
}

var _ DRAPlugin = (*fakeDRAPlugin)(nil)

func TestDRARegistrationMetadata(t *testing.T) {
	t.Parallel()

	wrapper, err := NewRegistrationPluginWrapper(&fakeDRAPlugin{}, nil, nil)
	if err != nil {
		t.Fatalf("NewRegistrationPluginWrapper() error = %v", err)
	}

	got, err := wrapper.GetInfo(context.Background(), &watcherapi.InfoRequest{})
	if err != nil {
		t.Fatalf("GetInfo() error = %v", err)
	}

	if got.Type != watcherapi.DRAPlugin {
		t.Errorf("GetInfo().Type = %q, want %q", got.Type, watcherapi.DRAPlugin)
	}
	if got.Name != fakeDRAPluginName {
		t.Errorf("GetInfo().Name = %q, want %q", got.Name, fakeDRAPluginName)
	}
	if got.Endpoint != "" {
		t.Errorf("GetInfo().Endpoint = %q, want empty endpoint", got.Endpoint)
	}

	wantVersions := []string{
		drapluginapi.DRAPluginService,
		drapluginapiv1beta1.DRAPluginService,
	}
	if !reflect.DeepEqual(got.SupportedVersions, wantVersions) {
		t.Errorf("GetInfo().SupportedVersions = %v, want %v", got.SupportedVersions, wantVersions)
	}
}

func TestDRAServerRegister(t *testing.T) {
	t.Parallel()

	wrapper, err := NewRegistrationPluginWrapper(&fakeDRAPlugin{}, nil, nil)
	if err != nil {
		t.Fatalf("NewRegistrationPluginWrapper() error = %v", err)
	}

	server := grpc.NewServer()
	wrapper.serverRegister(server)
	services := server.GetServiceInfo()

	for _, serviceName := range []string{
		drapluginapi.DRAPlugin_ServiceDesc.ServiceName,
		drapluginapiv1beta1.DRAPlugin_ServiceDesc.ServiceName,
	} {
		if _, ok := services[serviceName]; !ok {
			t.Errorf("registered services do not contain %q: %v", serviceName, services)
		}
	}
}

func TestDRAVersionCompatibility(t *testing.T) {
	t.Parallel()

	plugin := &fakeDRAPlugin{}
	wrapper, err := NewRegistrationPluginWrapper(plugin, nil, nil)
	if err != nil {
		t.Fatalf("NewRegistrationPluginWrapper() error = %v", err)
	}

	listener := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()
	wrapper.serverRegister(server)
	go func() {
		_ = server.Serve(listener)
	}()
	t.Cleanup(server.Stop)
	t.Cleanup(func() {
		_ = listener.Close()
	})

	ctx := context.Background()
	conn, err := grpc.DialContext(
		ctx,
		"bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return listener.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("grpc.DialContext() error = %v", err)
	}
	t.Cleanup(func() {
		_ = conn.Close()
	})

	if _, err := drapluginapi.NewDRAPluginClient(conn).NodePrepareResources(
		ctx,
		&drapluginapi.NodePrepareResourcesRequest{},
	); err != nil {
		t.Fatalf("v1 NodePrepareResources() error = %v", err)
	}

	if _, err := drapluginapiv1beta1.NewDRAPluginClient(conn).NodePrepareResources(
		ctx,
		&drapluginapiv1beta1.NodePrepareResourcesRequest{},
	); err != nil {
		t.Fatalf("v1beta1 NodePrepareResources() error = %v", err)
	}

	if plugin.prepareCalls != 2 {
		t.Errorf("NodePrepareResources() calls = %d, want 2", plugin.prepareCalls)
	}
}

func TestQRMRegistrationUnchanged(t *testing.T) {
	t.Parallel()

	wrapper, err := NewRegistrationPluginWrapper(&DummyQRMPlugin{}, nil, nil)
	if err != nil {
		t.Fatalf("NewRegistrationPluginWrapper() error = %v", err)
	}

	got, err := wrapper.GetInfo(context.Background(), &watcherapi.InfoRequest{})
	if err != nil {
		t.Fatalf("GetInfo() error = %v", err)
	}

	if got.Type != watcherapi.ResourcePlugin {
		t.Errorf("GetInfo().Type = %q, want %q", got.Type, watcherapi.ResourcePlugin)
	}
	if got.Name != "test_resource" {
		t.Errorf("GetInfo().Name = %q, want %q", got.Name, "test_resource")
	}
	if !reflect.DeepEqual(got.SupportedVersions, []string{qrmpluginapi.Version}) {
		t.Errorf(
			"GetInfo().SupportedVersions = %v, want %v",
			got.SupportedVersions,
			[]string{qrmpluginapi.Version},
		)
	}
}
