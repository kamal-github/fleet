/*
Copyright 2023 Rancher Labs, Inc.

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

// Code generated by main. DO NOT EDIT.

package v1alpha1

import (
	"context"
	"time"

	v1alpha1 "github.com/rancher/fleet/pkg/apis/fleet.cattle.io/v1alpha1"
	"github.com/rancher/lasso/pkg/client"
	"github.com/rancher/lasso/pkg/controller"
	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/condition"
	"github.com/rancher/wrangler/pkg/generic"
	"github.com/rancher/wrangler/pkg/kv"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type ImageScanHandler func(string, *v1alpha1.ImageScan) (*v1alpha1.ImageScan, error)

type ImageScanController interface {
	generic.ControllerMeta
	ImageScanClient

	OnChange(ctx context.Context, name string, sync ImageScanHandler)
	OnRemove(ctx context.Context, name string, sync ImageScanHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() ImageScanCache
}

type ImageScanClient interface {
	Create(*v1alpha1.ImageScan) (*v1alpha1.ImageScan, error)
	Update(*v1alpha1.ImageScan) (*v1alpha1.ImageScan, error)
	UpdateStatus(*v1alpha1.ImageScan) (*v1alpha1.ImageScan, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v1alpha1.ImageScan, error)
	List(namespace string, opts metav1.ListOptions) (*v1alpha1.ImageScanList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.ImageScan, err error)
}

type ImageScanCache interface {
	Get(namespace, name string) (*v1alpha1.ImageScan, error)
	List(namespace string, selector labels.Selector) ([]*v1alpha1.ImageScan, error)

	AddIndexer(indexName string, indexer ImageScanIndexer)
	GetByIndex(indexName, key string) ([]*v1alpha1.ImageScan, error)
}

type ImageScanIndexer func(obj *v1alpha1.ImageScan) ([]string, error)

type imageScanController struct {
	controller    controller.SharedController
	client        *client.Client
	gvk           schema.GroupVersionKind
	groupResource schema.GroupResource
}

func NewImageScanController(gvk schema.GroupVersionKind, resource string, namespaced bool, controller controller.SharedControllerFactory) ImageScanController {
	c := controller.ForResourceKind(gvk.GroupVersion().WithResource(resource), gvk.Kind, namespaced)
	return &imageScanController{
		controller: c,
		client:     c.Client(),
		gvk:        gvk,
		groupResource: schema.GroupResource{
			Group:    gvk.Group,
			Resource: resource,
		},
	}
}

func FromImageScanHandlerToHandler(sync ImageScanHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1alpha1.ImageScan
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1alpha1.ImageScan))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *imageScanController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1alpha1.ImageScan))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateImageScanDeepCopyOnChange(client ImageScanClient, obj *v1alpha1.ImageScan, handler func(obj *v1alpha1.ImageScan) (*v1alpha1.ImageScan, error)) (*v1alpha1.ImageScan, error) {
	if obj == nil {
		return obj, nil
	}

	copyObj := obj.DeepCopy()
	newObj, err := handler(copyObj)
	if newObj != nil {
		copyObj = newObj
	}
	if obj.ResourceVersion == copyObj.ResourceVersion && !equality.Semantic.DeepEqual(obj, copyObj) {
		return client.Update(copyObj)
	}

	return copyObj, err
}

func (c *imageScanController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controller.RegisterHandler(ctx, name, controller.SharedControllerHandlerFunc(handler))
}

func (c *imageScanController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), handler))
}

func (c *imageScanController) OnChange(ctx context.Context, name string, sync ImageScanHandler) {
	c.AddGenericHandler(ctx, name, FromImageScanHandlerToHandler(sync))
}

func (c *imageScanController) OnRemove(ctx context.Context, name string, sync ImageScanHandler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), FromImageScanHandlerToHandler(sync)))
}

func (c *imageScanController) Enqueue(namespace, name string) {
	c.controller.Enqueue(namespace, name)
}

func (c *imageScanController) EnqueueAfter(namespace, name string, duration time.Duration) {
	c.controller.EnqueueAfter(namespace, name, duration)
}

func (c *imageScanController) Informer() cache.SharedIndexInformer {
	return c.controller.Informer()
}

func (c *imageScanController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *imageScanController) Cache() ImageScanCache {
	return &imageScanCache{
		indexer:  c.Informer().GetIndexer(),
		resource: c.groupResource,
	}
}

func (c *imageScanController) Create(obj *v1alpha1.ImageScan) (*v1alpha1.ImageScan, error) {
	result := &v1alpha1.ImageScan{}
	return result, c.client.Create(context.TODO(), obj.Namespace, obj, result, metav1.CreateOptions{})
}

func (c *imageScanController) Update(obj *v1alpha1.ImageScan) (*v1alpha1.ImageScan, error) {
	result := &v1alpha1.ImageScan{}
	return result, c.client.Update(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *imageScanController) UpdateStatus(obj *v1alpha1.ImageScan) (*v1alpha1.ImageScan, error) {
	result := &v1alpha1.ImageScan{}
	return result, c.client.UpdateStatus(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *imageScanController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.client.Delete(context.TODO(), namespace, name, *options)
}

func (c *imageScanController) Get(namespace, name string, options metav1.GetOptions) (*v1alpha1.ImageScan, error) {
	result := &v1alpha1.ImageScan{}
	return result, c.client.Get(context.TODO(), namespace, name, result, options)
}

func (c *imageScanController) List(namespace string, opts metav1.ListOptions) (*v1alpha1.ImageScanList, error) {
	result := &v1alpha1.ImageScanList{}
	return result, c.client.List(context.TODO(), namespace, result, opts)
}

func (c *imageScanController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.client.Watch(context.TODO(), namespace, opts)
}

func (c *imageScanController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (*v1alpha1.ImageScan, error) {
	result := &v1alpha1.ImageScan{}
	return result, c.client.Patch(context.TODO(), namespace, name, pt, data, result, metav1.PatchOptions{}, subresources...)
}

type imageScanCache struct {
	indexer  cache.Indexer
	resource schema.GroupResource
}

func (c *imageScanCache) Get(namespace, name string) (*v1alpha1.ImageScan, error) {
	obj, exists, err := c.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(c.resource, name)
	}
	return obj.(*v1alpha1.ImageScan), nil
}

func (c *imageScanCache) List(namespace string, selector labels.Selector) (ret []*v1alpha1.ImageScan, err error) {

	err = cache.ListAllByNamespace(c.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.ImageScan))
	})

	return ret, err
}

func (c *imageScanCache) AddIndexer(indexName string, indexer ImageScanIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1alpha1.ImageScan))
		},
	}))
}

func (c *imageScanCache) GetByIndex(indexName, key string) (result []*v1alpha1.ImageScan, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v1alpha1.ImageScan, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v1alpha1.ImageScan))
	}
	return result, nil
}

type ImageScanStatusHandler func(obj *v1alpha1.ImageScan, status v1alpha1.ImageScanStatus) (v1alpha1.ImageScanStatus, error)

type ImageScanGeneratingHandler func(obj *v1alpha1.ImageScan, status v1alpha1.ImageScanStatus) ([]runtime.Object, v1alpha1.ImageScanStatus, error)

func RegisterImageScanStatusHandler(ctx context.Context, controller ImageScanController, condition condition.Cond, name string, handler ImageScanStatusHandler) {
	statusHandler := &imageScanStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, FromImageScanHandlerToHandler(statusHandler.sync))
}

func RegisterImageScanGeneratingHandler(ctx context.Context, controller ImageScanController, apply apply.Apply,
	condition condition.Cond, name string, handler ImageScanGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &imageScanGeneratingHandler{
		ImageScanGeneratingHandler: handler,
		apply:                      apply,
		name:                       name,
		gvk:                        controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	controller.OnChange(ctx, name, statusHandler.Remove)
	RegisterImageScanStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type imageScanStatusHandler struct {
	client    ImageScanClient
	condition condition.Cond
	handler   ImageScanStatusHandler
}

func (a *imageScanStatusHandler) sync(key string, obj *v1alpha1.ImageScan) (*v1alpha1.ImageScan, error) {
	if obj == nil {
		return obj, nil
	}

	origStatus := obj.Status.DeepCopy()
	obj = obj.DeepCopy()
	newStatus, err := a.handler(obj, obj.Status)
	if err != nil {
		// Revert to old status on error
		newStatus = *origStatus.DeepCopy()
	}

	if a.condition != "" {
		if errors.IsConflict(err) {
			a.condition.SetError(&newStatus, "", nil)
		} else {
			a.condition.SetError(&newStatus, "", err)
		}
	}
	if !equality.Semantic.DeepEqual(origStatus, &newStatus) {
		if a.condition != "" {
			// Since status has changed, update the lastUpdatedTime
			a.condition.LastUpdated(&newStatus, time.Now().UTC().Format(time.RFC3339))
		}

		var newErr error
		obj.Status = newStatus
		newObj, newErr := a.client.UpdateStatus(obj)
		if err == nil {
			err = newErr
		}
		if newErr == nil {
			obj = newObj
		}
	}
	return obj, err
}

type imageScanGeneratingHandler struct {
	ImageScanGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
}

func (a *imageScanGeneratingHandler) Remove(key string, obj *v1alpha1.ImageScan) (*v1alpha1.ImageScan, error) {
	if obj != nil {
		return obj, nil
	}

	obj = &v1alpha1.ImageScan{}
	obj.Namespace, obj.Name = kv.RSplit(key, "/")
	obj.SetGroupVersionKind(a.gvk)

	return nil, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects()
}

func (a *imageScanGeneratingHandler) Handle(obj *v1alpha1.ImageScan, status v1alpha1.ImageScanStatus) (v1alpha1.ImageScanStatus, error) {
	if !obj.DeletionTimestamp.IsZero() {
		return status, nil
	}

	objs, newStatus, err := a.ImageScanGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}

	return newStatus, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
}
