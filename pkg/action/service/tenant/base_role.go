package tenant

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

var _ service.IResourceService = &BaseRole{}

type BaseRole struct {
	service.Interface
}

func NewBaseRole(svcInterface service.Interface) *BaseRole {
	baseRole := &BaseRole{Interface: svcInterface}
	svcInterface.Install(k8s.BaseUser, baseRole)
	return baseRole
}

func (b *BaseRole) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := b.Interface.Get(namespace, k8s.BaseUser, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (b *BaseRole) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := b.Interface.List(namespace, k8s.BaseRole, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (b *BaseRole) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	item, isUpdate, err := b.Interface.Apply(namespace, k8s.BaseRole, name, unstructuredExtend)
	if err != nil {
		return nil, isUpdate, err
	}
	return item, isUpdate, nil
}

func (b *BaseRole) Delete(namespace, name string) error {
	err := b.Interface.Delete(namespace, k8s.BaseRole, name)
	if err != nil {
		return err
	}
	return nil
}
