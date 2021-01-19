package tenant

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

var _ service.IResourceService = &BaseDepartment{}

type BaseDepartment struct {
	service.Interface
}

func NewBaseDepartment(svcInterface service.Interface) *BaseDepartment {
	baseDepartment := &BaseDepartment{Interface: svcInterface}
	svcInterface.Install(k8s.BaseDepartment, baseDepartment)
	return baseDepartment
}

func (b *BaseDepartment) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := b.Interface.Get(namespace, k8s.BaseDepartment, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (b *BaseDepartment) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := b.Interface.List(namespace, k8s.BaseDepartment, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (b *BaseDepartment) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	item, isUpdate, err := b.Interface.Apply(namespace, k8s.BaseDepartment, name, unstructuredExtend)
	if err != nil {
		return nil, isUpdate, err
	}
	return item, isUpdate, nil
}

func (b *BaseDepartment) Delete(namespace, name string) error {
	err := b.Interface.Delete(namespace, k8s.BaseDepartment, name)
	if err != nil {
		return err
	}
	return nil
}
