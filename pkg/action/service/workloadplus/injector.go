package workloadplus

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

var _ service.IResourceService = &Injector{}

type Injector struct {
	service.Interface
}

func NewInjector(svcInterface service.Interface) *Injector {
	srv := &Injector{Interface: svcInterface}
	svcInterface.Install(k8s.Injector, srv)
	return srv
}

func (g *Injector) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := g.Interface.Get(namespace, k8s.Injector, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (g *Injector) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := g.Interface.List(namespace, k8s.Injector, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (g *Injector) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	item, isUpdate, err := g.Interface.Apply(namespace, k8s.Injector, name, unstructuredExtend)
	if err != nil {
		return nil, isUpdate, err
	}
	return item, isUpdate, nil
}

func (g *Injector) Delete(namespace, name string) error {
	err := g.Interface.Delete(namespace, k8s.Injector, name)
	if err != nil {
		return err
	}
	return nil
}
