package main

import (
	"fmt"
	"github.com/yametech/yamecloud/pkg/action/api"
	"github.com/yametech/yamecloud/pkg/action/api/workloadplus"
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/configure"
	"github.com/yametech/yamecloud/pkg/install"
	"github.com/yametech/yamecloud/pkg/k8s"
	"github.com/yametech/yamecloud/pkg/k8s/datasource"
	"github.com/yametech/yamecloud/pkg/k8s/types"
)

/*
export MICRO_SERVER_ADDRESS=0.0.0.0:8080
*/

const serviceName = "workloadplus"
const version = "latest"

var subscribeList = k8s.GVRMaps.Subscribe(
	k8s.Stone,
	k8s.StatefulSet1,
	k8s.Water,
	k8s.Injector,
)

func main() {
	config, err := configure.NewInstallConfigure(types.NewResourceITypes(subscribeList))
	if err != nil {
		panic(fmt.Sprintf("new install configure error %workloads", err))
	}

	_datasource := datasource.NewInterface(config)
	apiServer := api.NewServer(service.NewService(_datasource))
	apiServer.SetExtends(workloadplus.NewWorkloadPlusServer(serviceName, apiServer))

	microService, err := install.WebServiceInstall(serviceName, version, _datasource, apiServer)
	if err != nil {
		panic(fmt.Sprintf("web service install error %workloads", err))
	}

	if err := microService.Run(); err != nil {
		panic(err)
	}
}
