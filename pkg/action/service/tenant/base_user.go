package tenant

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

var _ service.IResourceService = &BaseUser{}

type BaseUser struct {
	service.Interface
}

func NewBaseUser(svcInterface service.Interface) *BaseUser {
	baseUser := &BaseUser{Interface: svcInterface}
	svcInterface.Install(k8s.BaseUser, baseUser)
	return baseUser
}

func (b *BaseUser) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := b.Interface.Get(namespace, k8s.BaseUser, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (b *BaseUser) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := b.Interface.List(namespace, k8s.BaseUser, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (b *BaseUser) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	item, isUpdate, err := b.Interface.Apply(namespace, k8s.BaseUser, name, unstructuredExtend)
	if err != nil {
		return nil, isUpdate, err
	}
	return item, isUpdate, nil
}

func (b *BaseUser) Delete(namespace, name string) error {
	err := b.Interface.Delete(namespace, k8s.BaseUser, name)
	if err != nil {
		return err
	}
	return nil
}
