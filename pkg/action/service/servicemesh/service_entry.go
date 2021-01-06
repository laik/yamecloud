package servicemesh

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

var _ service.IResourceService = &ServiceEntry{}

type ServiceEntry struct {
	service.Interface
}

func NewServiceEntry(svcInterface service.Interface) *ServiceEntry {
	serviceEntry := &ServiceEntry{Interface: svcInterface}
	svcInterface.Install(k8s.ServiceEntry, serviceEntry)
	return serviceEntry
}

func (s *ServiceEntry) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := s.Interface.Get(namespace, k8s.ServiceEntry, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (s *ServiceEntry) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := s.Interface.List(namespace, k8s.ServiceEntry, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (s *ServiceEntry) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, error) {
	item, err := s.Interface.Apply(namespace, k8s.ServiceEntry, name, unstructuredExtend)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (s *ServiceEntry) Delete(namespace, name string) error {
	err := s.Interface.Delete(namespace, k8s.ServiceEntry, name)
	if err != nil {
		return err
	}
	return nil
}
