//go:build !windows
// +build !windows

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
	"strings"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	pluginapi "github.com/kubewharf/katalyst-api/pkg/protocol/evictionplugin/v1alpha1"
)

const testEvictionPluginName = "test-eviction-plugin"

type testEvictionPlugin struct {
	pluginapi.UnimplementedEvictionPluginServer
}

func (p *testEvictionPlugin) Name() string { return testEvictionPluginName }
func (p *testEvictionPlugin) Start() error { return nil }
func (p *testEvictionPlugin) Stop() error  { return nil }
func (p *testEvictionPlugin) GetEvictPods(_ context.Context, _ *pluginapi.GetEvictPodsRequest) (
	*pluginapi.GetEvictPodsResponse, error,
) {
	return &pluginapi.GetEvictPodsResponse{}, nil
}

func TestPluginRegistrationWrapperGRPCMaxRecvMsgSize(t *testing.T) {
	t.Parallel()

	wrapper, err := NewRegistrationPluginWrapper(&testEvictionPlugin{}, []string{t.TempDir()}, nil)
	if err != nil {
		t.Fatalf("create plugin wrapper: %v", err)
	}
	if err := wrapper.serve(); err != nil {
		t.Fatalf("serve plugin: %v", err)
	}
	t.Cleanup(func() {
		if err := wrapper.stop(); err != nil {
			t.Errorf("stop plugin wrapper: %v", err)
		}
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, wrapper.sockets[0],
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
			return (&net.Dialer{}).DialContext(ctx, "unix", addr)
		}),
	)
	if err != nil {
		t.Fatalf("dial plugin: %v", err)
	}
	t.Cleanup(func() {
		if err := conn.Close(); err != nil {
			t.Errorf("close plugin connection: %v", err)
		}
	})

	client := pluginapi.NewEvictionPluginClient(conn)
	newRequest := func(size int) *pluginapi.GetEvictPodsRequest {
		return &pluginapi.GetEvictPodsRequest{
			ActivePods: []*corev1.Pod{{
				ObjectMeta: metav1.ObjectMeta{Annotations: map[string]string{"payload": strings.Repeat("a", size)}},
			}},
		}
	}

	if _, err := client.GetEvictPods(ctx, newRequest(5<<20)); err != nil {
		t.Fatalf("5 MiB request should be accepted: %v", err)
	}

	_, err = client.GetEvictPods(ctx, newRequest(9<<20))
	if status.Code(err) != codes.ResourceExhausted {
		t.Fatalf("9 MiB request should be rejected with ResourceExhausted, got: %v", err)
	}
}
