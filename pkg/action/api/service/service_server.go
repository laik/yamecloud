package service

import (
	"fmt"
	"github.com/yametech/yamecloud/pkg/action/api"
	"github.com/yametech/yamecloud/pkg/action/service/service"
)

type serviceServer struct {
	name string
	*api.Server
	// action services
	*service.Endpoint
	*service.Ingress
	*service.Service
}

func (s *serviceServer) Name() string {
	return s.name
}

func NewServiceServer(serviceName string, server *api.Server) *serviceServer {
	serviceServer := &serviceServer{
		name:   serviceName,
		Server: server,
		// action service
		Endpoint: service.NewEndpoint(server.Interface),
		Ingress:  service.NewIngress(server.Interface),
		Service:  service.NewService(server.Interface),
	}
	group := serviceServer.Group(fmt.Sprintf("/%s", serviceName))
	// service
	// service
	{
		group.GET("/api/v1/services", serviceServer.ListService)
		group.GET("/api/v1/namespaces/:namespace/services", serviceServer.ListService)
		group.GET("/api/v1/namespaces/:namespace/services/:name", serviceServer.GetService)
		group.POST("/api/v1/namespaces/:namespace/services", serviceServer.ApplyService)
		group.DELETE("/api/v1/namespaces/:namespace/services/:name", serviceServer.DeleteService)
		group.PUT("/api/v1/namespaces/:namespace/services/:name", serviceServer.UpdateService)
	}

	// endpoint
	{
		group.GET("/api/v1/endpoints", serviceServer.ListEndpoint)
		group.GET("/api/v1/namespaces/:namespace/endpoints", serviceServer.ListEndpoint)
		group.GET("/api/v1/namespaces/:namespace/endpoints/:name", serviceServer.GetEndpoint)
		group.POST("/api/v1/namespaces/:namespace/endpoints", serviceServer.ApplyEndpoint)
		group.DELETE("/api/v1/namespaces/:namespace/endpoints/:name", serviceServer.DeleteEndpoint)
	}

	// ingress
	{
		group.GET("/apis/extensions/v1beta1/ingresses", serviceServer.ListIngress)
		group.GET("/apis/extensions/v1beta1/namespaces/:namespace/ingresses", serviceServer.ListIngress)
		group.GET("/apis/extensions/v1beta1/namespaces/:namespace/ingresses/:name", serviceServer.GetIngress)
		group.POST("/apis/extensions/v1beta1/namespaces/:namespace/ingresses", serviceServer.ApplyIngress)
		group.DELETE("/apis/extensions/v1beta1/namespaces/:namespace/ingresses/:name", serviceServer.DeleteIngress)
		group.PUT("/apis/extensions/v1beta1/namespaces/:namespace/ingresses/:name", serviceServer.UpdateIngress)
	}

	return serviceServer
}
