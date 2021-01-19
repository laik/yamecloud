package tekton

import (
	"fmt"
	"github.com/yametech/yamecloud/pkg/action/api"
	"github.com/yametech/yamecloud/pkg/action/service/tekton"
)

type tektonServer struct {
	name string
	*api.Server
	// action services
	*tekton.Pipeline
	*tekton.PipelineRun
	*tekton.PipelineResource
	*tekton.Task
	*tekton.TaskRun
	*tekton.TektonStore
	*tekton.TektonWebHook
	*tekton.TektonGraph
}

func (s *tektonServer) Name() string {
	return s.name
}

func NewTektonServer(serviceName string, server *api.Server) *tektonServer {
	tektonServer := &tektonServer{
		name:   serviceName,
		Server: server,
		// action service
		Pipeline:         tekton.NewPipeline(server.Interface),
		PipelineRun:      tekton.NewPipelineRun(server.Interface),
		PipelineResource: tekton.NewPipelineResource(server.Interface),
		Task:             tekton.NewTask(server.Interface),
		TaskRun:          tekton.NewTaskRun(server.Interface),
		TektonStore:      tekton.NewTektonStore(server.Interface),
		TektonWebHook:    tekton.NewTektonWebHook(server.Interface),
		TektonGraph:      tekton.NewTektonGraph(server.Interface),
	}
	group := tektonServer.Group(fmt.Sprintf("/%s", serviceName))

	// tekton
	// pipeline
	{
		group.GET("/apis/tekton.dev/v1alpha1/pipelines", tektonServer.ListPipeline)
		group.GET("/apis/tekton.dev/v1alpha1/namespaces/:namespace/pipelines", tektonServer.ListPipeline)
		group.GET("/apis/tekton.dev/v1alpha1/namespaces/:namespace/pipelines/:name", tektonServer.GetPipeline)
		group.POST("/apis/tekton.dev/v1alpha1/namespaces/:namespace/pipelines", tektonServer.ApplyPipeline)
		group.POST("/apis/tekton.dev/v1alpha1/namespaces/:namespace/pipelines/:name/run", tektonServer.RunPipeline)
		group.DELETE("/apis/tekton.dev/v1alpha1/namespaces/:namespace/pipelines/:name", tektonServer.DeletePipeline)
	}

	// pipelineRun
	{
		group.GET("/apis/tekton.dev/v1alpha1/pipelineruns", tektonServer.ListPipelineRun)
		group.GET("/apis/tekton.dev/v1alpha1/namespaces/:namespace/pipelineruns", tektonServer.ListPipelineRun)
		group.GET("/apis/tekton.dev/v1alpha1/namespaces/:namespace/pipelineruns/:name", tektonServer.GetPipelineRun)
		group.POST("/apis/tekton.dev/v1alpha1/namespaces/:namespace/pipelineruns", tektonServer.ApplyPipelineRun)
		group.POST("/apis/tekton.dev/v1alpha1/namespaces/:namespace/pipelineruns/:name/rerun", tektonServer.RerunPipelineRun)
		group.DELETE("/apis/tekton.dev/v1alpha1/namespaces/:namespace/pipelineruns/:name", tektonServer.DeletePipelineRun)
	}

	// pipelineResource
	{
		group.GET("/apis/tekton.dev/v1alpha1/pipelineresources", tektonServer.ListPipelineResource)
		group.GET("/apis/tekton.dev/v1alpha1/namespaces/:namespace/pipelineresources", tektonServer.ListPipelineResource)
		group.GET("/apis/tekton.dev/v1alpha1/namespaces/:namespace/pipelineresources/:name", tektonServer.GetPipelineResource)
		group.POST("/apis/tekton.dev/v1alpha1/namespaces/:namespace/pipelineresources", tektonServer.ApplyPipelineResource)
		group.DELETE("/apis/tekton.dev/v1alpha1/namespaces/:namespace/pipelineresources/:name", tektonServer.DeletePipelineResource)
	}

	// task
	{
		group.GET("/apis/tekton.dev/v1alpha1/tasks", tektonServer.ListTask)
		group.GET("/apis/tekton.dev/v1alpha1/namespaces/:namespace/tasks", tektonServer.ListTask)
		group.GET("/apis/tekton.dev/v1alpha1/namespaces/:namespace/tasks/:name", tektonServer.GetTask)
		group.POST("/apis/tekton.dev/v1alpha1/namespaces/:namespace/tasks", tektonServer.ApplyTask)
		group.DELETE("/apis/tekton.dev/v1alpha1/namespaces/:namespace/tasks/:name", tektonServer.DeleteTask)
	}

	// taskRun
	{
		group.GET("/apis/tekton.dev/v1alpha1/taskruns", tektonServer.ListTaskRun)
		group.GET("/apis/tekton.dev/v1alpha1/namespaces/:namespace/taskruns", tektonServer.ListTaskRun)
		group.GET("/apis/tekton.dev/v1alpha1/namespaces/:namespace/taskruns/:name", tektonServer.GetTaskRun)
		group.POST("/apis/tekton.dev/v1alpha1/namespaces/:namespace/taskruns", tektonServer.ApplyTask)
		group.DELETE("/apis/tekton.dev/v1alpha1/namespaces/:namespace/taskruns/:name", tektonServer.DeleteTaskRun)
	}

	// tektonStore
	{
		group.GET("/apis/fuxi.nip.io/v1/tektonstores", tektonServer.ListTektonStore)
		group.GET("/apis/fuxi.nip.io/v1/namespaces/:namespace/tektonstores", tektonServer.ListTektonStore)
		group.GET("/apis/fuxi.nip.io/v1/namespaces/:namespace/tektonstores/:name", tektonServer.GetTektonStore)
		group.POST("/apis/fuxi.nip.io/v1/namespaces/:namespace/tektonstores", tektonServer.ApplyTektonStore)
		group.DELETE("/apis/fuxi.nip.io/v1/namespaces/:namespace/tektonstores/:name", tektonServer.DeleteTektonStore)
	}

	// tektonWebhook
	{
		group.GET("/apis/fuxi.nip.io/v1/tektonwebhooks", tektonServer.ListTektonWebHook)
		group.GET("/apis/fuxi.nip.io/v1/namespaces/:namespace/tektonwebhooks", tektonServer.ListTektonWebHook)
		group.GET("/apis/fuxi.nip.io/v1/namespaces/:namespace/tektonwebhooks/:name", tektonServer.GetTektonWebHook)
		group.POST("/apis/fuxi.nip.io/v1/namespaces/:namespace/tektonwebhooks", tektonServer.ApplyTektonWebHook)
		group.DELETE("/apis/fuxi.nip.io/v1/namespaces/:namespace/tektonwebhooks/:name", tektonServer.DeleteTektonWebHook)
	}

	// tektonGraph
	{
		group.GET("/apis/fuxi.nip.io/v1/tektongraphs", tektonServer.ListTektonGraph)
		group.GET("/apis/fuxi.nip.io/v1/namespaces/:namespace/tektongraphs", tektonServer.ListTektonGraph)
		group.GET("/apis/fuxi.nip.io/v1/namespaces/:namespace/tektongraphs/:name", tektonServer.GetTektonGraph)
		group.POST("/apis/fuxi.nip.io/v1/namespaces/:namespace/tektongraphs", tektonServer.ApplyTektonGraph)
		group.DELETE("/apis/fuxi.nip.io/v1/namespaces/:namespace/tektongraphs/:name", tektonServer.DeleteTektonGraph)
	}

	return tektonServer
}
