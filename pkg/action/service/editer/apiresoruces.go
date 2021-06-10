package editer

import (
	"fmt"
	"github.com/yametech/yamecloud/pkg/action/service"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"strings"
)

var _ service.IResourceService = &APIResources{}

type APIResources struct {
	service.Interface
}

func (g *APIResources) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	panic("implement me")
}

func (g *APIResources) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	panic("implement me")
}

func (g *APIResources) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	panic("implement me")
}

func NewAPIResources(svcInterface service.Interface) *APIResources {
	srv := &APIResources{Interface: svcInterface}
	return srv
}

func (g *APIResources) ListAPIResources() (map[string]*schema.GroupVersionResource, error) {
	_, APIResourceListSlice, err := g.Interface.DiscoveryClient().ServerGroupsAndResources()
	if err != nil {
		return nil, err
	}

	result := make(map[string]*schema.GroupVersionResource)

	for _, APIResourceList := range APIResourceListSlice {
		gv, err := schema.ParseGroupVersion(APIResourceList.GroupVersion)
		if err != nil {
			return nil, err
		}

		for _, APIResource := range APIResourceList.APIResources {
			// filter resource/status examples pod/status
			if strings.Contains(APIResource.Name, "/") {
				continue
			}

			result[fmt.Sprintf("%s/%s/%s", gv.Group, gv.Version, APIResource.Kind)] =
				&schema.GroupVersionResource{
					Group:    gv.Group,
					Version:  gv.Version,
					Resource: APIResource.Name,
				}
		}
	}

	return result, nil

}
