package tekton

import (
	"fmt"
	"github.com/yametech/yamecloud/pkg/action/api"
)

type serviceMeshServer struct {
	name string
	*api.Server
	// action services
}

func (s *serviceMeshServer) Name() string {
	return s.name
}

func NewServiceMeshServer(serviceName string, server *api.Server) *serviceMeshServer {
	serviceMeshServer := &serviceMeshServer{
		name:   serviceName,
		Server: server,
		// action service
	}
	group := serviceMeshServer.Group(fmt.Sprintf("/%s", serviceName))

	_ = group
	return serviceMeshServer
}
