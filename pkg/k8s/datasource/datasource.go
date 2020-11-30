package datasource

import (
	"github.com/yametech/yamecloud/pkg/configure"
	"github.com/yametech/yamecloud/pkg/k8s"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/watch"
)

var _ k8s.Interface = &dataSource{}

func NewInterface(configure *configure.InstallConfigure) k8s.Interface {
	return &dataSource{configure: configure}
}

type dataSource struct {
	configure *configure.InstallConfigure
}

func (d *dataSource) XGet(namespace string, resourceType k8s.ResourceType, name string) (*unstructured.Unstructured, error) {
	panic("implement me")
}

func (d *dataSource) Watch(namespace string, resourceType k8s.ResourceType, resourceVersion string, selector interface{}) (<-chan watch.Event, error) {
	panic("implement me")
}

func (d *dataSource) Apply(namespace string, resourceType k8s.ResourceType, name string, unstructured *unstructured.Unstructured, forceUpdate bool) (newUnstructured *unstructured.Unstructured, isUpdate bool, err error) {
	panic("implement me")
}

func (d *dataSource) Delete(namespace string, resourceType k8s.ResourceType, name string) error {
	panic("implement me")
}

func (d *dataSource) Patch(namespace string, resourceType k8s.ResourceType, name string, path string, data interface{}) error {
	panic("implement me")
}

func (d *dataSource) List(namespace string, resourceType k8s.ResourceType, selector interface{}) (*unstructured.UnstructuredList, error) {
	panic("implement me")
}

func (d *dataSource) Get(namespace string, resourceType k8s.ResourceType, name string) (*unstructured.Unstructured, error) {
	panic("implement me")
}

func (d *dataSource) Cache() k8s.ICache {
	panic("implement me")
}
