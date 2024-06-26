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

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha2

import (
	v1alpha2 "github.com/kubewharf/katalyst-api/pkg/apis/autoscaling/v1alpha2"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// VirtualWorkloadLister helps list VirtualWorkloads.
// All objects returned here must be treated as read-only.
type VirtualWorkloadLister interface {
	// List lists all VirtualWorkloads in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha2.VirtualWorkload, err error)
	// VirtualWorkloads returns an object that can list and get VirtualWorkloads.
	VirtualWorkloads(namespace string) VirtualWorkloadNamespaceLister
	VirtualWorkloadListerExpansion
}

// virtualWorkloadLister implements the VirtualWorkloadLister interface.
type virtualWorkloadLister struct {
	indexer cache.Indexer
}

// NewVirtualWorkloadLister returns a new VirtualWorkloadLister.
func NewVirtualWorkloadLister(indexer cache.Indexer) VirtualWorkloadLister {
	return &virtualWorkloadLister{indexer: indexer}
}

// List lists all VirtualWorkloads in the indexer.
func (s *virtualWorkloadLister) List(selector labels.Selector) (ret []*v1alpha2.VirtualWorkload, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha2.VirtualWorkload))
	})
	return ret, err
}

// VirtualWorkloads returns an object that can list and get VirtualWorkloads.
func (s *virtualWorkloadLister) VirtualWorkloads(namespace string) VirtualWorkloadNamespaceLister {
	return virtualWorkloadNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// VirtualWorkloadNamespaceLister helps list and get VirtualWorkloads.
// All objects returned here must be treated as read-only.
type VirtualWorkloadNamespaceLister interface {
	// List lists all VirtualWorkloads in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha2.VirtualWorkload, err error)
	// Get retrieves the VirtualWorkload from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha2.VirtualWorkload, error)
	VirtualWorkloadNamespaceListerExpansion
}

// virtualWorkloadNamespaceLister implements the VirtualWorkloadNamespaceLister
// interface.
type virtualWorkloadNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all VirtualWorkloads in the indexer for a given namespace.
func (s virtualWorkloadNamespaceLister) List(selector labels.Selector) (ret []*v1alpha2.VirtualWorkload, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha2.VirtualWorkload))
	})
	return ret, err
}

// Get retrieves the VirtualWorkload from the indexer for a given namespace and name.
func (s virtualWorkloadNamespaceLister) Get(name string) (*v1alpha2.VirtualWorkload, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha2.Resource("virtualworkload"), name)
	}
	return obj.(*v1alpha2.VirtualWorkload), nil
}
