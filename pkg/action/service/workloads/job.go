package workloads

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

var _ service.IResourceService = &Job{}

type Job struct {
	service.Interface
}

func NewJob(svcInterface service.Interface) *Job {
	job := &Job{Interface: svcInterface}
	svcInterface.Install(k8s.Job, job)
	return job
}

func (g *Job) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := g.Interface.Get(namespace, k8s.Job, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (g *Job) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := g.Interface.List(namespace, k8s.Job, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (g *Job) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	item, isUpdate, err := g.Interface.Apply(namespace, k8s.Job, name, unstructuredExtend)
	if err != nil {
		return nil, isUpdate, err
	}
	return item, isUpdate, nil
}

func (g *Job) Delete(namespace, name string) error {
	err := g.Interface.Delete(namespace, k8s.Job, name)
	if err != nil {
		return err
	}
	return nil
}
