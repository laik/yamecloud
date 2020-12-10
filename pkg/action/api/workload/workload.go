package workload

import (
	"fmt"
	"github.com/yametech/yamecloud/pkg/action/api"
	"github.com/yametech/yamecloud/pkg/action/service/dac"
)

type ServerImpl struct {
	name string
	*api.Server
	// action services
	*dac.ClusterRole
}

func NewWorkloadServer(serviceName string, server *api.Server) *ServerImpl {
	serverImpl := &ServerImpl{
		name:   serviceName,
		Server: server,
		// action service
		ClusterRole: dac.NewClusterRole(server.Interface),
	}
	group := serverImpl.Group(fmt.Sprintf("/%s", serviceName))

	// access control
	// clusterRole cluster level
	{
		group.GET("/apis/rbac.authorization.k8s.io/v1/clusterroles", serverImpl.ListClusterRole)
		group.GET("/apis/rbac.authorization.k8s.io/v1/clusterroles/:name", serverImpl.GetClusterRole)
		group.POST("/apis/rbac.authorization.k8s.io/v1/clusterroles", serverImpl.ApplyClusterRole)
	}
	// service account cluster level
	{

	}

	return serverImpl
}

func (s *ServerImpl) Name() string {
	return s.name
}
