package service

import (
	"fmt"
	"github.com/yametech/yamecloud/pkg/k8s"
	"github.com/yametech/yamecloud/pkg/utils"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"strconv"
	"strings"
)

var _ ObjectPatcher = &UnstructuredExtend{}
var _ ObjectPatcher = &UnstructuredListExtend{}

type ObjectPatcher interface {
	Set(path string, value interface{}) error
	Get(path string) (value interface{}, err error)
}

type UnstructuredListExtend struct {
	*unstructured.UnstructuredList
}

// ul.Set("0.metadata.labels","abc")
func (ul *UnstructuredListExtend) Set(path string, value interface{}) error {
	pathList := strings.Split(path, ".")
	if len(pathList) < 1 {
		return nil
	}
	index, err := strconv.ParseUint(pathList[0], 10, 64)
	if err != nil {
		return fmt.Errorf("not found index because parse index error %s", err)
	}
	if uint64(len(ul.Items)-1) > index {
		return fmt.Errorf("not found index %d item", index)
	}

	utils.Set(ul.Items[index].Object, strings.TrimPrefix(path, fmt.Sprintf("%d.", index)), value)

	return nil
}

// ul.Get("0.metadata.labels","abc")
func (ul *UnstructuredListExtend) Get(path string) (interface{}, error) {
	pathList := strings.Split(path, ".")
	if len(pathList) < 1 {
		return nil, fmt.Errorf("not found index")
	}
	index, err := strconv.ParseUint(pathList[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("not found index because parse index error %s", err)
	}
	if uint64(len(ul.Items)-1) > index {
		return nil, fmt.Errorf("not found index %d item", index)
	}

	return utils.Get(ul.Items[index].Object, strings.TrimPrefix(path, fmt.Sprintf("%d.", index))), nil
}

type UnstructuredExtend struct {
	*unstructured.Unstructured
}

// u.Set("metadata.labels","abc","merge/replace")
func (u *UnstructuredExtend) Set(path string, value interface{}) error {
	utils.Set(u.Object, path, value)
	return nil
}

// u.Get("metadata.name")
func (u *UnstructuredExtend) Get(path string) (interface{}, error) {
	return utils.Get(u.Object, path), nil
}

type IResourceService interface {
	Get(namespace, name string) (*UnstructuredExtend, error)
	List(namespace string, selector string) (*UnstructuredListExtend, error)
	Apply(namespace, name string, unstructuredExtend *UnstructuredExtend) (*UnstructuredExtend, bool, error)
}

type Interface interface {
	List(namespace string, resource k8s.ResourceType, selector string) (*UnstructuredListExtend, error)
	Get(namespace string, resource k8s.ResourceType, name string) (*UnstructuredExtend, error)
	Apply(namespace string, resource k8s.ResourceType, name string, unstructuredExtend *UnstructuredExtend) (*UnstructuredExtend, bool, error)
	ForceUpdate(namespace string, resource k8s.ResourceType, name string, unstructuredExtend *UnstructuredExtend) (*UnstructuredExtend, error)
	Delete(namespace string, resource k8s.ResourceType, name string) error
	Install(k8s.ResourceType, IResourceService)
}

var _ Interface = &Service{}

type Service struct {
	k8s.Interface
	services map[k8s.ResourceType]IResourceService
}

func (s *Service) Install(resourceType k8s.ResourceType, r IResourceService) {
	s.services[resourceType] = r
}

func (s *Service) Apply(namespace string, resource k8s.ResourceType, name string, unstructuredExtend *UnstructuredExtend) (*UnstructuredExtend, bool, error) {
	u, isUpdate, err := s.Interface.Apply(namespace, resource, name, unstructuredExtend.Unstructured, false)
	return &UnstructuredExtend{Unstructured: u}, isUpdate, err
}

func (s *Service) ForceUpdate(namespace string, resource k8s.ResourceType, name string, unstructuredExtend *UnstructuredExtend) (*UnstructuredExtend, error) {
	u, _, err := s.Interface.Apply(namespace, resource, name, unstructuredExtend.Unstructured, true)
	return &UnstructuredExtend{Unstructured: u}, err
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
		services:  make(map[k8s.ResourceType]IResourceService),
	}
}

var _ Interface = &FakeService{}

type FakeService struct {
	Data map[string]interface{}
}

func (f *FakeService) ForceUpdate(namespace string, resource k8s.ResourceType, name string, unstructuredExtend *UnstructuredExtend) (*UnstructuredExtend, error) {
	panic("implement me")
}

func (f *FakeService) Delete(namespace string, resource k8s.ResourceType, name string) error {
	panic("implement me")
}

func (f *FakeService) Apply(namespace string, resource k8s.ResourceType, name string, unstructuredExtend *UnstructuredExtend) (*UnstructuredExtend, bool, error) {
	panic("implement me")
}

func (f *FakeService) Install(resourceType k8s.ResourceType, i IResourceService) {
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
