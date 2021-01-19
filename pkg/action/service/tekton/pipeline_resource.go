package tekton

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

var _ service.IResourceService = &PipelineResource{}

type PipelineResource struct {
	service.Interface
}

func NewPipelineResource(svcInterface service.Interface) *PipelineResource {
	pipelineResource := &PipelineResource{Interface: svcInterface}
	svcInterface.Install(k8s.PipelineResource, pipelineResource)
	return pipelineResource
}

func (g *PipelineResource) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := g.Interface.Get(namespace, k8s.PipelineResource, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (g *PipelineResource) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := g.Interface.List(namespace, k8s.PipelineResource, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (g *PipelineResource) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	item, isUpdate, err := g.Interface.Apply(namespace, k8s.PipelineResource, name, unstructuredExtend)
	if err != nil {
		return nil, isUpdate, err
	}
	return item, isUpdate, nil
}

func (g *PipelineResource) Delete(namespace, name string) error {
	err := g.Interface.Delete(namespace, k8s.PipelineResource, name)
	if err != nil {
		return err
	}
	return nil
}
