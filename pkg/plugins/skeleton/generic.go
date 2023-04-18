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
	"context"
	"fmt"
	"net"
	"os"
	"path"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog/v2"
	watcherapi "k8s.io/kubelet/pkg/apis/pluginregistration/v1"
	qrmpluginapi "k8s.io/kubelet/pkg/apis/resourceplugin/v1alpha1"

	"github.com/kubewharf/katalyst-api/pkg/plugins/registration"
	evictionv1apha1 "github.com/kubewharf/katalyst-api/pkg/protocol/evictionplugin/v1alpha1"
	reporterv1apha1 "github.com/kubewharf/katalyst-api/pkg/protocol/reporterplugin/v1alpha1"
)

const (
	restartRetryInterval = 5 * time.Second

	grpcTimeout = 5 * time.Second
)

// MetricCallback is used to return key metric as a callback function
type MetricCallback func(key string, value int64)

func dummyMetricCallback(_ string, _ int64) {}

// GenericPlugin is used to define a skeleton to write a standard katalyst plugin;
// and every katalyst plugin should implement those functions along with the extra
// interfaces defined in each individual plugin
type GenericPlugin interface {
	Name() string
	// Start and Stop function should be none-block.
	// Start should initialize all channels or goroutines and can be called again
	// after calling Stop to restart, and Stop should close all channels to stop
	// plugin goroutines, otherwise the wrapped manager will not work as expected
	Start() error
	Stop() error
}

// PluginRegistrationWrapper is a decorator for GenericPlugin implementations, it is
// responsible to handle the running states of those plugins.
type PluginRegistrationWrapper struct {
	sync.Mutex
	sockets   []string
	servers   []*grpc.Server
	stopCh    chan struct{}
	restartCh chan struct{}

	plugin         GenericPlugin
	metricCallback MetricCallback
	serverRegister pluginServerRegister

	watcherapi.RegistrationServer
}

// pluginServerRegister is used to register the given grpc server
type pluginServerRegister func(server *grpc.Server)

// NewRegistrationPluginWrapper wrap an input plugin with PluginRegistrationWrapper; and
// the returned PluginRegistrationWrapper also works as a GenericPlugin, besides the standard
// GenericPlugin functionality, it also handles the running states. The input GenericPlugin
// must support calling Start to restart after calling Stop.
func NewRegistrationPluginWrapper(plugin GenericPlugin, pluginsRegistrationDirs []string,
	metricCallback MetricCallback) (*PluginRegistrationWrapper, error) {
	if plugin == nil {
		return nil, fmt.Errorf("input report plugin is nil")
	}

	var (
		registrationServer watcherapi.RegistrationServer
		serverRegister     pluginServerRegister
	)
	switch t := plugin.(type) {
	case EvictionPlugin:
		registrationServer = registration.NewRegistrationHandler(
			registration.EvictionPlugin, plugin.Name(), []string{registration.BaseVersion})
		serverRegister = func(server *grpc.Server) {
			evictionv1apha1.RegisterEvictionPluginServer(server, t)
		}
	case ReporterPlugin:
		registrationServer = registration.NewRegistrationHandler(
			registration.ReporterPlugin, plugin.Name(), []string{registration.BaseVersion})
		serverRegister = func(server *grpc.Server) {
			reporterv1apha1.RegisterReporterPluginServer(server, t)
		}
	case QRMPlugin:
		registrationServer = registration.NewRegistrationHandler(
			watcherapi.ResourcePlugin, t.ResourceName(), []string{qrmpluginapi.Version})
		serverRegister = func(server *grpc.Server) {
			qrmpluginapi.RegisterResourcePluginServer(server, t)
		}
	default:
		return nil, fmt.Errorf("unsupported plugin type: %v", t.Name())
	}

	if metricCallback == nil {
		metricCallback = dummyMetricCallback
	}

	p := &PluginRegistrationWrapper{
		stopCh:             make(chan struct{}),
		restartCh:          make(chan struct{}),
		plugin:             plugin,
		metricCallback:     metricCallback,
		serverRegister:     serverRegister,
		RegistrationServer: registrationServer,
	}

	if len(pluginsRegistrationDirs) > 0 {
		p.sockets = make([]string, 0, len(pluginsRegistrationDirs))
		for _, dir := range pluginsRegistrationDirs {
			p.sockets = append(p.sockets, path.Join(dir, fmt.Sprintf("%s.sock", plugin.Name())))
		}
	}

	return p, nil
}

// Name of this wrapped plugin
func (p *PluginRegistrationWrapper) Name() string {
	return p.plugin.Name()
}

// Start the plugin with auto restart logic, besides
// it will register to plugin grpc server
func (p *PluginRegistrationWrapper) Start() error {
	go func() {
		for {
			_ = wait.PollImmediateInfinite(restartRetryInterval, func() (bool, error) {
				err := p.start()
				if err != nil {
					p.metricCallback("plugin_start_failed", 1)
					klog.Errorf("start plugin %s failed with err: %v", p.Name(), err)
					return false, nil
				}

				return true, nil
			})

			select {
			case <-p.stopCh:
				klog.Infof("plugin %s stopped, return from start goroutine", p.Name())
				return
			case _, ok := <-p.restartCh:
				if !ok {
					klog.Infof("plugin %s is stopping", p.Name())
					return
				}
			}

			klog.Infof("plugin %s received restart signal, ready to restart", p.Name())

			_ = wait.PollImmediateInfinite(restartRetryInterval, func() (bool, error) {
				err := p.stop()
				if err != nil {
					p.metricCallback("plugin_stop_failed", 1)
					klog.Errorf("stop plugin %s failed with err: %v", p.Name(), err)
					return false, nil
				}

				return true, nil
			})
		}
	}()

	return nil
}

