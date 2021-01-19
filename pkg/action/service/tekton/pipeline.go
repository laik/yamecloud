package tekton

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

var _ service.IResourceService = &Pipeline{}

type Pipeline struct {
	service.Interface
}

func NewPipeline(svcInterface service.Interface) *Pipeline {
	pipeline := &Pipeline{Interface: svcInterface}
	svcInterface.Install(k8s.Pipeline, pipeline)
	return pipeline
}

func (g *Pipeline) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := g.Interface.Get(namespace, k8s.Pipeline, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (g *Pipeline) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := g.Interface.List(namespace, k8s.Pipeline, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (g *Pipeline) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	item, isUpdate, err := g.Interface.Apply(namespace, k8s.Pipeline, name, unstructuredExtend)
	if err != nil {
		return nil, isUpdate, err
	}
	return item, isUpdate, nil
}

func (g *Pipeline) Delete(namespace, name string) error {
	err := g.Interface.Delete(namespace, k8s.Pipeline, name)
	if err != nil {
		return err
	}
	return nil
}
