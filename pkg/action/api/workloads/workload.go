package workloads

import (
	"fmt"
	"github.com/yametech/yamecloud/pkg/action/api"
	"github.com/yametech/yamecloud/pkg/action/service/editer"
	workload_service "github.com/yametech/yamecloud/pkg/action/service/workloads"
	"github.com/yametech/yamecloud/pkg/helm"
	"helm.sh/helm/v3/pkg/action"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type workloadServer struct {
	name string
	*api.Server
	// action
	*workload_service.ConfigMap
	*workload_service.CronJob
	*workload_service.DaemonSet
	*workload_service.Deployment
	*workload_service.Event
	*workload_service.HPA
	*workload_service.Job
	*workload_service.Pod
	*workload_service.ReplicaSet
	*workload_service.ResourceQuota
	*workload_service.Secret
	*workload_service.StatefulSet
	*workload_service.Metrics
	*workload_service.Namespace
	*workload_service.Node
	*workload_service.CRD

	*workload_service.PV
	*workload_service.PVC
	*workload_service.StorageClass

	*workload_service.Template
	*editer.APIResources
	HelmAction helm.NewActionConfigWithSecretFunc

	config    *rest.Config
	clientset *kubernetes.Clientset
}

func (w *workloadServer) Name() string {
	return w.name
}

func NewWorkloadServer(serviceName string, server *api.Server, config *rest.Config, clientset *kubernetes.Clientset) *workloadServer {
	workloadServer := &workloadServer{
		name:      serviceName,
		Server:    server,
		config:    config,
		clientset: clientset,

		ConfigMap:  workload_service.NewConfigMap(server),
		CronJob:    workload_service.NewCronJob(server),
		DaemonSet:  workload_service.NewDaemonSet(server),
		Deployment: workload_service.NewDeployment(server),

		Event: workload_service.NewEvent(server),
		HPA:   workload_service.NewHPA(server),
		Job:   workload_service.NewJob(server),
		Pod:   workload_service.NewPod(server),

		ReplicaSet:    workload_service.NewReplicaSet(server),
		ResourceQuota: workload_service.NewResourceQuota(server),

		Secret:      workload_service.NewSecret(server),
		StatefulSet: workload_service.NewStatefulSet(server),
		Metrics:     workload_service.NewMetrics(server),

		Namespace: workload_service.NewNamespace(server),
		Node:      workload_service.NewNode(server),
		CRD:       workload_service.NewCRD(server),

		PV:  workload_service.NewPV(server),
		PVC: workload_service.NewPVC(server),

		StorageClass: workload_service.NewStorageClass(server),

		Template:     workload_service.NewTemplate(server),
		APIResources: editer.NewAPIResources(server),

		HelmAction: func(namespace string) *action.Configuration {
			cfg, err := helm.NewActionConfigWithSecret(config, clientset, namespace)
			if err != nil {
				panic(err)
			}
			return cfg
		},
	}

	group := workloadServer.Group(fmt.Sprintf("/%s", serviceName))

	//accesscontrol
	{
		group.POST("/accesscontrol", workloadServer.Stack)
	}
	//template
	{
		group.GET("/apis/yamecloud.io/v1/workloads", workloadServer.ListTemplate)
		group.POST("/apis/yamecloud.io/v1/workloads", workloadServer.CreateTemplate)
		group.PUT("/apis/yamecloud.io/v1/workloads/:name", workloadServer.UpdateTemplate)
		group.PUT("/apis/yamecloud.io/v1/workloads/:name/label", workloadServer.LabelTemplate)
		group.DELETE("/apis/yamecloud.io/v1/workloads/:name", workloadServer.DeleteTemplate)
		group.POST("/apis/yamecloud.io/v1/workloads/:name/deploy", workloadServer.DeployTemplate)
	}

	// configmap api
	{
		group.GET("/api/v1/configmaps", workloadServer.ListConfigMap)
		group.GET("/api/v1/namespaces/:namespace/configmaps", workloadServer.ListConfigMap)
		group.GET("/api/v1/namespaces/:namespace/configmaps/:name", workloadServer.GetConfigMap)
		group.POST("/api/v1/namespaces/:namespace/configmaps", workloadServer.ApplyConfigMap)
		group.PUT("/api/v1/namespaces/:namespace/configmaps/:name", workloadServer.UpdateConfigMap)
		group.DELETE("/api/v1/namespaces/:namespace/configmaps/:name", workloadServer.DeleteConfigMap)
	}

	// pod api
	{
		group.GET("/api/v1/pods", workloadServer.ListPod)
		group.GET("/api/v1/namespaces/:namespace/pods", workloadServer.ListPod)
		group.GET("/api/v1/namespaces/:namespace/pods/:name", workloadServer.GetPod)
		group.GET("/api/v1/namespaces/:namespace/pods/:name/log", workloadServer.LogsPod)
		group.DELETE("/api/v1/namespaces/:namespace/pods/:name", workloadServer.DeletePod)
	}

	// secret api
	{
		group.GET("/api/v1/secrets", workloadServer.ListSecret)
		group.GET("/api/v1/namespaces/:namespace/secrets", workloadServer.ListSecret)
		group.GET("/api/v1/namespaces/:namespace/secrets/:name", workloadServer.GetSecret)
		group.POST("/api/v1/namespaces/:namespace/secrets", workloadServer.ApplySecret)
		group.PUT("/api/v1/namespaces/:namespace/secrets/:name", workloadServer.UpdateSecret)
		group.DELETE("/api/v1/namespaces/:namespace/secrets/:name", workloadServer.DeleteSecret)
	}

	// hpa api
	{
		group.GET("/apis/autoscaling/v2beta1/horizontalpodautoscalers", workloadServer.ListHPA)
		group.GET("/apis/autoscaling/v2beta1/namespaces/:namespace/horizontalpodautoscalers", workloadServer.ListHPA)
		group.GET("/apis/autoscaling/v2beta1/namespaces/:namespace/horizontalpodautoscalers/:name", workloadServer.GetHPA)
		group.POST("/apis/autoscaling/v2beta1/namespaces/:namespace/horizontalpodautoscalers", workloadServer.ApplyHPA)
		group.PUT("/apis/autoscaling/v2beta1/namespaces/:namespace/horizontalpodautoscalers/:name", workloadServer.UpdateHPA)
		group.DELETE("/apis/autoscaling/v2beta1/namespaces/:namespace/horizontalpodautoscalers/:name", workloadServer.DeleteHPA)
	}

	// deployment api
	{
		group.GET("/apis/apps/v1/deployments", workloadServer.ListDeployment)
		group.GET("/apis/apps/v1/namespaces/:namespace/deployments", workloadServer.ListDeployment)
		group.GET("/apis/apps/v1/namespaces/:namespace/deployments/:name", workloadServer.GetDeployment)
		group.POST("/apis/apps/v1/namespaces/:namespace/deployments", workloadServer.ApplyDeployment)
		group.PUT("/apis/apps/v1/namespaces/:namespace/deployments/:name", workloadServer.UpdateDeployment)
		group.DELETE("/apis/apps/v1/namespaces/:namespace/deployments/:name", workloadServer.DeleteDeployment)
		// scale
		group.GET("/apis/apps/v1/namespaces/:namespace/deployments/:name/scale", workloadServer.DeploymentScaleInfo)
		group.PUT("/apis/apps/v1/namespaces/:namespace/deployments/:name/scale", workloadServer.DeploymentScale)
	}

	// daemonset api
	{
		group.GET("/apis/apps/v1/daemonsets", workloadServer.ListDaemonSet)
		group.GET("/apis/apps/v1/namespaces/:namespace/daemonsets", workloadServer.ListDaemonSet)
		group.GET("/apis/apps/v1/namespaces/:namespace/daemonsets/:name", workloadServer.GetDaemonSet)
		group.POST("/apis/apps/v1/namespaces/:namespace/daemonsets", workloadServer.ApplyDaemonSet)
		group.PUT("/apis/apps/v1/namespaces/:namespace/daemonsets/:name", workloadServer.UpdateDaemonSet)
		group.DELETE("/apis/apps/v1/namespaces/:namespace/daemonsets/:name", workloadServer.DeleteDaemonSet)
	}

	// statefulset api
	{
		group.GET("/apis/apps/v1/statefulsets", workloadServer.ListStatefulSet)
		group.GET("/apis/apps/v1/namespaces/:namespace/statefulsets", workloadServer.ListStatefulSet)
		group.GET("/apis/apps/v1/namespaces/:namespace/statefulsets/:name", workloadServer.GetStatefulSet)
		group.POST("/apis/apps/v1/namespaces/:namespace/statefulsets", workloadServer.ApplyStatefulSet)
		group.PUT("/apis/apps/v1/namespaces/:namespace/statefulsets/:name", workloadServer.UpdateStatefulSet)
		group.DELETE("/apis/apps/v1/namespaces/:namespace/statefulsets/:name", workloadServer.DeleteStatefulSet)
	}
	// job api
	{
		group.GET("/apis/batch/v1/jobs", workloadServer.ListJob)
		group.GET("/apis/batch/v1/namespaces/:namespace/jobs", workloadServer.ListJob)
		group.GET("/apis/batch/v1/namespaces/:namespace/jobs/:name", workloadServer.GetJob)
		group.POST("/apis/batch/v1/namespaces/:namespace/jobs", workloadServer.ApplyJob)
		group.PUT("/apis/batch/v1/namespaces/:namespace/jobs/:name", workloadServer.UpdateJob)
		group.DELETE("/apis/batch/v1/namespaces/:namespace/jobs/:name", workloadServer.DeleteJob)
	}

	// cronjob api
	{
		group.GET("/apis/batch/v1beta1/cronjobs", workloadServer.ListCronJob)
		group.GET("/apis/batch/v1beta1/namespaces/:namespace/cronjobs", workloadServer.ListCronJob)
		group.GET("/apis/batch/v1beta1/namespaces/:namespace/cronjobs/:name", workloadServer.GetCronJob)
		group.POST("/apis/batch/v1beta1/namespaces/:namespace/cronjobs", workloadServer.ApplyCronJob)
		group.PUT("/apis/batch/v1beta1/namespaces/:namespace/cronjobs/:name", workloadServer.UpdateCronJob)
		group.DELETE("/apis/batch/v1beta1/namespaces/:namespace/cronjobs/:name", workloadServer.DeleteCronJob)
	}

	// metrics
	{
		group.POST("/metrics", workloadServer.DefaultMetrics)
		group.GET("/apis/metrics.k8s.io/v1beta1/nodes", workloadServer.ListNodeMetrics)
		group.GET("/apis/metrics.k8s.io/v1beta1/pods", workloadServer.ListPodMetrics)
		group.GET("/apis/metrics.k8s.io/v1beta1/namespaces/:namespace/pods", workloadServer.GetPodMetrics)
	}

	// event
	{
		group.GET("/api/v1/events", workloadServer.ListEvent)
		group.GET("/api/v1/namespaces/:namespace/events", workloadServer.ListEvent)
		group.GET("/api/v1/namespaces/:namespace/events/:name", workloadServer.GetEvent)
	}

	// namespace
	{
		group.GET("/api/v1/namespaces", workloadServer.ListNamespace)
		group.GET("/api/v1/namespaces/:namespace", workloadServer.GetNamespace)
		group.POST("/api/v1/namespaces", workloadServer.ApplyNamespace)
		group.DELETE("/api/v1/namespaces/:namespace", workloadServer.DeleteNamespace)

		group.POST("/api/v1/namespaces/:namespace/annotate/node", workloadServer.AnnotateNamespaceAllowedNode)
		group.POST("/api/v1/namespaces/:namespace/annotate/networkattachment", workloadServer.AnnotateNamespaceNetworkAttach)
		group.POST("/api/v1/namespaces/:namespace/annotate/storageclass", workloadServer.AnnotateNamespaceAllowedStorageClass)
	}

	// resource quota
	{
		group.GET("/api/v1/resourcequotas", workloadServer.ListResourceQuota)
		group.GET("/api/v1/namespaces/:namespace/resourcequotas", workloadServer.ListResourceQuota)
		group.GET("/api/v1/namespaces/:namespace/resourcequotas/:name", workloadServer.GetResourceQuota)
		group.POST("/api/v1/namespaces/:namespace/resourcequotas", workloadServer.ApplyResourceQuota)
		group.DELETE("/api/v1/namespaces/:namespace/resourcequotas/:name", workloadServer.DeleteResourceQuota)
	}

	// node
	{
		group.GET("/api/v1/nodes", workloadServer.ListNode)
		group.GET("/api/v1/nodes/:name", workloadServer.GetNode)
		group.POST("/api/v1/nodes", workloadServer.ApplyNode)
		group.DELETE("/api/v1/nodes/:name", workloadServer.DeleteNode)
		group.POST("/api/v1/nodes/:name/annotate/geo", workloadServer.NodeGEO)
	}

	// replicasets
	{
		group.GET("apis/apps/v1/replicasets", workloadServer.ListReplicaSet)
	}

	// PersistentVolume
	{
		group.GET("/api/v1/persistentvolumes", workloadServer.ListPV)
		group.GET("/api/v1/persistentvolumes/:name", workloadServer.GetPV)
		group.DELETE("/api/v1/persistentvolumes/:name", workloadServer.DeletePV)
	}

	// PersistentVolumeClaims
	{
		group.GET("/api/v1/persistentvolumeclaims", workloadServer.ListPVC)
		group.GET("/api/v1/namespaces/:namespace/persistentvolumeclaims", workloadServer.ListPVC)
		group.GET("/api/v1/namespaces/:namespace/persistentvolumeclaims/:name", workloadServer.GetPVC)
		group.DELETE("/api/v1/namespaces/:namespace/persistentvolumeclaims/:name", workloadServer.DeletePVC)
	}

	// #storage.k8s.io
	{
		group.GET("/apis/storage.k8s.io/v1/storageclasses", workloadServer.ListStorageClass)
		group.GET("/apis/storage.k8s.io/v1/storageclasses/:name", workloadServer.GetStorageClass)
		group.POST("/apis/storage.k8s.io/v1/storageclasses", workloadServer.ApplyStorageClass)
		group.PUT("/apis/storage.k8s.io/v1/storageclasses/:name", workloadServer.UpdateStorageClass)
		group.DELETE("/apis/storage.k8s.io/v1/storageclasses/:name", workloadServer.DeleteStorageClass)
	}

	// helm
	{

		group.GET("/v2/charts", workloadServer.ListCharts)
		group.GET("/v2/charts/:repo/:chart", workloadServer.GetCharts)
		group.GET("/v2/charts/:repo/:chart/values", workloadServer.GetChartValues)
		//

		group.POST("/v2/releases", workloadServer.InstallChart)

		group.GET("/v2/releases", workloadServer.ListRelease)
		group.GET("/v2/releases/:namespace", workloadServer.ReleasesByNamespace)
		group.GET("/v2/releases/:namespace/:release", workloadServer.ReleaseByNamespace)

		group.GET("/v2/releases/:namespace/:release/values", workloadServer.ReleaseValueByName)

		group.DELETE("/v2/releases/:namespace/:release", workloadServer.DeleteRelease)
		group.PUT("/v2/releases/:namespace/:release", workloadServer.UpgradeRelease)
		group.PUT("/v2/releases/:namespace/:release/rollback", workloadServer.RollbackRelease)
		group.GET("/v2/releases/:namespace/:release/history", workloadServer.HistoryRelease)

	}

	// CRD
	// #apiextensions.k8s.io/v1beta1
	// #v1beta1, > 1.16 v1
	// CustomResourceDefinition
	{
		group.GET("/apis/apiextensions.k8s.io/v1/customresourcedefinitions", workloadServer.ListCustomResourceDefinition)

		ignores := []string{
			"yamecloud.io/v1/workloads",
			"yamecloud.io/v1/basedepartments",
			"yamecloud.io/v1/baseroles",
			"yamecloud.io/v1/baseusers",

			"yamecloud.io/v1/tektongraphs",
			"yamecloud.io/v1/tektonwebhooks",
			"yamecloud.io/v1/tektonstores",

			"tekton.dev/v1alpha1/pipelines",
			"tekton.dev/v1alpha1/pipelineruns",
			"tekton.dev/v1alpha1/pipelineresources",
			"tekton.dev/v1alpha1/tasks",
			"tekton.dev/v1alpha1/taskruns",

			"kubeovn.io/v1/ips",
			"kubeovn.io/v1/subnets",
			"kubeovn.io/v1/vlans",
		}

		apiVersions, err := workloadServer.ListCustomResourceRouter(ignores)
		if err != nil {
			panic(err)
		}
		routerPath := "/apis/%s"
		for _, apiVersion := range apiVersions {
			relativePath := fmt.Sprintf(routerPath, apiVersion)
			group.GET(relativePath, workloadServer.ListGeneralCustomResourceDefinition)
		}
	}

	return workloadServer
}
