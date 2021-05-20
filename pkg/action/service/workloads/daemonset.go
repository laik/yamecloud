package workloads

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

var _ service.IResourceService = &DaemonSet{}

type DaemonSet struct {
	service.Interface
}

func NewDaemonSet(svcInterface service.Interface) *DaemonSet {
	srv := &DaemonSet{Interface: svcInterface}
	svcInterface.Install(k8s.DaemonSet, srv)
	return srv
}

func (g *DaemonSet) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := g.Interface.Get(namespace, k8s.DaemonSet, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (g *DaemonSet) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := g.Interface.List(namespace, k8s.DaemonSet, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (g *DaemonSet) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	item, isUpdate, err := g.Interface.Apply(namespace, k8s.DaemonSet, name, unstructuredExtend)
	if err != nil {
		return nil, isUpdate, err
	}
	return item, isUpdate, nil
}

func (g *DaemonSet) Delete(namespace, name string) error {
	err := g.Interface.Delete(namespace, k8s.DaemonSet, name)
	if err != nil {
		return err
	}
	return nil
}
