package service

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

var _ service.IResourceService = &Endpoint{}

type Endpoint struct {
	service.Interface
}

func NewEndpoint(svcInterface service.Interface) *Endpoint {
	endpoint := &Endpoint{Interface: svcInterface}
	svcInterface.Install(k8s.Endpoint, endpoint)
	return endpoint
}

func (g *Endpoint) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := g.Interface.Get(namespace, k8s.Endpoint, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (g *Endpoint) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := g.Interface.List(namespace, k8s.Endpoint, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (g *Endpoint) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	item, isUpdate, err := g.Interface.Apply(namespace, k8s.Endpoint, name, unstructuredExtend)
	if err != nil {
		return nil, isUpdate, err
	}
	return item, isUpdate, nil
}

func (g *Endpoint) Delete(namespace, name string) error {
	err := g.Interface.Delete(namespace, k8s.Endpoint, name)
	if err != nil {
		return err
	}
	return nil
}
