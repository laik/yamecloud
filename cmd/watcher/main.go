package main

import (
	"fmt"

	"github.com/yametech/yamecloud/pkg/action/api"
	apiService "github.com/yametech/yamecloud/pkg/action/api/watcher"
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

const serviceName = "watcher"
const version = "latest"

var subscribeList = k8s.GVRMaps.Subscribe(
	k8s.Water,
	k8s.Deployment,
	k8s.Stone,
	k8s.StatefulSet,
	k8s.StatefulSet1,
	k8s.DaemonSet,
	k8s.Injector,
	k8s.Pod,
	k8s.Job,
	k8s.CronJobs,
	k8s.ReplicaSet,
	k8s.Event,
	k8s.Node,
	k8s.ConfigMaps,
	k8s.Secrets,
	k8s.ResourceQuota,
	k8s.Service,
	k8s.Ingress,
	k8s.NetworkPolicy,
	k8s.HorizontalPodAutoscaler,
	k8s.CustomResourceDefinition,
	k8s.PersistentVolume,
	k8s.PersistentVolumeClaims,
	k8s.StorageClass,
	k8s.ServiceAccount,
	k8s.Role,
	k8s.ClusterRole,
	k8s.RoleBinding,
	k8s.Namespace,
	k8s.PodSecurityPolicie,
	k8s.ClusterRoleBinding,
	k8s.Endpoint,

	// deployment resource workload template for CaaS
	k8s.WorkloadsTemplate,

	// form render
	//k8s.FormRender,
	//k8s.Page,
	//k8s.Form,
	//k8s.Field,

	// tenant for PaaS
	k8s.BaseDepartment,
	k8s.BaseTenant,
	k8s.BaseRole,
	k8s.BaseUser,
	k8s.BaseRoleUser,

	// network for container ovn/ovs control plant
	//k8s.IP,
	//k8s.SubNet,
	//k8s.Vlan,
	//k8s.NetworkAttachmentDefinition,

	//tekton
	k8s.Pipeline,
	k8s.PipelineRun,
	k8s.Task,
	k8s.TaskRun,
	k8s.PipelineResource,
	k8s.TektonGraph,
	k8s.OpsSecret,
	k8s.TektonWebHook,
	k8s.TektonStore,

	//Istio  NetWorking
	k8s.Gateway,
	k8s.DestinationRule,
	k8s.ServiceEntry,
	k8s.Sidecar,
	k8s.VirtualService,
	k8s.WorkloadEntry,
)

func main() {
	config, err := configure.NewInstallConfigure(types.NewResourceITypes(subscribeList))
	if err != nil {
		panic(fmt.Sprintf("new install configure error %s", err))
	}

	_datasource := datasource.NewInterface(config)
	apiServer := api.NewServer(service.NewService(_datasource))
	apiServer.SetExtends(apiService.NewWatcherServer(serviceName, apiServer))

	microService, err := install.WebServiceInstall(serviceName, version, _datasource, apiServer)
	if err != nil {
		panic(fmt.Sprintf("web service install error %s", err))
	}
	if err := microService.Run(); err != nil {
		panic(err)
	}
}
