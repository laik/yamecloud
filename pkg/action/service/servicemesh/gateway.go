package servicemesh

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

var _ service.IResourceService = &Gateway{}

type Gateway struct {
	service.Interface
}

func NewGateway(svcInterface service.Interface) *Gateway {
	gateway := &Gateway{Interface: svcInterface}
	svcInterface.Install(k8s.Gateway, gateway)
	return gateway
}

func (g *Gateway) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := g.Interface.Get(namespace, k8s.Gateway, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (g *Gateway) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := g.Interface.List(namespace, k8s.Gateway, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (g *Gateway) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	item, isUpdate, err := g.Interface.Apply(namespace, k8s.Gateway, name, unstructuredExtend)
	if err != nil {
		return nil, isUpdate, err
	}
	return item, isUpdate, nil
}

func (g *Gateway) Delete(namespace, name string) error {
	err := g.Interface.Delete(namespace, k8s.Gateway, name)
	if err != nil {
		return err
	}
	return nil
}
