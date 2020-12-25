package dac

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"reflect"
	"testing"
)

func Test_ClusterRoleList(t *testing.T) {
	expected := &service.UnstructuredListExtend{}
	fakeService := service.NewFakeService()
	cr := &ClusterRole{fakeService}
	ule, err := cr.List("", "apps=123")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(ule, expected) {
		t.Fatal("compare object not equal")
	}
}
