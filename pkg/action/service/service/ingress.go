package service

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

var _ service.IResourceService = &Ingress{}

type Ingress struct {
	service.Interface
}

func NewIngress(svcInterface service.Interface) *Ingress {
	ingress := &Ingress{Interface: svcInterface}
	svcInterface.Install(k8s.Ingress, ingress)
	return ingress
}

func (g *Ingress) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := g.Interface.Get(namespace, k8s.Ingress, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (g *Ingress) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := g.Interface.List(namespace, k8s.Ingress, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (g *Ingress) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	item, isUpdate, err := g.Interface.Apply(namespace, k8s.Ingress, name, unstructuredExtend)
	if err != nil {
		return nil, isUpdate, err
	}
	return item, isUpdate, nil
}

func (g *Ingress) Delete(namespace, name string) error {
	err := g.Interface.Delete(namespace, k8s.Ingress, name)
	if err != nil {
		return err
	}
	return nil
}
