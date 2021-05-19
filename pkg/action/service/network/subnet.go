package network

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

var _ service.IResourceService = &SubNet{}

type SubNet struct {
	service.Interface
}

func NewSubnet(svcInterface service.Interface) *SubNet {
	subnet := &SubNet{Interface: svcInterface}
	svcInterface.Install(k8s.SubNet, subnet)
	return subnet
}

func (b *SubNet) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := b.Interface.Get(namespace, k8s.SubNet, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (b *SubNet) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := b.Interface.List(namespace, k8s.SubNet, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (b *SubNet) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	item, isUpdate, err := b.Interface.Apply(namespace, k8s.SubNet, name, unstructuredExtend)
	if err != nil {
		return nil, isUpdate, err
	}
	return item, isUpdate, nil
}

func (b *SubNet) Delete(namespace, name string) error {
	err := b.Interface.Delete(namespace, k8s.SubNet, name)
	if err != nil {
		return err
	}
	return nil
}
