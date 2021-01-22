package service

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

var _ service.IResourceService = &Service{}

type Service struct {
	service.Interface
}

func NewService(svcInterface service.Interface) *Service {
	service := &Service{Interface: svcInterface}
	svcInterface.Install(k8s.Service, service)
	return service
}

func (g *Service) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := g.Interface.Get(namespace, k8s.Service, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (g *Service) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := g.Interface.List(namespace, k8s.Service, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (g *Service) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	item, isUpdate, err := g.Interface.Apply(namespace, k8s.Service, name, unstructuredExtend)
	if err != nil {
		return nil, isUpdate, err
	}
	return item, isUpdate, nil
}

func (g *Service) Delete(namespace, name string) error {
	err := g.Interface.Delete(namespace, k8s.Service, name)
	if err != nil {
		return err
	}
	return nil
}
