package globalconfig

import (
	"fmt"
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
)

var _ service.IResourceService = &GlobalConfig{}

type GlobalConfig struct {
	service.Interface
}

func NewGlobalConfig(svcInterface service.Interface) *GlobalConfig {
	serviceAccount := &GlobalConfig{Interface: svcInterface}
	svcInterface.Install(k8s.GlobalConfig, serviceAccount)
	return serviceAccount
}

func (c *GlobalConfig) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := c.Interface.Get("", k8s.GlobalConfig, "compass-config")
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (c *GlobalConfig) GetHelmRepos() (map[string]string, error) {
	item, err := c.Interface.Get("", k8s.GlobalConfig, "compass-config")
	if err != nil {
		return nil, err
	}
	helmReposInterface, _ := item.Get("spec.helmRepos")

	helmRepos, ok := helmReposInterface.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("can't not convert helmRepos")
	}

	result := make(map[string]string)
	for k, v := range helmRepos {
		result[k] = v.(string)
	}

	return result, nil
}

func (c *GlobalConfig) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	panic("implement me")
}

func (c *GlobalConfig) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	panic("implement me")
}
