package tekton

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

var _ service.IResourceService = &Task{}

type Task struct {
	service.Interface
}

func NewTask(svcInterface service.Interface) *Task {
	task := &Task{Interface: svcInterface}
	svcInterface.Install(k8s.Task, task)
	return task
}

func (g *Task) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := g.Interface.Get(namespace, k8s.Task, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (g *Task) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := g.Interface.List(namespace, k8s.Task, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (g *Task) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	item, isUpdate, err := g.Interface.Apply(namespace, k8s.Task, name, unstructuredExtend)
	if err != nil {
		return nil, isUpdate, err
	}
	return item, isUpdate, nil
}

func (g *Task) Delete(namespace, name string) error {
	err := g.Interface.Delete(namespace, k8s.Task, name)
	if err != nil {
		return err
	}
	return nil
}
