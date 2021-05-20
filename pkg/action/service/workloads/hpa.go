package workloads

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

var _ service.IResourceService = &HPA{}

type HPA struct {
	service.Interface
}

func NewHPA(svcInterface service.Interface) *HPA {
	srv := &HPA{Interface: svcInterface}
	svcInterface.Install(k8s.HorizontalPodAutoscaler, srv)
	return srv
}

func (g *HPA) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := g.Interface.Get(namespace, k8s.HorizontalPodAutoscaler, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (g *HPA) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := g.Interface.List(namespace, k8s.HorizontalPodAutoscaler, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (g *HPA) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	item, isUpdate, err := g.Interface.Apply(namespace, k8s.HorizontalPodAutoscaler, name, unstructuredExtend)
	if err != nil {
		return nil, isUpdate, err
	}
	return item, isUpdate, nil
}

func (g *HPA) Delete(namespace, name string) error {
	err := g.Interface.Delete(namespace, k8s.HorizontalPodAutoscaler, name)
	if err != nil {
		return err
	}
	return nil
}
