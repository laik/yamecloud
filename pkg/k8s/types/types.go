package types

import (
	"fmt"
	"github.com/yametech/yamecloud/pkg/k8s"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic/dynamicinformer"
)

var _ k8s.ITypes = &types{}

func NewResourceITypes(include []k8s.Resource) k8s.ITypes {
	return newTypes(include, nil)
}

func newTypes(include, exclude []k8s.Resource) *types {
	_types := &types{
		include: include,
		exclude: exclude,
		data:    make(map[k8s.ResourceType]schema.GroupVersionResource),
	}
	_types.Register(include...)
	return _types
}

type types struct {
	exclude []k8s.Resource
	include []k8s.Resource
	data    map[k8s.ResourceType]schema.GroupVersionResource
}

func (m *types) Register(resources ...k8s.Resource) {
	for _, resource := range resources {
		m.register(resource.Name, resource.Schema)
	}
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
		delete(m.data, k8s.ResourceType(value.Name))
	}
	for _, v := range m.data {
		value := v
		go d.ForResource(value).Informer().Run(stop)
	}
}

func (m *types) GroupVersionResource(s k8s.ResourceType) (schema.GroupVersionResource, error) {
	item, exist := m.data[s]
	if !exist {
		return schema.GroupVersionResource{}, fmt.Errorf("resource (%s) not exist", s)
	}
	return item, nil
}
