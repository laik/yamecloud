package main

import (
	"fmt"

	"github.com/yametech/yamecloud/pkg/action/api"
	apiService "github.com/yametech/yamecloud/pkg/action/api/service"
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/configure"
	"github.com/yametech/yamecloud/pkg/install"
	"github.com/yametech/yamecloud/pkg/k8s"
	"github.com/yametech/yamecloud/pkg/k8s/datasource"
	"github.com/yametech/yamecloud/pkg/k8s/types"
)

/*
#if the app runtime in kubernetes
export IN_CLUSTER=true

#if use etcd discover server
#argument additions
--registry etcd --registry_address ${etcd_addr}
*/

const serviceName = "service"
const version = "latest"

var subscribeList = k8s.GVRMaps.Subscribe(
	k8s.Service,
	k8s.Ingress,
	k8s.Endpoint,
)

func main() {
	config, err := configure.NewInstallConfigure(types.NewResourceITypes(subscribeList))
	if err != nil {
		panic(fmt.Sprintf("new install configure error %s", err))
	}

	_datasource := datasource.NewInterface(config)
	apiServer := api.NewServer(service.NewService(_datasource))
	apiServer.SetExtends(apiService.NewServiceServer(serviceName, apiServer))

	microService, err := install.WebServiceInstall(serviceName, version, _datasource, apiServer)
	if err != nil {
		panic(fmt.Sprintf("web service install error %s", err))
	}
	if err := microService.Run(); err != nil {
		panic(err)
	}
}
