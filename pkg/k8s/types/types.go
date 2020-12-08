package types

import (
	"fmt"
	"github.com/yametech/yamecloud/pkg/k8s"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"sort"
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

var _ sort.Interface = Resources{}

type Resources []k8s.Resource

func (r Resources) Len() int {
	return len(r)
}

func (r Resources) Less(i, j int) bool {
	return r[i].Name < r[j].Name
}

func (r Resources) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r Resources) Strings() []string {
	result := make([]string, 0)
	for _, item := range r {
		result = append(result, item.Name)
	}
	return result
}

func (r Resources) In(x string) bool {
	if sort.SearchStrings(r.Strings(), x) > len(r.Strings())-1 {
		return false
	}
	return true
}

type types struct {
	exclude Resources
	include Resources
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
		delete(m.data, value.Name)
	}

	for _type, value := range m.data {
		if m.include.In(_type) {
			continue
		}
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
