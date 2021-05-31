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

# environments
SUBLIST="";SUBTOPIC=tekton,ovn

SUBTOPIC=[tekton,ovn,istio] or *
*/

const serviceName = "watcher"
const version = "latest"

var defaultResources = []string{
	k8s.Water,
	k8s.StatefulSet1,
	k8s.Stone,
	k8s.Injector,

	k8s.Deployment,
	k8s.StatefulSet,
	k8s.DaemonSet,
	k8s.Pod,
	k8s.Job,
	k8s.CronJob,
	k8s.ReplicaSet,
	k8s.Event,
	k8s.Node,
	k8s.ConfigMap,
	k8s.Secret,
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
	k8s.Workloads,

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

var needDescSubscribeOVNList = []k8s.ResourceType{
	k8s.IP,
	k8s.SubNet,
	k8s.Vlan,
	k8s.NetworkAttachmentDefinition,
}

var needDescSubscribeTEKTONList = []k8s.ResourceType{
	k8s.Pipeline,
	k8s.PipelineRun,
	k8s.Task,
	k8s.TaskRun,
	k8s.PipelineResource,
	k8s.TektonGraph,
	k8s.OpsSecret,
	k8s.TektonWebHook,
	k8s.TektonStore,
}

var needDescSubscribeISTIOList = []k8s.ResourceType{
	k8s.Gateway,
	k8s.DestinationRule,
	k8s.ServiceEntry,
	k8s.Sidecar,
	k8s.VirtualService,
	k8s.WorkloadEntry,
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

func main() {
	subscribeListString := os.Getenv("SUBLIST")
	subscribeTopicString := os.Getenv("SUBTOPIC")

	subscribeList = append(subscribeList, defaultResources...)

	if subscribeListString == "*" {
		subscribeList = append(subscribeList, subscribeMapList()...)
	} else if subscribeListString != "" {
		needSubscribeList := strings.Split(subscribeListString, ",")
		subscribeList = append(subscribeList, subscribeMapList(needSubscribeList...)...)
	}

	if subscribeTopicString == "*" {
		subscribeList = append(subscribeList, subscribeMapList(needDescSubscribeOVNList...)...)
		subscribeList = append(subscribeList, subscribeMapList(needDescSubscribeTEKTONList...)...)
		subscribeList = append(subscribeList, subscribeMapList(needDescSubscribeISTIOList...)...)
	} else if subscribeTopicString != "" {
		topicList := strings.Split(subscribeTopicString, ",")
		for _, topic := range topicList {
			switch topic {
			case "tekton":
				subscribeList = append(subscribeList, subscribeMapList(needDescSubscribeTEKTONList...)...)
			case "istio":
				subscribeList = append(subscribeList, subscribeMapList(needDescSubscribeISTIOList...)...)
			case "ovn":
				subscribeList = append(subscribeList, subscribeMapList(needDescSubscribeOVNList...)...)
			}
		}
	}
	subscribeList = unique(subscribeList)

	fmt.Printf("[INFO] service watch resource: %v\n", subscribeList)

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

func unique(src []string) []string {
	keys := make(map[string]bool)
	list := make([]string, 0)
	for _, entry := range src {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
