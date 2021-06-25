package dac

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

var _ service.IResourceService = &PSP{}

type PSP struct {
	service.Interface
}

func NewPSP(svcInterface service.Interface) *PSP {
	role := &PSP{Interface: svcInterface}
	svcInterface.Install(k8s.PodSecurityPolicie, role)
	return role
}

func (c *PSP) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := c.Interface.Get(namespace, k8s.PodSecurityPolicie, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (c *PSP) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := c.Interface.List(namespace, k8s.PodSecurityPolicie, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (c *PSP) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	item, isUpdate, err := c.Interface.Apply(namespace, k8s.PodSecurityPolicie, name, unstructuredExtend)
	if err != nil {
		return nil, isUpdate, err
	}
	return item, isUpdate, nil
}
