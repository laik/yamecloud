package tekton

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

var _ service.IResourceService = &TaskRun{}

type TaskRun struct {
	service.Interface
}

func NewTaskRun(svcInterface service.Interface) *TaskRun {
	taskRun := &TaskRun{Interface: svcInterface}
	svcInterface.Install(k8s.TaskRun, taskRun)
	return taskRun
}

func (g *TaskRun) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := g.Interface.Get(namespace, k8s.TaskRun, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (g *TaskRun) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := g.Interface.List(namespace, k8s.TaskRun, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (g *TaskRun) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	item, isUpdate, err := g.Interface.Apply(namespace, k8s.TaskRun, name, unstructuredExtend)
	if err != nil {
		return nil, isUpdate, err
	}
	return item, isUpdate, nil
}

func (g *TaskRun) Delete(namespace, name string) error {
	err := g.Interface.Delete(namespace, k8s.TaskRun, name)
	if err != nil {
		return err
	}
	return nil
}
