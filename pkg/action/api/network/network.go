package network

import (
	"fmt"
	"github.com/yametech/yamecloud/pkg/action/api"
	"github.com/yametech/yamecloud/pkg/action/service/network"
)

type networkServer struct {
	name string
	*api.Server
	// action services
	*network.IP
	*network.NetworkAttachmentDefinition
}

func (s *networkServer) Name() string { return s.name }

func NewNetworkServer(serviceName string, server *api.Server) *networkServer {
	networkServer := &networkServer{
		name:   serviceName,
		Server: server,
		IP:     network.NewIP(server.Interface),
	}
	group := networkServer.Group(fmt.Sprintf("/%s", serviceName))

	// this route allow admin or network operations role maintain
	group.GET("/apis/kubeovn.io/v1/ips", networkServer.ListIPs)
	_ = group
	return networkServer
}
