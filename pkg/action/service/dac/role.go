package dac

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

var _ service.IResourceService = &Role{}

type Role struct {
	service.Interface
}

func NewRole(svcInterface service.Interface) *Role {
	role := &Role{Interface: svcInterface}
	svcInterface.Install(k8s.Role, role)
	return role
}

func (c *Role) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := c.Interface.Get(namespace, k8s.Role, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (c *Role) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := c.Interface.List(namespace, k8s.Role, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (c *Role) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	item, isUpdate, err := c.Interface.Apply(namespace, k8s.Role, name, unstructuredExtend)
	if err != nil {
		return nil, isUpdate, err
	}
	return item, isUpdate, nil
}

