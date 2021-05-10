package shell

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

var _ service.IResourceService = &Pod{}

type Pod struct {
	service.Interface
}

func NewPod(svcInterface service.Interface) *Pod {
	endpoint := &Pod{Interface: svcInterface}
	svcInterface.Install(k8s.Pod, endpoint)
	return endpoint
}

func (p Pod) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := p.Interface.Get(namespace, k8s.Pod, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (p Pod) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	panic("implement me")
}

func (p Pod) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	panic("implement me")
}
