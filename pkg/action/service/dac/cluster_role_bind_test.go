package dac

import (
	"fmt"
	"reflect"
	"testing"
)

func TestClusterRoleBinding_List(t *testing.T) {
	service := NewFakeService()
	crb := &ClusterRoleBinding{service}
	ule, err := crb.List("", "")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(ule)
	expected := ""
	if !reflect.DeepEqual(ule, expected) {
		t.Fatal("compare object not equal")
	}
}
