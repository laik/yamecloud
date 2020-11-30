package types

import (
	"fmt"
	"github.com/yametech/yamecloud/pkg/k8s"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic/dynamicinformer"
)

type ResourceRegister interface {
	Register(Resource)
}

type ResourceLister interface {
	Ranges(d dynamicinformer.DynamicSharedInformerFactory, stop <-chan struct{})
	GroupVersionResource(string) (schema.GroupVersionResource, error)
}

type ITypes interface {
	ResourceRegister
	ResourceLister
}

var _ ITypes = &types{}

type Resource struct {
	Name   string
	Schema schema.GroupVersionResource
}

type types struct {
	exclude []k8s.ResourceType
	include []k8s.ResourceType
	data    map[k8s.ResourceType]schema.GroupVersionResource
}

func NewResourceLister(include ...k8s.ResourceType) ResourceLister {
	return newTypes(include, nil)
}

func newTypes(include, exclude []k8s.ResourceType) *types {
	_types := &types{
		include: include,
		exclude: exclude,
		data:    make(map[k8s.ResourceType]schema.GroupVersionResource),
	}
	return _types
}

func (m *types) Register(resource Resource) {
	m.register(k8s.ResourceType(resource.Name), resource.Schema)
}

func (m *types) register(s k8s.ResourceType, resource schema.GroupVersionResource) {
	if _, exist := m.data[s]; exist {
		return
	}
	m.data[s] = resource
}

func (m *types) Ranges(d dynamicinformer.DynamicSharedInformerFactory, stop <-chan struct{}) {
	for _, v := range m.exclude {
		value := v
		delete(m.data, value)
	}
	for _, v := range m.data {
		value := v
		go d.ForResource(value).Informer().Run(stop)
	}
}

func (m *types) GroupVersionResource(s string) (schema.GroupVersionResource, error) {
	item, exist := m.data[k8s.ResourceType(s)]
	if !exist {
		return schema.GroupVersionResource{}, fmt.Errorf("resource (%s) not exist", s)
	}
	return item, nil
}
