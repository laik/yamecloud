package workloads

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

var _ service.IResourceService = &Namespace{}

type Namespace struct {
	service.Interface
}

func NewNamespace(svcInterface service.Interface) *Namespace {
	srv := &Namespace{Interface: svcInterface}
	svcInterface.Install(k8s.Namespace, srv)
	return srv
}

func (g *Namespace) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := g.Interface.Get(namespace, k8s.Namespace, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (g *Namespace) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := g.Interface.List(namespace, k8s.Namespace, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (g *Namespace) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	item, isUpdate, err := g.Interface.Apply(namespace, k8s.Namespace, name, unstructuredExtend)
	if err != nil {
		return nil, isUpdate, err
	}
	return item, isUpdate, nil
}

func (g *Namespace) Delete(namespace, name string) error {
	err := g.Interface.Delete(namespace, k8s.Namespace, name)
	if err != nil {
		return err
	}
	return nil
}
