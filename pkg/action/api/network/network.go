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
	*network.SubNet
}

func (s *networkServer) Name() string { return s.name }

func NewNetworkServer(serviceName string, server *api.Server) *networkServer {
	networkServer := &networkServer{
		name:   serviceName,
		Server: server,
		IP:     network.NewIP(server.Interface),
		SubNet: network.NewSubnet(server.Interface),

		NetworkAttachmentDefinition: network.NewNetworkAttachmentDefinition(server.Interface),
	}
	group := networkServer.Group(fmt.Sprintf("/%s", serviceName))

	// this route allow admin or sdn operations role maintain
	// ips
	{
		group.GET("/apis/kubeovn.io/v1/ips", networkServer.ListIP)
		group.GET("/apis/kubeovn.io/v1/ips/:name", networkServer.GetIP)
		group.DELETE("/apis/kubeovn.io/v1/ips/:name", networkServer.DeleteIP)
	}

	// subnet
	{
		group.GET("/apis/kubeovn.io/v1/subnets", networkServer.ListSubNet)
		group.GET("/apis/kubeovn.io/v1/subnets/:name", networkServer.GetSubNet)
		group.POST("/apis/kubeovn.io/v1/subnets", networkServer.ApplySubNet)
		group.PUT("/apis/kubeovn.io/v1/subnets/:name", networkServer.UpdateSubNet)
		group.DELETE("/apis/kubeovn.io/v1/subnets/:name", networkServer.DeleteSubNet)
	}

	// NetworkAttachmentDefinition
	{
		group.GET("/apis/k8s.cni.cncf.io/v1/network-attachment-definitions", networkServer.ListNetWorkAttachment)
		group.GET("/apis/k8s.cni.cncf.io/v1/namespaces/:namespace/network-attachment-definitions/:name", networkServer.GetNetWorkAttachment)
		group.POST("/apis/k8s.cni.cncf.io/v1/namespaces/:namespace/network-attachment-definitions", networkServer.ApplyNetWorkAttachment)
		group.PUT("/apis/k8s.cni.cncf.io/v1/namespaces/:namespace/network-attachment-definitions/:name", networkServer.UpdateNetWorkAttachment)
		group.DELETE("/apis/k8s.cni.cncf.io/v1/namespaces/:namespace/network-attachment-definitions/:name", networkServer.DeleteNetWorkAttachment)
	}

	_ = group
	return networkServer
}
