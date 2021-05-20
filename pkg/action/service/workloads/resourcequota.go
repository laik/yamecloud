package workloads

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

var _ service.IResourceService = &ResourceQuota{}

type ResourceQuota struct {
	service.Interface
}

func NewResourceQuota(svcInterface service.Interface) *ResourceQuota {
	srv := &ResourceQuota{Interface: svcInterface}
	svcInterface.Install(k8s.ResourceQuota, srv)
	return srv
}

func (g *ResourceQuota) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := g.Interface.Get(namespace, k8s.ResourceQuota, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (g *ResourceQuota) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := g.Interface.List(namespace, k8s.ResourceQuota, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (g *ResourceQuota) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	item, isUpdate, err := g.Interface.Apply(namespace, k8s.ResourceQuota, name, unstructuredExtend)
	if err != nil {
		return nil, isUpdate, err
	}
	return item, isUpdate, nil
}

func (g *ResourceQuota) Delete(namespace, name string) error {
	err := g.Interface.Delete(namespace, k8s.ResourceQuota, name)
	if err != nil {
		return err
	}
	return nil
}
