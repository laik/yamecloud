package configuration

import (
	"fmt"
	"github.com/yametech/yamecloud/pkg/action/api"
)

type configurationServer struct {
	name string
	*api.Server
	// action types
}

func (s *configurationServer) Name() string { return s.name }

func NewConfigurationServer(serviceName string, server *api.Server) *configurationServer {
	configurationServer := &configurationServer{
		name:   serviceName,
		Server: server,
	}
	group := configurationServer.Group(fmt.Sprintf("/%s", serviceName))

	_ = group
	return configurationServer
}
