package k8s

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type ResourceType = string

type Resource struct {
	Name   ResourceType
	Schema schema.GroupVersionResource
}

type ResourceRegister interface {
	Register(...Resource)
}

type ResourceLister interface {
	Ranges(d dynamicinformer.DynamicSharedInformerFactory, stop <-chan struct{})
	GroupVersionResource(resourceType ResourceType) (schema.GroupVersionResource, error)
}

type ITypes interface {
	ResourceRegister
	ResourceLister
}

type ICache interface {
	XGet(namespace string, resourceType ResourceType, name string) (*unstructured.Unstructured, error)
}

type Lister interface {
	List(namespace string, resourceType ResourceType, selector string) (*unstructured.UnstructuredList, error)
	Get(namespace string, resourceType ResourceType, name string) (*unstructured.Unstructured, error)
	Cache() ICache
}

type Watcher interface {
	Watch(namespace string, resourceType ResourceType, resourceVersion string, selector string) (<-chan watch.Event, error)
}

type IDataOperator interface {
	Apply(namespace string, resourceType ResourceType, name string, unstructured *unstructured.Unstructured, forceUpdate bool) (newUnstructured *unstructured.Unstructured, isUpdate bool, err error)
	Delete(namespace string, resourceType ResourceType, name string) error
	Patch(namespace string, resourceType ResourceType, name string, data []byte) (*unstructured.Unstructured, error)
}

type RESTClient interface {
	RESTClient() rest.Interface
	ClientSet() *kubernetes.Clientset
}

type Interface interface {
	Lister
	Watcher
	ICache
	IDataOperator
	RESTClient
}
