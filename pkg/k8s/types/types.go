package types

import (
	"fmt"
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
	excluded []string
	Data     map[string]schema.GroupVersionResource
}

func NewResources(excluded []string) *types {
	rs := &types{
		excluded: excluded,
		Data:     make(map[string]schema.GroupVersionResource),
	}
	return rs
}

func (m *types) Register(resource Resource) { m.register(resource.Name, resource.Schema) }

func (m *types) register(s string, resource schema.GroupVersionResource) {
	if _, exist := m.Data[s]; exist {
		return
	}
	m.Data[s] = resource
}

func (m *types) Ranges(d dynamicinformer.DynamicSharedInformerFactory, stop <-chan struct{}) {
	for _, v := range m.excluded {
		value := v
		delete(m.Data, value)
	}
	for _, v := range m.Data {
		value := v
		go d.ForResource(value).Informer().Run(stop)
	}
}

func (m *types) GroupVersionResource(s string) (schema.GroupVersionResource, error) {
	item, exist := m.Data[s]
	if !exist {
		return schema.GroupVersionResource{}, fmt.Errorf("resource (%s) not exist", s)
	}
	return item, nil
}
