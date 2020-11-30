package k8s

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/watch"
)

type ICache interface {
	XGet(namespace string, resourceType ResourceType, name string) (*unstructured.Unstructured, error)
}

type Lister interface {
	List(namespace string, resourceType ResourceType, selector interface{}) (*unstructured.UnstructuredList, error)
	Get(namespace string, resourceType ResourceType, name string) (*unstructured.Unstructured, error)
	Cache() ICache
}

type Watcher interface {
	Watch(namespace string, resourceType ResourceType, resourceVersion string, selector interface{}) (<-chan watch.Event, error)
}

type IDataOperator interface {
	Apply(namespace string, resourceType ResourceType, name string, unstructured *unstructured.Unstructured, forceUpdate bool) (newUnstructured *unstructured.Unstructured, isUpdate bool, err error)
	Delete(namespace string, resourceType ResourceType, name string) error
	Patch(namespace string, resourceType ResourceType, name string, path string, data interface{}) error
}

type Interface interface {
	Lister
	Watcher
	ICache
	IDataOperator
}
