package workloadplus

import (
	"fmt"
	"github.com/yametech/yamecloud/pkg/action/api"
	"github.com/yametech/yamecloud/pkg/action/service/workloadplus"
)

type workloadPlusServer struct {
	name string
	*api.Server

	*workloadplus.Stone
	*workloadplus.StatefulSet
	*workloadplus.Injector
}

func (s *workloadPlusServer) Name() string {
	return s.name
}

func NewWorkloadPlusServer(serviceName string, server *api.Server) *workloadPlusServer {
	workloadPlusServer := &workloadPlusServer{
		name:        serviceName,
		Server:      server,
		Stone:       workloadplus.NewStone(server),
		StatefulSet: workloadplus.NewStatefulSet(server),
		Injector:    workloadplus.NewInjector(server),
	}
	group := workloadPlusServer.Group(fmt.Sprintf("/%s", serviceName))
	//stone
	{
		group.GET("/apis/nuwa.nip.io/v1/stones", workloadPlusServer.ListStone)
		group.GET("/apis/nuwa.nip.io/v1/namespaces/:namespace/stones", workloadPlusServer.ListStone)
		group.GET("/apis/nuwa.nip.io/v1/namespaces/:namespace/stones/:name", workloadPlusServer.GetStone)
		group.POST("/apis/nuwa.nip.io/v1/stones", workloadPlusServer.ApplyStone)
		group.PUT("/apis/nuwa.nip.io/v1/namespaces/:namespace/stones/:name", workloadPlusServer.UpdateStone)
		group.DELETE("/apis/nuwa.nip.io/v1/namespaces/:namespace/stones/:name", workloadPlusServer.DeleteStone)

	}

	//statefulset
	{
		group.GET("/apis/nuwa.nip.io/v1/statefulsets", workloadPlusServer.ListStatefulSet)
		group.GET("/apis/nuwa.nip.io/v1/namespaces/:namespace/statefulsets", workloadPlusServer.ListStatefulSet)
		group.GET("/apis/nuwa.nip.io/v1/namespaces/:namespace/statefulsets/:name", workloadPlusServer.GetStatefulSet)
		group.POST("/apis/nuwa.nip.io/v1/statefulsets", workloadPlusServer.ApplyStatefulSet)
		group.PUT("/apis/nuwa.nip.io/v1/namespaces/:namespace/statefulsets/:name", workloadPlusServer.UpdateStatefulSet)
		group.DELETE("/apis/nuwa.nip.io/v1/namespaces/:namespace/statefulsets/:name", workloadPlusServer.DeleteStatefulSet)
	}

	//injector
	{
		group.GET("/apis/nuwa.nip.io/v1/injectors", workloadPlusServer.ListInjector)
		group.GET("/apis/nuwa.nip.io/v1/namespaces/:namespace/injectors", workloadPlusServer.ListInjector)
		group.GET("/apis/nuwa.nip.io/v1/namespaces/:namespace/injectors/:name", workloadPlusServer.GetInjector)
		group.POST("/apis/nuwa.nip.io/v1/injectors", workloadPlusServer.ApplyInjector)
		group.PUT("/apis/nuwa.nip.io/v1/namespaces/:namespace/injectors/:name", workloadPlusServer.UpdateInjector)
		group.DELETE("/apis/nuwa.nip.io/v1/namespaces/:namespace/injectors/:name", workloadPlusServer.DeleteInjector)

	}
	_ = group
	return workloadPlusServer
}
