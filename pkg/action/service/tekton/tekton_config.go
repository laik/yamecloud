package tekton

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/k8s"
	"github.com/yametech/yamecloud/pkg/utils"
	"strings"
)

var _ service.IResourceService = &TektonConfig{}

type TektonConfig struct {
	service.Interface
}

func NewTektonConfig(svcInterface service.Interface) *TektonConfig {
	taskRun := &TektonConfig{Interface: svcInterface}
	svcInterface.Install(k8s.TektonConfig, taskRun)
	return taskRun
}

func (g *TektonConfig) Get(namespace, name string) (*service.UnstructuredExtend, error) {
	item, err := g.Interface.Get(namespace, k8s.Secret, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (g *TektonConfig) List(namespace string, selector string) (*service.UnstructuredListExtend, error) {
	list, err := g.Interface.List(namespace, k8s.Secret, selector)
	if err != nil {
		return nil, err
	}
	for index, _ := range list.Items {
		utils.Set(list.Items[index].Object, "metadata.selfLink", strings.Replace(list.Items[index].GetSelfLink(), "secrets", "tektonconfig", -1))
	}
	return list, nil
}

func (g *TektonConfig) Apply(namespace, name string, unstructuredExtend *service.UnstructuredExtend) (*service.UnstructuredExtend, bool, error) {
	item, isUpdate, err := g.Interface.Apply(namespace, k8s.Secret, name, unstructuredExtend)
	if err != nil {
		return nil, isUpdate, err
	}
	return item, isUpdate, nil
}

func (g *TektonConfig) Delete(namespace, name string) error {
	err := g.Interface.Delete(namespace, k8s.Secret, name)
	if err != nil {
		return err
	}
	return nil
}
