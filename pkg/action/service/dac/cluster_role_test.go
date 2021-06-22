package dac

import (
	"fmt"
	"github.com/yametech/yamecloud/pkg/action/api"
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/configure"
	"github.com/yametech/yamecloud/pkg/k8s"
	"github.com/yametech/yamecloud/pkg/k8s/datasource"
	"github.com/yametech/yamecloud/pkg/k8s/types"
	"reflect"
	"testing"
)

var subscribeList = k8s.GVRMaps.Subscribe(
	k8s.ClusterRole,
	k8s.ServiceAccount,
	k8s.RoleBinding,
	k8s.ClusterRoleBinding,
)

func NewFakeService() *api.Server {
	config, err := configure.NewInstallConfigure(types.NewResourceITypes(subscribeList))
	if err != nil {
		panic(fmt.Sprintf("new install configure error %s", err))
	}

	_datasource := datasource.NewInterface(config)
	return api.NewServer(service.NewService(_datasource))
}

func Test_ClusterRoleList(t *testing.T) {
	expected := &service.UnstructuredListExtend{}
	fakeService := NewFakeService()
	cr := &ClusterRole{fakeService}
	ule, err := cr.List("", "apps=123")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(ule, expected) {
		t.Fatal("compare object not equal")
	}
}
