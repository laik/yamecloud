package workloadplus

import (
	"fmt"
	"github.com/yametech/yamecloud/pkg/action/api"
)

type workloadPlusServer struct {
	name string
	*api.Server
}

func (s *workloadPlusServer) Name() string {
	return s.name
}

func NewWorkloadPlusServer(serviceName string, server *api.Server) *workloadPlusServer {
	workloadPlusServer := &workloadPlusServer{
		name:   serviceName,
		Server: server,
	}
	group := workloadPlusServer.Group(fmt.Sprintf("/%s", serviceName))

	_ = group
	return workloadPlusServer
}
