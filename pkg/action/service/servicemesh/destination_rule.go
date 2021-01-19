package servicemesh

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

var _ service.IResourceService = &DestinationRule{}

type DestinationRule struct {
	service.Interface
}

func NewDestinationRule(svcInterface service.Interface) *DestinationRule {
	destinationRule := &DestinationRule{Interface: svcInterface}
	svcInterface.Install(k8s.DestinationRule, destinationRule)
	return destinationRule
}

func (d *DestinationRule) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := d.Interface.Get(namespace, k8s.DestinationRule, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (d *DestinationRule) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := d.Interface.List(namespace, k8s.DestinationRule, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (d *DestinationRule) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	item, isUpdate, err := d.Interface.Apply(namespace, k8s.DestinationRule, name, unstructuredExtend)
	if err != nil {
		return nil, isUpdate, err
	}
	return item, isUpdate, nil
}

func (d *DestinationRule) Delete(namespace, name string) error {
	err := d.Interface.Delete(namespace, k8s.DestinationRule, name)
	if err != nil {
		return err
	}
	return nil
}
