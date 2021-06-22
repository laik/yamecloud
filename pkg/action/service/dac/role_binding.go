package dac

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

var _ service.IResourceService = &RoleBinding{}

type RoleBinding struct {
	service.Interface
}

func NewRoleBinding(svcInterface service.Interface) *RoleBinding {
	roleBinding := &RoleBinding{Interface: svcInterface}
	svcInterface.Install(k8s.RoleBinding, roleBinding)
	return roleBinding
}

func (c *RoleBinding) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := c.Interface.Get(namespace, k8s.RoleBinding, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (c *RoleBinding) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := c.Interface.List(namespace, k8s.RoleBinding, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (c *RoleBinding) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	item, isUpdate, err := c.Interface.Apply(namespace, k8s.RoleBinding, name, unstructuredExtend)
	if err != nil {
		return nil, isUpdate, err
	}
	return item, isUpdate, nil
}

