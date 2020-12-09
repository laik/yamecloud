package workload

import (
	"fmt"
	"github.com/yametech/yamecloud/pkg/action/api"
)

type ServerImpl struct {
	name string
	*api.Server
}

func NewWorkloadServer(serviceName string, server *api.Server) *ServerImpl {
	serverImpl := &ServerImpl{
		name:   serviceName,
		Server: server,
	}
	group := serverImpl.Group(fmt.Sprintf("/%s", serviceName))

	// clusterRole
	{
		group.GET("/apis/rbac.authorization.k8s.io/v1/clusterroles", serverImpl.ListClusterRole)
		group.GET("/apis/rbac.authorization.k8s.io/v1/namespaces/:namespace/clusterroles/:name", serverImpl.GetClusterRole)
		group.POST("/apis/rbac.authorization.k8s.io/v1/clusterroles", serverImpl.ApplyClusterRole)
	}

	return serverImpl
}

func (s *ServerImpl) Name() string {
	return s.name
}