// Stop will trigger this plugin to stop completely
func (p *PluginRegistrationWrapper) Stop() error {
	defer func() {
		close(p.stopCh)
		close(p.restartCh)
	}()
	return p.stop()
}

// Restart will trigger this plugin to restart if plugin
// has been stopped, it will return error
func (p *PluginRegistrationWrapper) Restart() (restartErr error) {
	// recover panic when restart a stopped plugin
	defer func() {
		if err := recover(); err != nil {
			restartErr = fmt.Errorf("restart plugin %s failed with err: %v", p.Name(), err)
		}
	}()

	p.restartCh <- struct{}{}
	return nil
}

// initialize plugin servers
func (p *PluginRegistrationWrapper) initialize() {
	p.servers = make([]*grpc.Server, 0, len(p.sockets))
}

func (p *PluginRegistrationWrapper) start() (startError error) {
	p.Lock()
	defer func() {
		p.Unlock()
		if startError != nil {
			klog.Errorf("start %s failed with error: %v", p.Name(), startError)
			_ = p.stop()
		}
	}()

	if len(p.sockets) == 0 {
		return fmt.Errorf("%s plugin get empty sockets", p.Name())
	}

	// we may call plugin.Stop multiple times, and even before plugin.Start
	// to ensure everything has been cleaned up before starting;
	// and this potentially requires all plugins should implement
	// stopping logic in a reentrant way (and not panic)
	_ = p.plugin.Stop()
	err := p.plugin.Start()
	if err != nil {
		return err
	}

	klog.Infof("starting to serve %s on %+v", p.Name(), p.sockets)

	if err := p.serve(); err != nil {
		return err
	}

	// reporter plugin relies on KubeCrane plugin manager automatic discovery, so needn't register by itself.
	klog.Infof("successfully started to serve %s on %+v", p.Name(), p.sockets)

	return nil
}

func (p *PluginRegistrationWrapper) stop() (stopErr error) {
	p.Lock()
	defer p.Unlock()

	// recover panic when stop plugin more than once
	defer func() {
		if err := recover(); err != nil {
			stopErr = fmt.Errorf("stop plugin %s failed with err: %v", p.Name(), err)
		}
	}()

	klog.Infof("stop to serve %s on %+v", p.Name(), p.sockets)

	for _, server := range p.servers {
		if server != nil {
			server.Stop()
		}
	}

	err := p.cleanup()
	if err != nil {
		return fmt.Errorf("cleanup failed for %s with error: %v", p.Name(), err)
	}

	err = p.plugin.Stop()
	if err != nil {
		return fmt.Errorf("stop wrapped plugin of %s failed with err: %v", p.Name(), err)
	}

	klog.Infof("stop reporter plugin %s successfully", p.Name())

	return nil
}

func (p *PluginRegistrationWrapper) cleanup() error {
	p.servers = nil

	for _, socket := range p.sockets {
		if err := os.Remove(socket); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("remove %s for %s failed with error: %v", socket, p.Name(), err)
		}
	}

	return nil
}

func (p *PluginRegistrationWrapper) serve() error {
	err := p.cleanup()
	if err != nil {
		return err
	}

	p.initialize()

	for _, socket := range p.sockets {
		curSocket := socket

		sock, err := net.Listen("unix", curSocket)
		if err != nil {
			return fmt.Errorf("listen for %s at socket: %s faield with err: %v", p.Name(), socket, err)
		}

		server := grpc.NewServer()
		// register reporter plugin server by real reporter Plugin which
		// only need to implement simple reporter plugin Start/Stop/Get/ListAndWatch
		// function without concerning plugin registration related logic
		p.serverRegister(server)
		watcherapi.RegisterRegistrationServer(server, p)

		go func() {
			lastCrashTime := time.Now()
			restartCount := 0
			for {
				klog.Infof("starting GRPC Server for %s at socket: %s", p.Name(), curSocket)
				err := server.Serve(sock)
				if err == nil {
					break
				}

				klog.Errorf("GRPC server for %s crashed with error: %v at socket: %s", p.Name(), err, curSocket)
				p.metricCallback("plugin_restart", 1)

				if restartCount > 5 {
					klog.Errorf("GRPC server for %s at socket: %s has repeatedly crashed recently. Quitting", p.Name(), curSocket)
					_ = p.Restart()
				}

				timeSinceLastCrash := time.Since(lastCrashTime).Seconds()
				lastCrashTime = time.Now()
				if timeSinceLastCrash > 3600 {
					restartCount = 1
				} else {
					restartCount++
				}
			}
		}()

		err = func() error {
			ctx, cancel := context.WithTimeout(context.Background(), grpcTimeout)
			defer cancel()

			conn, err := grpc.DialContext(ctx, curSocket,
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithBlock(),
				grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
					return (&net.Dialer{}).DialContext(ctx, "unix", addr)
				}),
			)

			if err != nil {
				return err
			}

			_ = conn.Close()

			return nil
		}()

		if err != nil {
			server.Stop()
			return fmt.Errorf("dial check for %s at socket: %s failed with err: %v", p.Name(), curSocket, err)
		}

		p.servers = append(p.servers, server)
		klog.Infof("serve %s successfully on %s", p.Name(), curSocket)
	}

	return nil
}
