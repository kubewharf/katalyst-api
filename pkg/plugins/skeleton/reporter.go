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

	"k8s.io/klog/v2"

	"github.com/kubewharf/katalyst-api/pkg/protocol/reporterplugin/v1alpha1"
)

const (
	fakeReporterPluginName = "fake-reporter-plugin"
)

// ReporterPlugin performs report actions based on agent requirements.
type ReporterPlugin interface {
	GenericPlugin
	v1alpha1.ReporterPluginServer
}

// DummyReporterPlugin performs dummy report actions
type DummyReporterPlugin struct {
	v1alpha1.UnimplementedReporterPluginServer
}

// Name of the dummy reporter plugin
func (DummyReporterPlugin) Name() string { return fakeReporterPluginName }

// Start the dummy reporter plugin
func (DummyReporterPlugin) Start() error { return nil }

// Stop the dummy reporter plugin
func (DummyReporterPlugin) Stop() error { return nil }

var _ ReporterPlugin = &DummyReporterPlugin{
	v1alpha1.UnimplementedReporterPluginServer{},
}

// ReporterPluginStub is a stub for test reporter plugin manager
type ReporterPluginStub struct {
	started bool
	stop    chan struct{}
	update  chan []*v1alpha1.ReportContent

	content []*v1alpha1.ReportContent
	socket  string
	name    string

	v1alpha1.UnimplementedReporterPluginServer
}

var _ ReporterPlugin = &ReporterPluginStub{}

// NewReporterPluginStub initialize a reporter plugin stub which will report the input content.
func NewReporterPluginStub(content []*v1alpha1.ReportContent, name string) *ReporterPluginStub {
	return &ReporterPluginStub{
		started: false,
		content: content,
		name:    name,
		update:  make(chan []*v1alpha1.ReportContent),
	}
}

// Name of reporter plugin stub
func (m *ReporterPluginStub) Name() string {
	return m.name
}

// GetReportContent get report content from cache
func (m *ReporterPluginStub) GetReportContent(_ context.Context, _ *v1alpha1.Empty) (*v1alpha1.GetReportContentResponse, error) {
	return &v1alpha1.GetReportContentResponse{
		Content: m.content,
	}, nil
}

// ListAndWatchReportContent watch updateCh channel and send report content to plugin manager
func (m *ReporterPluginStub) ListAndWatchReportContent(empty *v1alpha1.Empty, server v1alpha1.ReporterPlugin_ListAndWatchReportContentServer) error {
	klog.Infof("plugin %s ListAndWatchReportContent", m.name)

	_ = server.Send(&v1alpha1.GetReportContentResponse{
		Content: m.content,
	})

	for {
		select {
		case <-m.stop:
			return nil
		case updated := <-m.update:
			err := server.Send(&v1alpha1.GetReportContentResponse{
				Content: updated,
			})
			if err != nil {
				klog.Errorf("plugin %s ListAndWatchReportContent send response failed, %v", m.name, err)
			}
		}
	}
}

// Update send report content to trigger list/watch
func (m *ReporterPluginStub) Update(content []*v1alpha1.ReportContent) {
	m.update <- content
}

// Start initialize some variables to start plugin
func (m *ReporterPluginStub) Start() (err error) {
	defer func() {
		if err == nil {
			m.started = true
		}
	}()

	if m.started {
		return
	}

	m.stop = make(chan struct{})
	return
}

// Stop clean some variables when stop plugin
func (m *ReporterPluginStub) Stop() error {
	defer func() {
		m.started = false
	}()

	if !m.started {
		return nil
	}

	close(m.stop)
	return nil
}
