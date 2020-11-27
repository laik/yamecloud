package k8s

type ResourceType string

const (
	//// CI && CD YameCloudExtensions resources
	//CI = "cis"
	//CD = "cds"

	// Tekton resources
	Pipeline         ResourceType = "pipelines"
	PipelineRun      ResourceType = "pipelineruns"
	Task             ResourceType = "tasks"
	TaskRun          ResourceType = "taskruns"
	PipelineResource ResourceType = "pipelineresources"

	// Extend Tekton Pipeline & PipelineRun resource Graph
	TektonGraph  ResourceType = "tektongraphs"
	TektonConfig ResourceType = "secrets"

	// Stone deployment resource
	Stone ResourceType = "stones"

	// Kubernetes
	ServiceAccount ResourceType = "serviceaccounts"
	Namespace      ResourceType = "namespaces"
	Pod            ResourceType = "pods"

	// App v1
	Deployment  ResourceType = "deployments"
	Statefulset ResourceType = "statefulsets"
)

var (
	// TektonAlphaV1 tekton resources
	TektonAlphaV1 = []ResourceType{
		Pipeline,
		PipelineRun,
		Task,
		TaskRun,
		PipelineResource,
	}

	// KubernetesAppV1 App v1 resources
	KubernetesAppV1 = []ResourceType{
		Deployment,
		Statefulset,
	}

	// KubernetesV1 Native resources
	KubernetesV1 = []ResourceType{
		"pods",
		"namespaces",
		"serviceaccounts",
	}
)
