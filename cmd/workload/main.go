package main

import (
	"flag"
	"github.com/yametech/yamecloud/pkg/action/api"
	"github.com/yametech/yamecloud/pkg/action/api/workload"
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/action/service/dac"
	"github.com/yametech/yamecloud/pkg/configure"
	"github.com/yametech/yamecloud/pkg/k8s"
	"github.com/yametech/yamecloud/pkg/k8s/datasource"
	"github.com/yametech/yamecloud/pkg/k8s/types"
)

func main() {
	flag.Parse()
	config, err := configure.NewInstallConfigure(
		types.NewResourceITypes(
			k8s.GVRMaps.List(
				k8s.Deployment,
				k8s.StatefulSet,
				k8s.ClusterRole,
			)),
	)
	if err != nil {
		panic(err)
	}
	svc := service.NewService(datasource.NewInterface(config))
	app := workload.NewWorkloadServer(
		api.NewServer(),
		dac.NewClusterRole(svc),
	)
	if err := app.Run(":8080"); err != nil {
		panic(err)
	}
}
