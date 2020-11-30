package workload

import (
	"github.com/yametech/yamecloud/pkg/action/api"
	"github.com/yametech/yamecloud/pkg/action/service/dac"
)

type Workload struct {
	*api.APIServer
	clusterRole *dac.ClusterRole
}

func NewWorkloadServer(server *api.APIServer, clusterRole *dac.ClusterRole) *Workload {
	return &Workload{
		APIServer:   server,
		clusterRole: clusterRole,
	}
}
