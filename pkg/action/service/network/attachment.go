package network

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

var _ service.IResourceService = &NetworkAttachmentDefinition{}

type NetworkAttachmentDefinition struct {
	service.Interface
}

func NewNetworkAttachmentDefinition(svcInterface service.Interface) *NetworkAttachmentDefinition {
	attachment := &NetworkAttachmentDefinition{Interface: svcInterface}
	svcInterface.Install(k8s.NetworkAttachmentDefinition, attachment)
	return attachment
}

func (b *NetworkAttachmentDefinition) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := b.Interface.Get(namespace, k8s.NetworkAttachmentDefinition, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (b *NetworkAttachmentDefinition) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := b.Interface.List(namespace, k8s.NetworkAttachmentDefinition, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (b *NetworkAttachmentDefinition) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	item, isUpdate, err := b.Interface.Apply(namespace, k8s.NetworkAttachmentDefinition, name, unstructuredExtend)
	if err != nil {
		return nil, isUpdate, err
	}
	return item, isUpdate, nil
}

func (b *NetworkAttachmentDefinition) Delete(namespace, name string) error {
	err := b.Interface.Delete(namespace, k8s.NetworkAttachmentDefinition, name)
	if err != nil {
		return err
	}
	return nil
}
