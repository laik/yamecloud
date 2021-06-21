package dac

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

var _ service.IResourceService = &ServiceAccount{}

type ServiceAccount struct {
	service.Interface
}

func NewServiceAccount(svcInterface service.Interface) *ServiceAccount {
	serviceAccount := &ServiceAccount{Interface: svcInterface}
	svcInterface.Install(k8s.ServiceAccount, serviceAccount)
	return serviceAccount
}

func (c *ServiceAccount) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := c.Interface.Get(namespace, k8s.ServiceAccount, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (c *ServiceAccount) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := c.Interface.List(namespace, k8s.ServiceAccount, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (c *ServiceAccount) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	item, isUpdate, err := c.Interface.Apply(namespace, k8s.ServiceAccount, name, unstructuredExtend)
	if err != nil {
		return nil, isUpdate, err
	}
	return item, isUpdate, nil
}

