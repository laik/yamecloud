package workloads

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

var _ service.IResourceService = &ConfigMap{}

type ConfigMap struct {
	service.Interface
}

func NewConfigMap(svcInterface service.Interface) *ConfigMap {
	srv := &ConfigMap{Interface: svcInterface}
	svcInterface.Install(k8s.ConfigMap, srv)
	return srv
}

func (g *ConfigMap) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := g.Interface.Get(namespace, k8s.ConfigMap, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (g *ConfigMap) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := g.Interface.List(namespace, k8s.ConfigMap, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (g *ConfigMap) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	item, isUpdate, err := g.Interface.Apply(namespace, k8s.ConfigMap, name, unstructuredExtend)
	if err != nil {
		return nil, isUpdate, err
	}
	return item, isUpdate, nil
}

func (g *ConfigMap) Delete(namespace, name string) error {
	err := g.Interface.Delete(namespace, k8s.ConfigMap, name)
	if err != nil {
		return err
	}
	return nil
}
