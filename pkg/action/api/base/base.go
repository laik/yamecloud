package base

import (
	"fmt"
	"github.com/yametech/yamecloud/pkg/action/api"
	"github.com/yametech/yamecloud/pkg/action/service/tenant"
)

type baseServer struct {
	name string
	*api.Server
	// action services
	*tenant.BaseDepartment
}

func (s *baseServer) Name() string { return s.name }

func NewBaseServer(serviceName string, server *api.Server) *baseServer {
	baseServer := &baseServer{
		name:   serviceName,
		Server: server,
	}
	group := baseServer.Group(fmt.Sprintf("/%s", serviceName))

	_ = group
	return baseServer
}
