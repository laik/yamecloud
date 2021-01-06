package k8s

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sort"
)

const (
	Water                    ResourceType = "water"
	Deployment               ResourceType = "deployment"
	Stone                    ResourceType = "stone"
	StatefulSet              ResourceType = "statefulset"
	StatefulSet1             ResourceType = "statefulset1"
	DaemonSet                ResourceType = "daemonsets"
	Injector                 ResourceType = "injector"
	Pod                      ResourceType = "pod"
	Job                      ResourceType = "jobs"
	CronJobs                 ResourceType = "cronjobs"
	ReplicaSet               ResourceType = "replicasets"
	Event                    ResourceType = "events"
	Node                     ResourceType = "nodes"
	ConfigMaps               ResourceType = "configmaps"
	Secrets                  ResourceType = "secrets"
	ResourceQuota            ResourceType = "resourcequotas"
	Service                  ResourceType = "micro"
	Ingress                  ResourceType = "ingresses"
	NetworkPolicy            ResourceType = "networkpolicies"
	HorizontalPodAutoscaler  ResourceType = "horizontalpodautoscalers"
	CustomResourceDefinition ResourceType = "customresourcedefinitions"
	PersistentVolume         ResourceType = "persistentvolumes"
	PersistentVolumeClaims   ResourceType = "persistentvolumeclaims"
	StorageClass             ResourceType = "storageclasses"
	ServiceAccount           ResourceType = "serviceaccounts"
	Role                     ResourceType = "roles"
	ClusterRole              ResourceType = "clusterroles"
	RoleBinding              ResourceType = "rolebindings"
	Namespace                ResourceType = "namesapces"
	PodSecurityPolicie       ResourceType = "podsecuritypolicies"
	ClusterRoleBinding       ResourceType = "clusterrolebindings"
	Endpoint                 ResourceType = "endpoints"

	// deployment resource workload template for CaaS
	WorkloadsTemplate ResourceType = "workloads"

	// form render
	FormRender ResourceType = "formrenders"
	Page       ResourceType = "pages"
	Form       ResourceType = "forms"
	Field      ResourceType = "fields"

	// tenant for PaaS
	BaseDepartment ResourceType = "basedepartments"
	BaseTenant     ResourceType = "basedtenants"
	BaseRole       ResourceType = "baseroles"
	BaseUser       ResourceType = "baseusers"
	BaseRoleUser   ResourceType = "baseroleusers"

	// network for container ovn/ovs control plant
	IP                          ResourceType = "ips"
	SubNet                      ResourceType = "subnets"
	Vlan                        ResourceType = "vlans"
	NetworkAttachmentDefinition ResourceType = "network-attachment-definitions"

	//tekton
	Pipeline         ResourceType = "pipelines"
	PipelineRun      ResourceType = "pipelineruns"
	Task             ResourceType = "tasks"
	TaskRun          ResourceType = "taskruns"
	PipelineResource ResourceType = "pipelineresources"
	TektonGraph      ResourceType = "tektongraphs"
	TektonWebHook    ResourceType = "tektonwebhooks"
	TektonStore      ResourceType = "tektonstores"

	//Istio  NetWorking
	Gateway         ResourceType = "gateways"
	DestinationRule ResourceType = "destinationrules"
	ServiceEntry    ResourceType = "serviceentries"
	Sidecar         ResourceType = "sidecars"
	VirtualService  ResourceType = "virtualservices"
	WorkloadEntry   ResourceType = "workloadentries"
)

type IGVRMaps interface {
	Subscribe(include ...string) []Resource
}

type groupVersionCollection map[ResourceType]schema.GroupVersionResource

func (c *groupVersionCollection) Subscribe(include ...string) []Resource {
	result := make([]Resource, 0)
	for key, value := range *c {
		if len(include) == 0 {
			continue
		}
		if sort.SearchStrings(include, key) > len(include)-1 {
			continue
		}
		result = append(result, Resource{key, value})
	}
	return result
}

