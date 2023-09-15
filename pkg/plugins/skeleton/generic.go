//go:build !windows
// +build !windows

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
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog/v2"
	watcherapi "k8s.io/kubelet/pkg/apis/pluginregistration/v1"
	qrmpluginapi "k8s.io/kubelet/pkg/apis/resourceplugin/v1alpha1"
	utilfs "k8s.io/kubernetes/pkg/util/filesystem"

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

	metricCallback MetricCallback
	serverRegister pluginServerRegister

	GenericPlugin
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
		metricCallback:     metricCallback,
		serverRegister:     serverRegister,
		GenericPlugin:      plugin,
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

// Start the plugin with auto restart logic, besides
// it will register to plugin grpc server
func (p *PluginRegistrationWrapper) Start() error {
	go func() {
		for {
			_ = wait.PollImmediateInfinite(restartRetryInterval, func() (bool, error) {

				select {
				case <-p.stopCh:
					return false, fmt.Errorf("stop channel closed during polling start")
				case _, ok := <-p.restartCh:
					if !ok {
						return false, fmt.Errorf("restart channel closed during polling start")
					}

					klog.Infof("receive restart signal during polling start, continue to poll")
					return false, nil
				default:
				}

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

func (p *PluginRegistrationWrapper) initializeSocketDirs() error {
	filteredSockets := make([]string, 0, len(p.sockets))
	inodeSet := make(map[uint64]string)
	fs := utilfs.DefaultFs{}

	for _, socket := range p.sockets {
		socketDir := path.Dir(socket)

		if _, err := fs.Stat(socketDir); err != nil {
			// MkdirAll returns nil if directory already exists.
			if err := fs.MkdirAll(socketDir, 0755); err != nil {
				return fmt.Errorf("create socket dir: %s failed with error: %v", socketDir, err)
			}
		}

		var socketDirStat syscall.Stat_t
		if err := syscall.Stat(socketDir, &socketDirStat); err != nil {
			return fmt.Errorf("stat socket dir: %s failed with error: %v", socketDir, err)
		}

		klog.Infof("detect socketDir: %s with inode: %d", socketDir, socketDirStat.Ino)

		if prevDir, found := inodeSet[socketDirStat.Ino]; found {
			klog.Warningf("found socket dir: %s duplicated with %s, same inode: %d",
				socketDir, prevDir, socketDirStat.Ino)
			continue
		}

		inodeSet[socketDirStat.Ino] = socketDir
		filteredSockets = append(filteredSockets, socket)
	}

	if len(filteredSockets) == 0 {
		return fmt.Errorf("filteredSockets is empty")
	}

	p.sockets = filteredSockets
	return nil
}

// initialize plugin servers
func (p *PluginRegistrationWrapper) initialize() error {
	err := p.initializeSocketDirs()
	if err != nil {
		return fmt.Errorf("initializeSocketDirs failed with error: %v", err)
	}

	p.servers = make([]*grpc.Server, 0, len(p.sockets))
	return nil
}

func (p *PluginRegistrationWrapper) start() (startError error) {
	p.Lock()
	defer func() {
		p.Unlock()
		if startError != nil {
			// call stop to revert executed start steps for all servers
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
	_ = p.GenericPlugin.Stop()
	err := p.GenericPlugin.Start()
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

	err = p.GenericPlugin.Stop()
	if err != nil {
		return fmt.Errorf("stop wrapped plugin of %s failed with err: %v", p.Name(), err)
	}

	klog.Infof("stop plugin %s successfully", p.Name())

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

	if err := p.initialize(); err != nil {
		return fmt.Errorf("initialize failed %v", err)
	}

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

		// server.Serve works in a separate goroutine; we will retry several times if
		// it crashes, and trigger restart if it exceeds retry thresholds.
		// if the server stops, server.Serve will return will nil error, so the for loop
		// can be break successfully without causing goroutine leaks.
		go func() {
			lastCrashTime := time.Now()
			restartCount := 0
			for {
				klog.Infof("starting GRPC Server for %s at socket: %s", p.Name(), curSocket)
				err := server.Serve(sock)
				if err == nil {
					klog.Infof("GRPC Server for %s at socket: %s stops serving", p.Name(), curSocket)
					break
				} else if err == grpc.ErrServerStopped {
					klog.Infof("GRPC Server for %s at socket: %s already stopped, break from serving goroutine", p.Name(), curSocket)
					break
				}

				klog.Errorf("GRPC server for %s crashed with error: %v at socket: %s", p.Name(), err, curSocket)
				p.metricCallback("plugin_restart", 1)

				if restartCount > 5 {
					klog.Errorf("GRPC server for %s at socket: %s has repeatedly crashed recently. Quitting", p.Name(), curSocket)
					_ = p.Restart()
					break
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

		// try to connect with the server to ensure the serving works as expected
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
