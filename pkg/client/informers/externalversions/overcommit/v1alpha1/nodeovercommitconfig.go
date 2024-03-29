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

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	time "time"

	overcommitv1alpha1 "github.com/kubewharf/katalyst-api/pkg/apis/overcommit/v1alpha1"
	versioned "github.com/kubewharf/katalyst-api/pkg/client/clientset/versioned"
	internalinterfaces "github.com/kubewharf/katalyst-api/pkg/client/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/kubewharf/katalyst-api/pkg/client/listers/overcommit/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// NodeOvercommitConfigInformer provides access to a shared informer and lister for
// NodeOvercommitConfigs.
type NodeOvercommitConfigInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.NodeOvercommitConfigLister
}

type nodeOvercommitConfigInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewNodeOvercommitConfigInformer constructs a new informer for NodeOvercommitConfig type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewNodeOvercommitConfigInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredNodeOvercommitConfigInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredNodeOvercommitConfigInformer constructs a new informer for NodeOvercommitConfig type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredNodeOvercommitConfigInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.OvercommitV1alpha1().NodeOvercommitConfigs().List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.OvercommitV1alpha1().NodeOvercommitConfigs().Watch(context.TODO(), options)
			},
		},
		&overcommitv1alpha1.NodeOvercommitConfig{},
		resyncPeriod,
		indexers,
	)
}

func (f *nodeOvercommitConfigInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredNodeOvercommitConfigInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *nodeOvercommitConfigInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&overcommitv1alpha1.NodeOvercommitConfig{}, f.defaultInformer)
}

func (f *nodeOvercommitConfigInformer) Lister() v1alpha1.NodeOvercommitConfigLister {
	return v1alpha1.NewNodeOvercommitConfigLister(f.Informer().GetIndexer())
}
