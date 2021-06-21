package dac

import (
	"fmt"
	"reflect"
	"testing"
)

func TestServiceAccount_List(t *testing.T) {
	service := NewFakeService()
	sa := &ServiceAccount{service}
	ule, err := sa.List("xxx", "")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(ule)
	expected := ""
	if !reflect.DeepEqual(ule, expected) {
		t.Fatal("compare object not equal")
	}
}
