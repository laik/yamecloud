package tenant

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

var _ service.IResourceService = &BaseTenant{}

type BaseTenant struct {
	service.Interface
}

func NewBaseTenant(svcInterface service.Interface) *BaseTenant {
	baseTenant := &BaseTenant{Interface: svcInterface}
	svcInterface.Install(k8s.BaseDepartment, baseTenant)
	return baseTenant
}

func (b *BaseTenant) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := b.Interface.Get(namespace, k8s.BaseDepartment, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (b *BaseTenant) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := b.Interface.List(namespace, k8s.BaseDepartment, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (b *BaseTenant) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	item, isUpdate, err := b.Interface.Apply(namespace, k8s.BaseDepartment, name, unstructuredExtend)
	if err != nil {
		return nil, isUpdate, err
	}
	return item, isUpdate, nil
}

func (b *BaseTenant) Delete(namespace, name string) error {
	err := b.Interface.Delete(namespace, k8s.BaseDepartment, name)
	if err != nil {
		return err
	}
	return nil
}
