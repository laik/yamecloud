package main

import (
	"flag"
	"github.com/yametech/yamecloud/pkg/action/api"
	"github.com/yametech/yamecloud/pkg/configure"
	"github.com/yametech/yamecloud/pkg/install"
	"github.com/yametech/yamecloud/pkg/k8s"
	"github.com/yametech/yamecloud/pkg/k8s/datasource"
	"github.com/yametech/yamecloud/pkg/k8s/types"
)

const version = "latest"

var subscribeList = k8s.GVRMaps.List(
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
	apiServer := api.NewServer()
	microService, err := install.WebServiceInstall("workload", version, datasource.NewInterface(config), apiServer)
	if err != nil {
		panic(err)
	}

	_ = microService

	if err := apiServer.Run(":8080"); err != nil {
		panic(err)
	}
}
