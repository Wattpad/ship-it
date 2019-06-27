/*
Copyright The Kubernetes Authors.

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

package v1alpha1

import (
	"time"

	scheme "ship-it/pkg/generated/clientset/versioned/scheme"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
	v1alpha1 "ship-it/pkg/apis/helmreleases.k8s.wattpad.com/v1alpha1"
)

// HelmReleasesGetter has a method to return a HelmReleaseInterface.
// A group's client should implement this interface.
type HelmReleasesGetter interface {
	HelmReleases(namespace string) HelmReleaseInterface
}

// HelmReleaseInterface has methods to work with HelmRelease resources.
type HelmReleaseInterface interface {
	Create(*v1alpha1.HelmRelease) (*v1alpha1.HelmRelease, error)
	Update(*v1alpha1.HelmRelease) (*v1alpha1.HelmRelease, error)
	UpdateStatus(*v1alpha1.HelmRelease) (*v1alpha1.HelmRelease, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.HelmRelease, error)
	List(opts v1.ListOptions) (*v1alpha1.HelmReleaseList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.HelmRelease, err error)
	HelmReleaseExpansion
}

// helmReleases implements HelmReleaseInterface
type helmReleases struct {
	client rest.Interface
	ns     string
}

// newHelmReleases returns a HelmReleases
func newHelmReleases(c *HelmreleasesV1alpha1Client, namespace string) *helmReleases {
	return &helmReleases{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the helmRelease, and returns the corresponding helmRelease object, and an error if there is any.
func (c *helmReleases) Get(name string, options v1.GetOptions) (result *v1alpha1.HelmRelease, err error) {
	result = &v1alpha1.HelmRelease{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("helmreleases").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of HelmReleases that match those selectors.
func (c *helmReleases) List(opts v1.ListOptions) (result *v1alpha1.HelmReleaseList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.HelmReleaseList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("helmreleases").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested helmReleases.
func (c *helmReleases) Watch(opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("helmreleases").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a helmRelease and creates it.  Returns the server's representation of the helmRelease, and an error, if there is any.
func (c *helmReleases) Create(helmRelease *v1alpha1.HelmRelease) (result *v1alpha1.HelmRelease, err error) {
	result = &v1alpha1.HelmRelease{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("helmreleases").
		Body(helmRelease).
		Do().
		Into(result)
	return
}

// Update takes the representation of a helmRelease and updates it. Returns the server's representation of the helmRelease, and an error, if there is any.
func (c *helmReleases) Update(helmRelease *v1alpha1.HelmRelease) (result *v1alpha1.HelmRelease, err error) {
	result = &v1alpha1.HelmRelease{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("helmreleases").
		Name(helmRelease.Name).
		Body(helmRelease).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *helmReleases) UpdateStatus(helmRelease *v1alpha1.HelmRelease) (result *v1alpha1.HelmRelease, err error) {
	result = &v1alpha1.HelmRelease{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("helmreleases").
		Name(helmRelease.Name).
		SubResource("status").
		Body(helmRelease).
		Do().
		Into(result)
	return
}

// Delete takes name of the helmRelease and deletes it. Returns an error if one occurs.
func (c *helmReleases) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("helmreleases").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *helmReleases) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("helmreleases").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched helmRelease.
func (c *helmReleases) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.HelmRelease, err error) {
	result = &v1alpha1.HelmRelease{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("helmreleases").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
