package workloads

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

var _ service.IResourceService = &Secret{}

type Secret struct {
	service.Interface
}

func NewSecret(svcInterface service.Interface) *Secret {
	srv := &Secret{Interface: svcInterface}
	svcInterface.Install(k8s.Secret, srv)
	return srv
}

func (g *Secret) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := g.Interface.Get(namespace, k8s.Secret, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (g *Secret) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := g.Interface.List(namespace, k8s.Secret, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (g *Secret) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	item, isUpdate, err := g.Interface.Apply(namespace, k8s.Secret, name, unstructuredExtend)
	if err != nil {
		return nil, isUpdate, err
	}
	return item, isUpdate, nil
}

func (g *Secret) Delete(namespace, name string) error {
	err := g.Interface.Delete(namespace, k8s.Secret, name)
	if err != nil {
		return err
	}
	return nil
}
