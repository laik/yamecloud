package network

import (
	"fmt"
	"github.com/yametech/yamecloud/pkg/action/api"
)

type networkServer struct {
	name string
	*api.Server
	// action services
}

func (s *networkServer) Name() string { return s.name }

func NewNetworkServer(serviceName string, server *api.Server) *networkServer {
	networkServer := &networkServer{
		name:   serviceName,
		Server: server,
	}
	group := networkServer.Group(fmt.Sprintf("/%s", serviceName))

	_ = group
	return networkServer
}
