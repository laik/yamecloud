package cluster

import (
	"fmt"
	"github.com/yametech/yamecloud/pkg/action/api"
)

type clusterServer struct {
	name string
	*api.Server
	// action types
}

func (s *clusterServer) Name() string { return s.name }

func NewClusterServer(serviceName string, server *api.Server) *clusterServer {
	clusterServer := &clusterServer{
		name:   serviceName,
		Server: server,
	}
	group := clusterServer.Group(fmt.Sprintf("/%s", serviceName))

	_ = group
	return clusterServer
}
