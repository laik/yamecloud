package workloads

import (
	"fmt"
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

var _ service.IResourceService = &Template{}

type Template struct {
	service.Interface
}

func NewTemplate(svcInterface service.Interface) *Template {
	srv := &Template{Interface: svcInterface}
	svcInterface.Install(k8s.Workloads, srv)
	return srv
}

func (g *Template) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := g.Interface.Get(namespace, k8s.Workloads, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (g *Template) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := g.Interface.List(namespace, k8s.Workloads, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (g *Template) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	item, isUpdate, err := g.Interface.Apply(namespace, k8s.Workloads, name, unstructuredExtend)
	if err != nil {
		return nil, isUpdate, err
	}
	return item, isUpdate, nil
}

func (g *Template) CreateStone(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	item, isUpdate, err := g.Interface.Apply(namespace, k8s.Stone, name, unstructuredExtend)
	if err != nil {
		return nil, isUpdate, err
	}
	return item, isUpdate, nil
}

func (g *Template) CreateDeployment(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	item, isUpdate, err := g.Interface.Apply(namespace, k8s.Deployment, name, unstructuredExtend)
	if err != nil {
		return nil, isUpdate, err
	}
	return item, isUpdate, nil
}

func (g *Template) LabelKV(namespace, name, k, v string) (*service.UnstructuredExtend, error) {
	data := fmt.Sprintf(`{"metadata":{"labels":{"%s":"%s"}}}`, k, v)
	item, err := g.Interface.Patch(namespace, k8s.Workloads, name, []byte(data))
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (g *Template) Delete(namespace, name string) error {
	err := g.Interface.Delete(namespace, k8s.Workloads, name)
	if err != nil {
		return err
	}
	return nil
}
