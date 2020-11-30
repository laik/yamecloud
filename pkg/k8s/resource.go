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
	TektonWebHook    ResourceType = "tektonwebhooks"
	TektonStore      ResourceType = "tektonstores"
	TektonGraph      ResourceType = "tektongraphs"
	TektonConfig     ResourceType = "secrets"

	// nuwa resource
	Stone        ResourceType = "stones"
	Water        ResourceType = "water"
	StatefulSet1 ResourceType = "statefulset1"
	Injector     ResourceType = "injector"

	// dac (access control)
	ServiceAccount     ResourceType = "serviceaccounts"
	Role               ResourceType = "roles"
	ClusterRole        ResourceType = "clusterroles"
	RoleBinding        ResourceType = "rolebindings"
	ClusterRoleBinding ResourceType = "clusterrolebindings"

	// kubernetes
	Namespace                ResourceType = "namespaces"
	Pod                      ResourceType = "pods"
	Deployment               ResourceType = "deployments"
	Statefulset              ResourceType = "statefulsets"
	DaemonSet                ResourceType = "daemonsets"
	Job                      ResourceType = "jobs"
	CronJobs                 ResourceType = "cronjobs"
	ReplicaSet               ResourceType = "replicasets"
	Event                    ResourceType = "events"
	Node                     ResourceType = "nodes"
	ConfigMaps               ResourceType = "configmaps"
	Secrets                  ResourceType = "secrets"
	ResourceQuota            ResourceType = "resourcequotas"
	Service                  ResourceType = "services"
	Ingress                  ResourceType = "ingresses"
	NetworkPolicy            ResourceType = "networkpolicies"
	HorizontalPodAutoscaler  ResourceType = "horizontalpodautoscalers"
	CustomResourceDefinition ResourceType = "customresourcedefinitions"
	PersistentVolume         ResourceType = "persistentvolumes"
	PersistentVolumeClaims   ResourceType = "persistentvolumeclaims"
	StorageClass             ResourceType = "storageclasses"
	Endpoint                 ResourceType = "endpoints"
	PodSecurityPolicie       ResourceType = "podsecuritypolicies"

	// fuxi
	WorkloadsTemplate ResourceType = "workloads"

	// dynamic page
	FormRender ResourceType = "formrenders"
	Page       ResourceType = "pages"
	Form       ResourceType = "forms"
	Field      ResourceType = "fields"

	// tenant
	BaseDepartment ResourceType = "basedepartments"
	BasePermission ResourceType = "basepermissions"
	BaseRole       ResourceType = "baseroles"
	BaseUser       ResourceType = "baseusers"
	BaseRoleUser   ResourceType = "baseroleusers"

	// network
	IP                          ResourceType = "ips"
	SubNet                      ResourceType = "subnets"
	NetworkAttachmentDefinition ResourceType = "network-attachment-definitions"

	// istio
	Gateway         ResourceType = "gateways"
	DestinationRule ResourceType = "destinationrules"
	ServiceEntry    ResourceType = "serviceentries"
	Sidecar         ResourceType = "sidecars"
	VirtualService  ResourceType = "virtualservices"
	WorkloadEntry   ResourceType = "workloadentries"
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
