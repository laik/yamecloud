package servicemesh

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

var _ service.IResourceService = &VirtualService{}

type VirtualService struct {
	service.Interface
}

func NewVirtualService(svcInterface service.Interface) *VirtualService {
	virtualService := &VirtualService{Interface: svcInterface}
	svcInterface.Install(k8s.VirtualService, virtualService)
	return virtualService
}

func (v *VirtualService) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := v.Interface.Get(namespace, k8s.VirtualService, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (v *VirtualService) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := v.Interface.List(namespace, k8s.VirtualService, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (v *VirtualService) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	item, isUpdate, err := v.Interface.Apply(namespace, k8s.VirtualService, name, unstructuredExtend)
	if err != nil {
		return nil, isUpdate, err
	}
	return item, isUpdate, nil
}

func (v *VirtualService) Delete(namespace, name string) error {
	err := v.Interface.Delete(namespace, k8s.VirtualService, name)
	if err != nil {
		return err
	}
	return nil
}
