package servicemesh

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

var _ service.IResourceService = &Sidecar{}

type Sidecar struct {
	service.Interface
}

func NewSidecar(svcInterface service.Interface) *Sidecar {
	sidecar := &Sidecar{Interface: svcInterface}
	svcInterface.Install(k8s.Sidecar, sidecar)
	return sidecar
}

func (s *Sidecar) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := s.Interface.Get(namespace, k8s.Sidecar, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Sidecar) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := s.Interface.List(namespace, k8s.Sidecar, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (s *Sidecar) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, error) {
	item, err := s.Interface.Apply(namespace, k8s.Sidecar, name, unstructuredExtend)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Sidecar) Delete(namespace, name string) error {
	err := s.Interface.Delete(namespace, k8s.Sidecar, name)
	if err != nil {
		return err
	}
	return nil
}
