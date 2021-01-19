package tekton

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

var _ service.IResourceService = &TektonStore{}

type TektonStore struct {
	service.Interface
}

func NewTektonStore(svcInterface service.Interface) *TektonStore {
	tektonStore := &TektonStore{Interface: svcInterface}
	svcInterface.Install(k8s.TektonStore, tektonStore)
	return tektonStore
}

func (g *TektonStore) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := g.Interface.Get(namespace, k8s.TektonStore, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (g *TektonStore) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := g.Interface.List(namespace, k8s.TektonStore, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (g *TektonStore) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	item, isUpdate, err := g.Interface.Apply(namespace, k8s.TektonStore, name, unstructuredExtend)
	if err != nil {
		return nil, isUpdate, err
	}
	return item, isUpdate, nil
}

func (g *TektonStore) Delete(namespace, name string) error {
	err := g.Interface.Delete(namespace, k8s.TektonStore, name)
	if err != nil {
		return err
	}
	return nil
}
