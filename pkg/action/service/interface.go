package service

import (
	"github.com/yametech/yamecloud/pkg/k8s"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type ObjectPatcher interface {
	Set(path string, value interface{}) error
	Get(path string) (value interface{}, err error)
}

type UnstructuredListExtend struct {
	*unstructured.UnstructuredList
}

func (u *UnstructuredListExtend) Set(path string, value interface{}) error {
	return nil
}

func (u *UnstructuredListExtend) Get(path string) (interface{}, error) {
	return nil, nil
}

type UnstructuredExtend struct {
	*unstructured.Unstructured
}

// u.Set(".metadata.labels","abc","merge/replace")
func (u *UnstructuredExtend) Set(path string, value interface{}) error {
	return nil
}

// u.Get(".metadata.name")
func (u *UnstructuredExtend) Get(path string) (interface{}, error) {
	return nil, nil
}

type ResourceLister interface {
	Get(namespace, name string) (*UnstructuredExtend, error)
	List(namespace string, selector string) (*UnstructuredListExtend, error)
}

type Interface interface {
	List(namespace string, resource k8s.ResourceType, selector string) (*UnstructuredListExtend, error)
	Get(namespace string, resource k8s.ResourceType, name string) (*UnstructuredExtend, error)
	Install(k8s.ResourceType, ResourceLister)
}

var _ Interface = &Service{}

type Service struct {
	k8s.Interface
	services map[k8s.ResourceType]ResourceLister
}

func (s *Service) Install(resourceType k8s.ResourceType, r ResourceLister) {
	s.services[resourceType] = r
}

func (s *Service) List(namespace string, resource k8s.ResourceType, selector string) (*UnstructuredListExtend, error) {
	list, err := s.Interface.List(namespace, resource, selector)
	if err != nil {
		return nil, err
	}
	return &UnstructuredListExtend{list}, nil
}

func (s *Service) Get(namespace string, resource k8s.ResourceType, name string) (*UnstructuredExtend, error) {
	item, err := s.Interface.Get(namespace, resource, name)
	if err != nil {
		return nil, err
	}
	return &UnstructuredExtend{item}, nil
}

func NewService(k8sInterface k8s.Interface) *Service {
	if k8sInterface == nil {
		panic("datasource interface is nil")
	}
	return &Service{
		Interface: k8sInterface,
		services:  make(map[k8s.ResourceType]ResourceLister),
	}
}

var _ Interface = &FakeService{}

type FakeService struct {
	Data map[string]interface{}
}

func (f *FakeService) Install(resourceType k8s.ResourceType, i ResourceLister) {
	panic("implement me")
}

func (f *FakeService) List(namespace string, resource k8s.ResourceType, selector string) (*UnstructuredListExtend, error) {
	panic("implement me")
}

func (f FakeService) Get(namespace string, resource k8s.ResourceType, name string) (*UnstructuredExtend, error) {
	panic("implement me")
}
func NewFakeService() *FakeService {
	return &FakeService{
		Data: make(map[string]interface{}),
	}
}
