package main

import (
	"fmt"
	"github.com/yametech/yamecloud/pkg/action/api/tekton"

	"github.com/yametech/yamecloud/pkg/action/api"
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

const serviceName = "tekton"
const version = "latest"

var subscribeList = k8s.GVRMaps.Subscribe(
	k8s.Pipeline,
	k8s.PipelineRun,
	k8s.PipelineResource,
	k8s.Task,
	k8s.TaskRun,
	k8s.TektonStore,
	k8s.TektonWebHook,
)

func main() {
	config, err := configure.NewInstallConfigure(types.NewResourceITypes(subscribeList))
	if err != nil {
		panic(fmt.Sprintf("new install configure error %s", err))
	}

	_datasource := datasource.NewInterface(config)
	apiServer := api.NewServer(service.NewService(_datasource))
	apiServer.SetExtends(tekton.NewTektonServer(serviceName, apiServer))

	microService, err := install.WebServiceInstall(serviceName, version, _datasource, apiServer)
	if err != nil {
		panic(fmt.Sprintf("web service install error %s", err))
	}
	if err := microService.Run(); err != nil {
		panic(err)
	}
}
