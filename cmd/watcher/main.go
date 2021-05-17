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
	"os"
	"strings"
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

var defaultResources = []string{
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

	// tenant for PaaS
	k8s.BaseDepartment,
	k8s.BaseTenant,
	k8s.BaseRole,
	k8s.BaseUser,
	k8s.BaseRoleUser,
}

var needDescSubscribeMap = map[string]k8s.ResourceType{
	"ip":                          k8s.IP,
	"subnet":                      k8s.SubNet,
	"vlan":                        k8s.Vlan,
	"networkAttachmentDefinition": k8s.NetworkAttachmentDefinition,

	//tekton
	"pipeline":         k8s.Pipeline,
	"pipelinerun":      k8s.PipelineRun,
	"task":             k8s.Task,
	"taskrun":          k8s.TaskRun,
	"pipelineresource": k8s.PipelineResource,
	"tektongraph":      k8s.TektonGraph,
	"opssecret":        k8s.OpsSecret,
	"tektonwebhook":    k8s.TektonWebHook,
	"tektonstore":      k8s.TektonStore,

	//Istio  NetWorking
	"gateway":         k8s.Gateway,
	"destinationrule": k8s.DestinationRule,
	"serviceentry":    k8s.ServiceEntry,
	"sidecar":         k8s.Sidecar,
	"virtualservice":  k8s.VirtualService,
	"workloadentry":   k8s.WorkloadEntry,
}

func subscribeMapList(includes ...string) []string {
	result := make([]string, 0)
	for _, v := range needDescSubscribeMap {
		pass := false
		if len(includes) != 0 {
			for _, item := range includes {
				if v == item {
					pass = true
				}
			}
		}
		if len(includes) != 0 && !pass {
			continue
		}
		result = append(result, v)
	}

	return result
}

var subscribeList []string

var subscribeListString string

func main() {
	subscribeListString := os.Getenv("SUBLIST")

	subscribeList = append(subscribeList, defaultResources...)

	if subscribeListString == "*" {
		subscribeList = append(subscribeList, subscribeMapList()...)
	} else if subscribeListString != "" {
		subscribeList = strings.Split(subscribeListString, ",")
		subscribeList = append(subscribeList, subscribeMapList(subscribeList...)...)
	}

	fmt.Printf("service watch resource: %v\n", subscribeList)

	resourceList := k8s.GVRMaps.Subscribe(subscribeList...)

	config, err := configure.NewInstallConfigure(types.NewResourceITypes(resourceList))
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
