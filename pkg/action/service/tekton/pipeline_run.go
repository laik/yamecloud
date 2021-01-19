package tekton

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

var _ service.IResourceService = &PipelineRun{}

type PipelineRun struct {
	service.Interface
}

func NewPipelineRun(svcInterface service.Interface) *PipelineRun {
	pipelineRun := &PipelineRun{Interface: svcInterface}
	svcInterface.Install(k8s.PipelineRun, pipelineRun)
	return pipelineRun
}

func (g *PipelineRun) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := g.Interface.Get(namespace, k8s.PipelineRun, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (g *PipelineRun) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := g.Interface.List(namespace, k8s.PipelineRun, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (g *PipelineRun) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	item, isUpdate, err := g.Interface.Apply(namespace, k8s.PipelineRun, name, unstructuredExtend)
	if err != nil {
		return nil, isUpdate, err
	}
	return item, isUpdate, nil
}

func (g *PipelineRun) Delete(namespace, name string) error {
	err := g.Interface.Delete(namespace, k8s.PipelineRun, name)
	if err != nil {
		return err
	}
	return nil
}
