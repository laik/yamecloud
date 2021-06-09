package main

import (
	"fmt"

	"github.com/yametech/yamecloud/pkg/action/api"
	"github.com/yametech/yamecloud/pkg/action/api/workloads"
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

const serviceName = "workloads"
const version = "latest"

var subscribeList = k8s.GVRMaps.Subscribe(
	k8s.Pod,
	k8s.Deployment,
	k8s.StatefulSet,
	k8s.DaemonSet,
	k8s.Job,
	k8s.CronJob,
	k8s.ReplicaSet,

	k8s.ConfigMap,
	k8s.Secret,
	k8s.HorizontalPodAutoscaler,
	k8s.ResourceQuota,

	k8s.Namespace,
	k8s.Node,
	k8s.Event,
	k8s.CustomResourceDefinition,

	// storage
	k8s.PersistentVolume,
	k8s.PersistentVolumeClaims,
	k8s.StorageClass,

	// workloads template
	k8s.Workloads,
	k8s.Stone,

	//
	k8s.Service,
)

func main() {
	config, err := configure.NewInstallConfigure(types.NewResourceITypes(subscribeList))
	if err != nil {
		panic(fmt.Sprintf("new install configure error %s", err))
	}

	_datasource := datasource.NewInterface(config)
	apiServer := api.NewServer(service.NewService(_datasource))
	apiServer.SetExtends(workloads.NewWorkloadServer(serviceName, apiServer))

	microService, err := install.WebServiceInstall(serviceName, version, _datasource, apiServer)
	if err != nil {
		panic(fmt.Sprintf("web service install error %s", err))
	}

	if err := microService.Run(); err != nil {
		panic(err)
	}
}
