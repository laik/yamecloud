package workloads

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

var _ service.IResourceService = &StorageClass{}

type StorageClass struct {
	service.Interface
}

func NewStorageClass(svcInterface service.Interface) *StorageClass {
	srv := &StorageClass{Interface: svcInterface}
	svcInterface.Install(k8s.StorageClass, srv)
	return srv
}

func (g *StorageClass) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := g.Interface.Get(namespace, k8s.StorageClass, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (g *StorageClass) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := g.Interface.List(namespace, k8s.StorageClass, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (g *StorageClass) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	item, isUpdate, err := g.Interface.Apply(namespace, k8s.StorageClass, name, unstructuredExtend)
	if err != nil {
		return nil, isUpdate, err
	}
	return item, isUpdate, nil
}

func (g *StorageClass) Delete(namespace, name string) error {
	err := g.Interface.Delete(namespace, k8s.StorageClass, name)
	if err != nil {
		return err
	}
	return nil
}
