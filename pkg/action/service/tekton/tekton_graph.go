package tekton

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

var _ service.IResourceService = &TektonGraph{}

type TektonGraph struct {
	service.Interface
}

func NewTektonGraph(svcInterface service.Interface) *TektonGraph {
	tektonGraph := &TektonGraph{Interface: svcInterface}
	svcInterface.Install(k8s.TektonGraph, tektonGraph)
	return tektonGraph
}

func (g *TektonGraph) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := g.Interface.Get(namespace, k8s.TektonGraph, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (g *TektonGraph) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := g.Interface.List(namespace, k8s.TektonGraph, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (g *TektonGraph) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	item, isUpdate, err := g.Interface.Apply(namespace, k8s.TektonGraph, name, unstructuredExtend)
	if err != nil {
		return nil, isUpdate, err
	}
	return item, isUpdate, nil
}

func (g *TektonGraph) Delete(namespace, name string) error {
	err := g.Interface.Delete(namespace, k8s.TektonGraph, name)
	if err != nil {
		return err
	}
	return nil
}
