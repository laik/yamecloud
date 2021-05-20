package workloads

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

var _ service.IResourceService = &Event{}

type Event struct {
	service.Interface
}

func NewEvent(svcInterface service.Interface) *Event {
	srv := &Event{Interface: svcInterface}
	svcInterface.Install(k8s.Event, srv)
	return srv
}

func (g *Event) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := g.Interface.Get(namespace, k8s.Event, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (g *Event) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := g.Interface.List(namespace, k8s.Event, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (g *Event) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	item, isUpdate, err := g.Interface.Apply(namespace, k8s.Event, name, unstructuredExtend)
	if err != nil {
		return nil, isUpdate, err
	}
	return item, isUpdate, nil
}

func (g *Event) Delete(namespace, name string) error {
	err := g.Interface.Delete(namespace, k8s.Event, name)
	if err != nil {
		return err
	}
	return nil
}
