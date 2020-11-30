package main

import (
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
	config, err := configure.NewInstallConfigure(
		types.NewResourceLister([]k8s.ResourceType{}...),
	)
	if err != nil {
		panic(err)
	}
	svc := service.NewService(datasource.NewInterface(config))
	app := workload.NewWorkloadServer(
		api.NewAPIServer(),
		dac.NewClusterRole(svc),
	)
	if err := app.Run(":8080"); err != nil {
		panic(err)
	}
}
