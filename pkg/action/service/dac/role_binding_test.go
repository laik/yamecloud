package dac

import (
	"reflect"
	"testing"
)

func TestRoleBinding_List(t *testing.T) {
	service := NewFakeService()
	rb := &RoleBinding{service}
	ule, err := rb.List("", "")
	if err != nil {
		t.Fatal(err)
	}

	expected := ""
	if !reflect.DeepEqual(ule, expected) {
		t.Fatal("compare object not equal")
	}
}
