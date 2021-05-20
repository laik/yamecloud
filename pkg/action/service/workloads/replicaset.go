package workloads

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

var _ service.IResourceService = &ReplicaSet{}

type ReplicaSet struct {
	service.Interface
}

func NewReplicaSet(svcInterface service.Interface) *ReplicaSet {
	srv := &ReplicaSet{Interface: svcInterface}
	svcInterface.Install(k8s.ReplicaSet, srv)
	return srv
}

func (g *ReplicaSet) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := g.Interface.Get(namespace, k8s.ReplicaSet, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (g *ReplicaSet) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := g.Interface.List(namespace, k8s.ReplicaSet, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (g *ReplicaSet) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	item, isUpdate, err := g.Interface.Apply(namespace, k8s.ReplicaSet, name, unstructuredExtend)
	if err != nil {
		return nil, isUpdate, err
	}
	return item, isUpdate, nil
}

func (g *ReplicaSet) Delete(namespace, name string) error {
	err := g.Interface.Delete(namespace, k8s.ReplicaSet, name)
	if err != nil {
		return err
	}
	return nil
}
