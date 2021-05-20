package workloads

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

var _ service.IResourceService = &CronJob{}

type CronJob struct {
	service.Interface
}

func NewCronJob(svcInterface service.Interface) *CronJob {
	srv := &CronJob{Interface: svcInterface}
	svcInterface.Install(k8s.CronJob, srv)
	return srv
}

func (g *CronJob) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := g.Interface.Get(namespace, k8s.CronJob, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (g *CronJob) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := g.Interface.List(namespace, k8s.CronJob, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (g *CronJob) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	item, isUpdate, err := g.Interface.Apply(namespace, k8s.CronJob, name, unstructuredExtend)
	if err != nil {
		return nil, isUpdate, err
	}
	return item, isUpdate, nil
}

func (g *CronJob) Delete(namespace, name string) error {
	err := g.Interface.Delete(namespace, k8s.CronJob, name)
	if err != nil {
		return err
	}
	return nil
}
