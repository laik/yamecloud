package main

import (
	"flag"
	"github.com/yametech/yamecloud/pkg/action/api"
	"github.com/yametech/yamecloud/pkg/action/api/workload"
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/action/service/dac"
	"github.com/yametech/yamecloud/pkg/configure"
	"github.com/yametech/yamecloud/pkg/install"
	"github.com/yametech/yamecloud/pkg/k8s"
	"github.com/yametech/yamecloud/pkg/k8s/datasource"
	"github.com/yametech/yamecloud/pkg/k8s/types"
)

const serviceName = "workload"
const version = "latest"

var subscribeList = k8s.GVRMaps.Subscribe(
	k8s.Deployment,
	k8s.StatefulSet,
	k8s.ClusterRole,
)

func main() {
	flag.Parse()
	config, err := configure.NewInstallConfigure(types.NewResourceITypes(subscribeList))
	if err != nil {
		panic(err)
	}
	_datasource := datasource.NewInterface(config)
	actionService := service.NewService(_datasource)

	apiServer := api.NewServer()
	apiServer.SetIResourceServiceMaps(
		api.IResourceServiceMaps{
			k8s.ClusterRole: dac.NewClusterRole(actionService),
		},
	)
	apiServer.SetExtends(workload.NewWorkloadServer(serviceName, apiServer))

	microService, err := install.WebServiceInstall(serviceName, version, _datasource, apiServer)
	if err != nil {
		panic(err)
	}
	if err := microService.Run(); err != nil {
		panic(err)
	}
}
