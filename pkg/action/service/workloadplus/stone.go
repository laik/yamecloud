package workloadplus

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

var _ service.IResourceService = &Stone{}

type Stone struct {
	service.Interface
}

func NewStone(svcInterface service.Interface) *Stone {
	srv := &Stone{Interface: svcInterface}
	svcInterface.Install(k8s.Stone, srv)
	return srv
}

func (g *Stone) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := g.Interface.Get(namespace, k8s.Stone, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (g *Stone) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := g.Interface.List(namespace, k8s.Stone, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (g *Stone) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	item, isUpdate, err := g.Interface.Apply(namespace, k8s.Stone, name, unstructuredExtend)
	if err != nil {
		return nil, isUpdate, err
	}
	return item, isUpdate, nil
}

func (g *Stone) Delete(namespace, name string) error {
	err := g.Interface.Delete(namespace, k8s.Stone, name)
	if err != nil {
		return err
	}
	return nil
}
