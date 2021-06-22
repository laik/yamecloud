package dac

import (
"github.com/yametech/yamecloud/pkg/action/service"
"github.com/yametech/yamecloud/pkg/k8s"
)

var _ service.IResourceService = &ClusterRoleBinding{}

type ClusterRoleBinding struct {
	service.Interface
}

func NewClusterRoleBinding(svcInterface service.Interface) *ClusterRoleBinding {
	clusterRoleBinding := &ClusterRoleBinding{Interface: svcInterface}
	svcInterface.Install(k8s.ClusterRoleBinding, clusterRoleBinding)
	return clusterRoleBinding
}

func (c *ClusterRoleBinding) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := c.Interface.Get(namespace, k8s.ClusterRoleBinding, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (c *ClusterRoleBinding) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := c.Interface.List(namespace, k8s.ClusterRoleBinding, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (c *ClusterRoleBinding) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	item, isUpdate, err := c.Interface.Apply(namespace, k8s.ClusterRoleBinding, name, unstructuredExtend)
	if err != nil {
		return nil, isUpdate, err
	}
	return item, isUpdate, nil
}
