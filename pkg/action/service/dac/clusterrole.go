package dac

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

type ClusterRole struct {
	service.Interface
}

func NewClusterRole(svcInterface service.Interface) *ClusterRole {
	clusterRole := &ClusterRole{svcInterface}
	svcInterface.Install(k8s.ClusterRole, clusterRole)
	return clusterRole
}

func (c *ClusterRole) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := c.Interface.Get(namespace, k8s.ClusterRole, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (c *ClusterRole) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := c.Interface.List(namespace, k8s.ClusterRole, selector)
	if err != nil {
		return nil, err
	}
	return list, nil
}
