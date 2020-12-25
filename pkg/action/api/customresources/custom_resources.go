package customresources

import (
	"fmt"
	"github.com/yametech/yamecloud/pkg/action/api"
)

type customResourcesServer struct {
	name string
	*api.Server
	// action types
}

func (s *customResourcesServer) Name() string { return s.name }

func NewConfigurationServer(serviceName string, server *api.Server) *customResourcesServer {
	customResourcesServer := &customResourcesServer{
		name:   serviceName,
		Server: server,
	}
	group := customResourcesServer.Group(fmt.Sprintf("/%s", serviceName))

	_ = group
	return customResourcesServer
}
