package main

import (
	"fmt"
	"github.com/yametech/yamecloud/pkg/action/api"
	"github.com/yametech/yamecloud/pkg/action/api/editer"
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

const serviceName = "editer"
const version = "latest"

var subscribeList = k8s.GVRMaps.Subscribe()

func main() {
	config, err := configure.NewInstallConfigure(types.NewResourceITypes(subscribeList))
	if err != nil {
		panic(fmt.Sprintf("new install configure error %s", err))
	}

	_datasource := datasource.NewInterface(config)
	apiServer := api.NewServer(service.NewService(_datasource))
	apiServer.SetExtends(editer.NewEditerServer(serviceName, apiServer))

	microService, err := install.WebServiceInstall(serviceName, version, _datasource, apiServer)
	if err != nil {
		panic(fmt.Sprintf("web service install error %s", err))
	}
	if err := microService.Run(); err != nil {
		panic(err)
	}
}
