package workloads

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

var _ service.IResourceService = &PV{}

type PV struct {
	service.Interface
}

func NewPV(svcInterface service.Interface) *PV {
	srv := &PV{Interface: svcInterface}
	svcInterface.Install(k8s.PersistentVolume, srv)
	return srv
}

func (g *PV) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := g.Interface.Get(namespace, k8s.PersistentVolume, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (g *PV) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := g.Interface.List(namespace, k8s.PersistentVolume, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (g *PV) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	item, isUpdate, err := g.Interface.Apply(namespace, k8s.PersistentVolume, name, unstructuredExtend)
	if err != nil {
		return nil, isUpdate, err
	}
	return item, isUpdate, nil
}

func (g *PV) Delete(namespace, name string) error {
	err := g.Interface.Delete(namespace, k8s.PersistentVolume, name)
	if err != nil {
		return err
	}
	return nil
}
