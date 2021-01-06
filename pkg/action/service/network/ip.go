package network

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

var _ service.IResourceService = &IP{}

type IP struct {
	service.Interface
}

func NewIP(svcInterface service.Interface) *IP {
	ip := &IP{Interface: svcInterface}
	svcInterface.Install(k8s.IP, ip)
	return ip
}

func (b *IP) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := b.Interface.Get(namespace, k8s.IP, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (b *IP) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := b.Interface.List(namespace, k8s.IP, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (b *IP) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, error) {
	item, err := b.Interface.Apply(namespace, k8s.IP, name, unstructuredExtend)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (b *IP) Delete(namespace, name string) error {
	err := b.Interface.Delete(namespace, k8s.IP, name)
	if err != nil {
		return err
	}
	return nil
}
