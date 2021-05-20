package workloads

import (
	"fmt"
	"github.com/yametech/yamecloud/pkg/action/api"
	workload_service "github.com/yametech/yamecloud/pkg/action/service/workloads"
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
}

func (s *workloadServer) Name() string {
	return s.name
}

func NewWorkloadServer(serviceName string, server *api.Server) *workloadServer {
	workloadServer := &workloadServer{
		name:      serviceName,
		Server:    server,
		ConfigMap: workload_service.NewConfigMap(server),
		CronJob:   workload_service.NewCronJob(server),
		DaemonSet: workload_service.NewDaemonSet(server),

		Event: workload_service.NewEvent(server),
		HPA:   workload_service.NewHPA(server),
		Job:   workload_service.NewJob(server),
		Pod:   workload_service.NewPod(server),

		ReplicaSet:    workload_service.NewReplicaSet(server),
		ResourceQuota: workload_service.NewResourceQuota(server),

		Secret:      workload_service.NewSecret(server),
		StatefulSet: workload_service.NewStatefulSet(server),
	}

	group := workloadServer.Group(fmt.Sprintf("/%s", serviceName))

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

	_ = group
	return workloadServer
}
