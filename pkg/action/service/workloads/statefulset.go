package workloads

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

var _ service.IResourceService = &StatefulSet{}

type StatefulSet struct {
	service.Interface
}

func NewStatefulSet(svcInterface service.Interface) *StatefulSet {
	srv := &StatefulSet{Interface: svcInterface}
	svcInterface.Install(k8s.StatefulSet, srv)
	return srv
}

func (g *StatefulSet) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := g.Interface.Get(namespace, k8s.StatefulSet, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (g *StatefulSet) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := g.Interface.List(namespace, k8s.StatefulSet, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (g *StatefulSet) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	item, isUpdate, err := g.Interface.Apply(namespace, k8s.StatefulSet, name, unstructuredExtend)
	if err != nil {
		return nil, isUpdate, err
	}
	return item, isUpdate, nil
}

func (g *StatefulSet) Delete(namespace, name string) error {
	err := g.Interface.Delete(namespace, k8s.StatefulSet, name)
	if err != nil {
		return err
	}
	return nil
}
