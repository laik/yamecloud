package utils

import (
	"encoding/json"
	"github.com/tidwall/sjson"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"reflect"
	"sort"
)

func UnmarshalObject(object runtime.Object, target interface{}) error {
	bytesData, err := json.Marshal(object)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytesData, target)
}

func ObjectToUnstructured(object runtime.Object) (*unstructured.Unstructured, error) {
	obj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(object)
	if err != nil {
		return nil, err
	}
	return &unstructured.Unstructured{
		Object: obj,
	}, nil
}

func CompareSpecByUnstructured(source, target *unstructured.Unstructured) bool {
	if source == nil || target == nil {
		return false
	}
	srcUnstructuredSpec, exist := source.Object["spec"]
	if !exist {
		return false
	}
	targetUnstructuredSpec, exist := target.Object["spec"]
	if !exist {
		return false
	}
	if !reflect.DeepEqual(srcUnstructuredSpec, targetUnstructuredSpec) {
		return false
	}
	return true
}

func ContainStringItem(list []string, item string) bool {
	if sort.SearchStrings(list, item) == len(list) {
		return false
	}
	return true
}

func CloneNewObject(src *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	bytes, err := src.MarshalJSON()
	if err != nil {
		return nil, err
	}

	delete := func(res string, paths []string) (string, error) {
		var err error
		for _, path := range paths {
			res, err = sjson.Delete(res, path)
			if err != nil {
				return "", err
			}
		}
		return res, nil
	}

	dest, err := delete(string(bytes),
		[]string{
			"metadata.creationTimestamp",
			"metadata.generation",
			"metadata.managedFields",
			"metadata.resourceVersion",
			"metadata.selfLink",
			"metadata.uid",
			"status",
		})

	obj := make(map[string]interface{})
	if err := json.Unmarshal([]byte(dest), &obj); err != nil {
		return nil, err
	}

	return &unstructured.Unstructured{Object: obj}, nil
}
