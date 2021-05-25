package workloads

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var _ service.IResourceService = &CRD{}

type CRD struct {
	service.Interface
}

func NewCRD(svcInterface service.Interface) *CRD {
	srv := &CRD{Interface: svcInterface}
	svcInterface.Install(k8s.CustomResourceDefinition, srv)
	return srv
}

func (g *CRD) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := g.Interface.Get(namespace, k8s.CustomResourceDefinition, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (g *CRD) ListGVR(namespace string, gvr schema.GroupVersionResource, selector string) (*service.UnstructuredListExtend, error) {
	list, err := g.Interface.ListGVR(namespace, gvr, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (g *CRD) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := g.Interface.List(namespace, k8s.CustomResourceDefinition, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (g *CRD) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	item, isUpdate, err := g.Interface.Apply(namespace, k8s.CustomResourceDefinition, name, unstructuredExtend)
	if err != nil {
		return nil, isUpdate, err
	}
	return item, isUpdate, nil
}

func (g *CRD) Delete(namespace, name string) error {
	err := g.Interface.Delete(namespace, k8s.CustomResourceDefinition, name)
	if err != nil {
		return err
	}
	return nil
}
