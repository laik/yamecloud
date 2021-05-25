package workloads

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

var _ service.IResourceService = &PVC{}

type PVC struct {
	service.Interface
}

func NewPVC(svcInterface service.Interface) *PVC {
	srv := &PVC{Interface: svcInterface}
	svcInterface.Install(k8s.PersistentVolumeClaims, srv)
	return srv
}

func (g *PVC) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := g.Interface.Get(namespace, k8s.PersistentVolumeClaims, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (g *PVC) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := g.Interface.List(namespace, k8s.PersistentVolumeClaims, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (g *PVC) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	item, isUpdate, err := g.Interface.Apply(namespace, k8s.PersistentVolumeClaims, name, unstructuredExtend)
	if err != nil {
		return nil, isUpdate, err
	}
	return item, isUpdate, nil
}

func (g *PVC) Delete(namespace, name string) error {
	err := g.Interface.Delete(namespace, k8s.PersistentVolumeClaims, name)
	if err != nil {
		return err
	}
	return nil
}