// describe resource collection
var GVRMaps IGVRMaps = &groupVersionCollection{
	Water:        {Group: "nuwa.nip.io", Version: "v1", Resource: "waters"},
	Stone:        {Group: "nuwa.nip.io", Version: "v1", Resource: "stones"},
	StatefulSet1: {Group: "nuwa.nip.io", Version: "v1", Resource: "statefulsets"},
	Injector:     {Group: "nuwa.nip.io", Version: "v1", Resource: "injectors"},

	WorkloadsTemplate: {Group: "fuxi.nip.io", Version: "v1", Resource: "workloads"},

	Page:  {Group: "fuxi.nip.io", Version: "v1", Resource: "pages"},
	Form:  {Group: "fuxi.nip.io", Version: "v1", Resource: "forms"},
	Field: {Group: "fuxi.nip.io", Version: "v1", Resource: "fields"},

	// PaaS RBAC
	BaseDepartment: {Group: "fuxi.nip.io", Version: "v1", Resource: "basedepartments"},
	BaseRole:       {Group: "fuxi.nip.io", Version: "v1", Resource: "baseroles"},
	BaseUser:       {Group: "fuxi.nip.io", Version: "v1", Resource: "baseusers"},
	BaseRoleUser:   {Group: "fuxi.nip.io", Version: "v1", Resource: "baseroleusers"},

	StatefulSet: {Group: "apps", Version: "v1", Resource: "statefulsets"},
	DaemonSet:   {Group: "apps", Version: "v1", Resource: "daemonsets"},
	ReplicaSet:  {Group: "apps", Version: "v1", Resource: "replicasets"},
	Deployment:  {Group: "apps", Version: "v1", Resource: "deployments"},

	Job:      {Group: "batch", Version: "v1", Resource: "jobs"},
	CronJobs: {Group: "batch", Version: "v1beta1", Resource: "cronjobs"},

	Pod:                    {Group: "", Version: "v1", Resource: "pods"},
	Node:                   {Group: "", Version: "v1", Resource: "nodes"},
	Event:                  {Group: "", Version: "v1", Resource: "events"},
	ConfigMaps:             {Group: "", Version: "v1", Resource: "configmaps"},
	Secrets:                {Group: "", Version: "v1", Resource: "secrets"},
	ResourceQuota:          {Group: "", Version: "v1", Resource: "resourcequotas"},
	Service:                {Group: "", Version: "v1", Resource: "micro"},
	Namespace:              {Group: "", Version: "v1", Resource: "namespaces"},
	PersistentVolume:       {Group: "", Version: "v1", Resource: "persistentvolumes"},
	PersistentVolumeClaims: {Group: "", Version: "v1", Resource: "persistentvolumeclaims"},
	ServiceAccount:         {Group: "", Version: "v1", Resource: "serviceaccounts"},
	Endpoint:               {Group: "", Version: "v1", Resource: "endpoints"},

	Ingress:                  {Group: "extensions", Version: "v1beta1", Resource: "ingresses"},
	NetworkPolicy:            {Group: "networking.k8s.io", Version: "v1", Resource: "networkpolicies"},
	HorizontalPodAutoscaler:  {Group: "autoscaling", Version: "v2beta1", Resource: "horizontalpodautoscalers"},
	CustomResourceDefinition: {Group: "apiextensions.k8s.io", Version: "v1beta1", Resource: "customresourcedefinitions"},

	StorageClass: {Group: "storage.k8s.io", Version: "v1", Resource: "storageclasses"},

	ClusterRole:        {Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "clusterroles"},
	Role:               {Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "roles"},
	ClusterRoleBinding: {Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "clusterrolebindings"},
	RoleBinding:        {Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "rolebindings"},

	IP:     {Group: "kubeovn.io", Version: "v1", Resource: "ips"},
	SubNet: {Group: "kubeovn.io", Version: "v1", Resource: "subnets"},
	Vlan:   {Group: "kubeovn.io", Version: "v1", Resource: "vlans"},

	NetworkAttachmentDefinition: {Group: "k8s.cni.cncf.io", Version: "v1", Resource: "network-attachment-definitions"},

	// tekton.dev resource view
	Pipeline:         {Group: "tekton.dev", Version: "v1alpha1", Resource: "pipelines"},
	PipelineRun:      {Group: "tekton.dev", Version: "v1alpha1", Resource: "pipelineruns"},
	Task:             {Group: "tekton.dev", Version: "v1alpha1", Resource: "tasks"},
	TaskRun:          {Group: "tekton.dev", Version: "v1alpha1", Resource: "taskruns"},
	PipelineResource: {Group: "tekton.dev", Version: "v1alpha1", Resource: "pipelineresources"},
	TektonGraph:      {Group: "fuxi.nip.io", Version: "v1", Resource: "tektongraphs"},
	TektonWebHook:    {Group: "fuxi.nip.io", Version: "v1", Resource: "tektonwebhooks"},
	TektonStore:      {Group: "fuxi.nip.io", Version: "v1", Resource: "tektonstores"},

	PodSecurityPolicie: {Group: "policy", Version: "v1beta1", Resource: "podsecuritypolicies"},

	//Istio Networking
	Gateway:         {Group: "networking.istio.io", Version: "v1beta1", Resource: "gateways"},
	DestinationRule: {Group: "networking.istio.io", Version: "v1beta1", Resource: "destinationrules"},
	ServiceEntry:    {Group: "networking.istio.io", Version: "v1beta1", Resource: "serviceentries"},
	Sidecar:         {Group: "networking.istio.io", Version: "v1beta1", Resource: "sidecars"},
	VirtualService:  {Group: "networking.istio.io", Version: "v1beta1", Resource: "virtualservices"},
	WorkloadEntry:   {Group: "networking.istio.io", Version: "v1beta1", Resource: "workloadentries"},
}
