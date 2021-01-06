package servicemesh

import (
	"fmt"
	"github.com/yametech/yamecloud/pkg/action/api"
	"github.com/yametech/yamecloud/pkg/action/service/servicemesh"
)

type serviceMeshServer struct {
	name string
	*api.Server
	// action services
	*servicemesh.Gateway
	*servicemesh.ServiceEntry
	*servicemesh.DestinationRule
	*servicemesh.Sidecar
	*servicemesh.VirtualService
	*servicemesh.WorkloadEntry
}

func (s *serviceMeshServer) Name() string {
	return s.name
}

func NewServiceMeshServer(serviceName string, server *api.Server) *serviceMeshServer {
	serviceMeshServer := &serviceMeshServer{
		name:   serviceName,
		Server: server,
		// action service
		Gateway:         servicemesh.NewGateway(server.Interface),
		DestinationRule: servicemesh.NewDestinationRule(server.Interface),
		ServiceEntry:    servicemesh.NewServiceEntry(server.Interface),
		Sidecar:         servicemesh.NewSidecar(server.Interface),
		VirtualService:  servicemesh.NewVirtualService(server.Interface),
		WorkloadEntry:   servicemesh.NewWorkloadEntry(server.Interface),
	}
	group := serviceMeshServer.Group(fmt.Sprintf("/%s", serviceName))

	// service mesh
	// gateway
	{
		group.GET("/apis/networking.istio.io/v1beta1/gateways", serviceMeshServer.ListGateway)
		group.GET("/apis/networking.istio.io/v1beta1/namespaces/:namespace/gateways", serviceMeshServer.ListGateway)
		group.GET("/apis/networking.istio.io/v1beta1/namespaces/:namespace/gateways/:name", serviceMeshServer.GetGateway)
		group.POST("/apis/networking.istio.io/v1beta1/gateways", serviceMeshServer.ApplyGateway)
	}
	// service entry
	{
		group.GET("/apis/networking.istio.io/v1beta1/serviceentrys", serviceMeshServer.ListServiceEntry)
		group.GET("/apis/networking.istio.io/v1beta1/namespaces/:namespace/serviceentrys", serviceMeshServer.ListServiceEntry)
		group.GET("/apis/networking.istio.io/v1beta1/namespaces/:namespace/serviceentrys/:name", serviceMeshServer.GetServiceEntry)
		group.POST("/apis/networking.istio.io/v1beta1/serviceentrys", serviceMeshServer.ApplyServiceEntry)
	}
	// virtual service
	{
		group.GET("/apis/networking.istio.io/v1beta1/virtualservices", serviceMeshServer.ListVirtualService)
		group.GET("/apis/networking.istio.io/v1beta1/namespaces/:namespace/virtualservices", serviceMeshServer.ListVirtualService)
		group.GET("/apis/networking.istio.io/v1beta1/namespaces/:namespace/virtualservices/:name", serviceMeshServer.GetVirtualService)
		group.POST("/apis/networking.istio.io/v1beta1/virtualservices", serviceMeshServer.ApplyVirtualService)
	}
	// destination rule
	{
		group.GET("/apis/networking.istio.io/v1beta1/destinationrules", serviceMeshServer.ListDestinationRule)
		group.GET("/apis/networking.istio.io/v1beta1/namespaces/:namespace/destinationrules", serviceMeshServer.ListDestinationRule)
		group.GET("/apis/networking.istio.io/v1beta1/namespaces/:namespace/destinationrules/:name", serviceMeshServer.GetDestinationRule)
		group.POST("/apis/networking.istio.io/v1beta1/destinationrules", serviceMeshServer.ApplyDestinationRule)
	}

	// workload entry
	{
		group.GET("/apis/networking.istio.io/v1beta1/workloadentries", serviceMeshServer.ListWorkloadEntry)
		group.GET("/apis/networking.istio.io/v1beta1/namespaces/:namespace/workloadentries", serviceMeshServer.ListWorkloadEntry)
		group.GET("/apis/networking.istio.io/v1beta1/namespaces/:namespace/workloadentries/:name", serviceMeshServer.GetWorkloadEntry)
		group.POST("/apis/networking.istio.io/v1beta1/workloadentries", serviceMeshServer.ApplyWorkloadEntry)
	}

	// sidecar
	{
		group.GET("/apis/networking.istio.io/v1beta1/sidecars", serviceMeshServer.ListSidecar)
		group.GET("/apis/networking.istio.io/v1beta1/namespaces/:namespace/sidecars", serviceMeshServer.ListSidecar)
		group.GET("/apis/networking.istio.io/v1beta1/namespaces/:namespace/sidecars/:name", serviceMeshServer.GetSidecar)
		group.POST("/apis/networking.istio.io/v1beta1/sidecars", serviceMeshServer.ApplySidecar)
	}

	return serviceMeshServer
}
