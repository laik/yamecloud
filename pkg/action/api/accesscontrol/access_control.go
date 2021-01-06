package accesscontrol

import (
	"fmt"
	"github.com/yametech/yamecloud/pkg/action/api"
	"github.com/yametech/yamecloud/pkg/action/service/dac"
)

type accessControlServer struct {
	name string
	*api.Server
	// action services
	*dac.ClusterRole
}

func (s *accessControlServer) Name() string {
	return s.name
}

func NewApiServer(serviceName string, server *api.Server) *accessControlServer {
	apiServer := &accessControlServer{
		name:   serviceName,
		Server: server,
		// action service
		ClusterRole: dac.NewClusterRole(server.Interface),
	}
	group := apiServer.Group(fmt.Sprintf("/%s", serviceName))

	// access control
	// clusterRole cluster level
	{
		group.GET("/apis/rbac.authorization.k8s.io/v1/clusterroles", apiServer.ListClusterRole)
		group.GET("/apis/rbac.authorization.k8s.io/v1/clusterroles/:name", apiServer.GetClusterRole)
		group.POST("/apis/rbac.authorization.k8s.io/v1/clusterroles", apiServer.ApplyClusterRole)
	}
	// service account cluster level
	{

	}

	return apiServer
}
