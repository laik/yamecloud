package tekton

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

var _ service.IResourceService = &TektonWebHook{}

type TektonWebHook struct {
	service.Interface
}

func NewTektonWebHook(svcInterface service.Interface) *TektonWebHook {
	tektonWebHook := &TektonWebHook{Interface: svcInterface}
	svcInterface.Install(k8s.TektonWebHook, tektonWebHook)
	return tektonWebHook
}

func (g *TektonWebHook) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := g.Interface.Get(namespace, k8s.TektonWebHook, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (g *TektonWebHook) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := g.Interface.List(namespace, k8s.TektonWebHook, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (g *TektonWebHook) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	item, isUpdate, err := g.Interface.Apply(namespace, k8s.TektonWebHook, name, unstructuredExtend)
	if err != nil {
		return nil, isUpdate, err
	}
	return item, isUpdate, nil
}

func (g *TektonWebHook) Delete(namespace, name string) error {
	err := g.Interface.Delete(namespace, k8s.TektonWebHook, name)
	if err != nil {
		return err
	}
	return nil
}
