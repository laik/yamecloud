package workload

import (
	"github.com/yametech/yamecloud/pkg/action/api"
	"github.com/yametech/yamecloud/pkg/action/service/dac"
)

type Workload struct {
	*api.Server
	clusterRole *dac.ClusterRole
}

func NewWorkloadServer(server *api.Server, clusterRole *dac.ClusterRole) *Workload {
	return &Workload{
		Server:      server,
		clusterRole: clusterRole,
	}
}
