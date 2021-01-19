package tenant

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

var _ service.IResourceService = &BaseRoleUser{}

type BaseRoleUser struct {
	service.Interface
}

func NewBaseRoleUser(svcInterface service.Interface) *BaseRoleUser {
	baseRoleUser := &BaseRoleUser{Interface: svcInterface}
	svcInterface.Install(k8s.BaseRoleUser, baseRoleUser)
	return baseRoleUser
}

func (b *BaseRoleUser) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := b.Interface.Get(namespace, k8s.BaseRoleUser, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (b *BaseRoleUser) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := b.Interface.List(namespace, k8s.BaseRoleUser, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (b *BaseRoleUser) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	item, isUpdate, err := b.Interface.Apply(namespace, k8s.BaseRoleUser, name, unstructuredExtend)
	if err != nil {
		return nil, isUpdate, err
	}
	return item, isUpdate, nil
}

func (b *BaseRoleUser) Delete(namespace, name string) error {
	err := b.Interface.Delete(namespace, k8s.BaseRoleUser, name)
	if err != nil {
		return err
	}
	return nil
}
