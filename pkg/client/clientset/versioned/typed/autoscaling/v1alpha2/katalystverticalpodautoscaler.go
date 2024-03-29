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

// Code generated by client-gen. DO NOT EDIT.

package v1alpha2

import (
	"context"
	"time"

	v1alpha2 "github.com/kubewharf/katalyst-api/pkg/apis/autoscaling/v1alpha2"
	scheme "github.com/kubewharf/katalyst-api/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// KatalystVerticalPodAutoscalersGetter has a method to return a KatalystVerticalPodAutoscalerInterface.
// A group's client should implement this interface.
type KatalystVerticalPodAutoscalersGetter interface {
	KatalystVerticalPodAutoscalers(namespace string) KatalystVerticalPodAutoscalerInterface
}

// KatalystVerticalPodAutoscalerInterface has methods to work with KatalystVerticalPodAutoscaler resources.
type KatalystVerticalPodAutoscalerInterface interface {
	Create(ctx context.Context, katalystVerticalPodAutoscaler *v1alpha2.KatalystVerticalPodAutoscaler, opts v1.CreateOptions) (*v1alpha2.KatalystVerticalPodAutoscaler, error)
	Update(ctx context.Context, katalystVerticalPodAutoscaler *v1alpha2.KatalystVerticalPodAutoscaler, opts v1.UpdateOptions) (*v1alpha2.KatalystVerticalPodAutoscaler, error)
	UpdateStatus(ctx context.Context, katalystVerticalPodAutoscaler *v1alpha2.KatalystVerticalPodAutoscaler, opts v1.UpdateOptions) (*v1alpha2.KatalystVerticalPodAutoscaler, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha2.KatalystVerticalPodAutoscaler, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha2.KatalystVerticalPodAutoscalerList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha2.KatalystVerticalPodAutoscaler, err error)
	KatalystVerticalPodAutoscalerExpansion
}

// katalystVerticalPodAutoscalers implements KatalystVerticalPodAutoscalerInterface
type katalystVerticalPodAutoscalers struct {
	client rest.Interface
	ns     string
}

// newKatalystVerticalPodAutoscalers returns a KatalystVerticalPodAutoscalers
func newKatalystVerticalPodAutoscalers(c *AutoscalingV1alpha2Client, namespace string) *katalystVerticalPodAutoscalers {
	return &katalystVerticalPodAutoscalers{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the katalystVerticalPodAutoscaler, and returns the corresponding katalystVerticalPodAutoscaler object, and an error if there is any.
func (c *katalystVerticalPodAutoscalers) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha2.KatalystVerticalPodAutoscaler, err error) {
	result = &v1alpha2.KatalystVerticalPodAutoscaler{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("katalystverticalpodautoscalers").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of KatalystVerticalPodAutoscalers that match those selectors.
func (c *katalystVerticalPodAutoscalers) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha2.KatalystVerticalPodAutoscalerList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha2.KatalystVerticalPodAutoscalerList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("katalystverticalpodautoscalers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested katalystVerticalPodAutoscalers.
func (c *katalystVerticalPodAutoscalers) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("katalystverticalpodautoscalers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a katalystVerticalPodAutoscaler and creates it.  Returns the server's representation of the katalystVerticalPodAutoscaler, and an error, if there is any.
func (c *katalystVerticalPodAutoscalers) Create(ctx context.Context, katalystVerticalPodAutoscaler *v1alpha2.KatalystVerticalPodAutoscaler, opts v1.CreateOptions) (result *v1alpha2.KatalystVerticalPodAutoscaler, err error) {
	result = &v1alpha2.KatalystVerticalPodAutoscaler{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("katalystverticalpodautoscalers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(katalystVerticalPodAutoscaler).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a katalystVerticalPodAutoscaler and updates it. Returns the server's representation of the katalystVerticalPodAutoscaler, and an error, if there is any.
func (c *katalystVerticalPodAutoscalers) Update(ctx context.Context, katalystVerticalPodAutoscaler *v1alpha2.KatalystVerticalPodAutoscaler, opts v1.UpdateOptions) (result *v1alpha2.KatalystVerticalPodAutoscaler, err error) {
	result = &v1alpha2.KatalystVerticalPodAutoscaler{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("katalystverticalpodautoscalers").
		Name(katalystVerticalPodAutoscaler.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(katalystVerticalPodAutoscaler).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *katalystVerticalPodAutoscalers) UpdateStatus(ctx context.Context, katalystVerticalPodAutoscaler *v1alpha2.KatalystVerticalPodAutoscaler, opts v1.UpdateOptions) (result *v1alpha2.KatalystVerticalPodAutoscaler, err error) {
	result = &v1alpha2.KatalystVerticalPodAutoscaler{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("katalystverticalpodautoscalers").
		Name(katalystVerticalPodAutoscaler.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(katalystVerticalPodAutoscaler).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the katalystVerticalPodAutoscaler and deletes it. Returns an error if one occurs.
func (c *katalystVerticalPodAutoscalers) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("katalystverticalpodautoscalers").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *katalystVerticalPodAutoscalers) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("katalystverticalpodautoscalers").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched katalystVerticalPodAutoscaler.
func (c *katalystVerticalPodAutoscalers) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha2.KatalystVerticalPodAutoscaler, err error) {
	result = &v1alpha2.KatalystVerticalPodAutoscaler{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("katalystverticalpodautoscalers").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
