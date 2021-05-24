package workloads

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

var _ service.IResourceService = &Node{}

type Node struct {
	service.Interface
}

func NewNode(svcInterface service.Interface) *Node {
	srv := &Node{Interface: svcInterface}
	svcInterface.Install(k8s.Node, srv)
	return srv
}

func (g *Node) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := g.Interface.Get(namespace, k8s.Node, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (g *Node) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := g.Interface.List(namespace, k8s.Node, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (g *Node) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	item, isUpdate, err := g.Interface.Apply(namespace, k8s.Node, name, unstructuredExtend)
	if err != nil {
		return nil, isUpdate, err
	}
	return item, isUpdate, nil
}

func (g *Node) Delete(namespace, name string) error {
	err := g.Interface.Delete(namespace, k8s.Node, name)
	if err != nil {
		return err
	}
	return nil
}
