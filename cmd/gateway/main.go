package main

import (
	"flag"
	"fmt"
	"github.com/yametech/yamecloud/pkg/action/api"
	"github.com/yametech/yamecloud/pkg/action/api/gateway"
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/configure"
	"github.com/yametech/yamecloud/pkg/install"
	"github.com/yametech/yamecloud/pkg/k8s"
	"github.com/yametech/yamecloud/pkg/k8s/datasource"
	"github.com/yametech/yamecloud/pkg/k8s/types"

	_ "github.com/micro/go-plugins/registry/etcd"
)

/*
export MICRO_SERVER_ADDRESS=0.0.0.0:8080
*/

const serviceName = "gateway"

var subscribeList = k8s.GVRMaps.Subscribe(
	//k8s.BaseTenant,
	k8s.BaseRole,
	k8s.BaseUser,
	k8s.BaseRoleUser,
)

func main() {
	flag.Parse()

	config, err := configure.NewInstallConfigure(types.NewResourceITypes(subscribeList))
	if err != nil {
		panic(fmt.Sprintf("new install configure error %s", err))
	}

	_datasource := datasource.NewInterface(config)
	apiServer := api.NewServer(service.NewService(_datasource))
	apiServer.SetExtends(gateway.NewGatewayServer(serviceName, apiServer))

	gatewayService, err := install.GatewayInstall(_datasource, apiServer)
	if err != nil {
		panic(fmt.Sprintf("gateway service install error %s", err))
	}

	if err := gatewayService.Run(); err != nil {
		panic(err)
	}
}
