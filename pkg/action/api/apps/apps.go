package apps

import (
	"fmt"
	"github.com/yametech/yamecloud/pkg/action/api"
	"github.com/yametech/yamecloud/pkg/action/service/tenant"
)

type appsServer struct {
	name string
	*api.Server
	// action services
	*tenant.BaseDepartment
}

func (s *appsServer) Name() string { return s.name }

func NewAppsServer(serviceName string, server *api.Server) *appsServer {
	appsServer := &appsServer{
		name:   serviceName,
		Server: server,
	}
	group := appsServer.Group(fmt.Sprintf("/%s", serviceName))

	_ = group
	return appsServer
}
