package workloads

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

var _ service.IResourceService = &Deployment{}

type Deployment struct {
	service.Interface
}

func NewDeployment(svcInterface service.Interface) *Deployment {
	srv := &Deployment{Interface: svcInterface}
	svcInterface.Install(k8s.Deployment, srv)
	return srv
}

func (g *Deployment) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := g.Interface.Get(namespace, k8s.Deployment, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (g *Deployment) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := g.Interface.List(namespace, k8s.Deployment, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (g *Deployment) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	item, isUpdate, err := g.Interface.Apply(namespace, k8s.Deployment, name, unstructuredExtend)
	if err != nil {
		return nil, isUpdate, err
	}
	return item, isUpdate, nil
}

func (g *Deployment) Delete(namespace, name string) error {
	err := g.Interface.Delete(namespace, k8s.Deployment, name)
	if err != nil {
		return err
	}
	return nil
}
