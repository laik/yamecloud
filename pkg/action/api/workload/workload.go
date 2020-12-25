package workload

import (
	"fmt"
	"github.com/yametech/yamecloud/pkg/action/api"
)

type workloadServer struct {
	name string
	*api.Server
}

func (s *workloadServer) Name() string {
	return s.name
}

func NewWorkloadServer(serviceName string, server *api.Server) *workloadServer {
	workloadServer := &workloadServer{
		name:   serviceName,
		Server: server,
	}
	group := workloadServer.Group(fmt.Sprintf("/%s", serviceName))

	_ = group
	return workloadServer
}
