package utils

import (
	"encoding/json"
	gyaml "github.com/ghodss/yaml"
	"github.com/tidwall/gjson"
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
	_unstructured, err := runtime.DefaultUnstructuredConverter.ToUnstructured(object)
	if err != nil {
		return nil, err
	}
	return &unstructured.Unstructured{
		Object: _unstructured,
	}, nil
}

func UnmarshalBytesToUnstructured(data []byte, unstructured *unstructured.Unstructured) error {
	return json.Unmarshal(data, unstructured)
}

func UnmarshalUnstructuredList(obj *unstructured.UnstructuredList, target interface{}) error {
	bytesData, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytesData, target)
}

func UnmarshalUnstructured(unstructured *unstructured.Unstructured, target interface{}) error {
	bytesData, err := json.Marshal(unstructured)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytesData, target)
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
	if sort.SearchStrings(list, item) >= 0 {
		return true
	}
	return false
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

	dest, err := delete(string(bytes), []string{
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

func SetObjectOwner(object []byte, apiVersion, kind, name, uid string) (*unstructured.Unstructured, error) {
	type ownerReference struct {
		ApiVersion         string `json:"apiVersion"`
		Kind               string `json:"kind"`
		Name               string `json:"name"`
		UID                string `json:"uid"`
		Controller         bool   `json:"controller"`
		BlockOwnerDeletion bool   `json:"blockOwnerDeletion"`
	}

	s, err := sjson.Set(string(object), "metadata.ownerReferences", []ownerReference{{
		ApiVersion:         apiVersion,
		Kind:               kind,
		Name:               name,
		UID:                uid,
		Controller:         false,
		BlockOwnerDeletion: false,
	}})
	if err != nil {
		return nil, err
	}

	obj := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &obj); err != nil {
		return nil, err
	}

	return &unstructured.Unstructured{Object: obj}, nil
}

func SetYamlValue(yamlData []byte, path string, value interface{}) ([]byte, error) {
	jsonData, err := gyaml.YAMLToJSON(yamlData)
	if err != nil {
		return []byte(""), err
	}
	j, err := sjson.Set(string(jsonData), path, value)
	if err != nil {
		return []byte(""), err
	}
	return gyaml.JSONToYAML([]byte(j))
}

func GetYamlValue(yamlData []byte, path string) (gjson.Result, error) {
	jsonData, err := gyaml.YAMLToJSON(yamlData)
	if err != nil {
		return gjson.Result{}, err
	}
	return gjson.Get(string(jsonData), path), nil
}

type Result = gjson.Result

func GetObjectValue(object interface{}, path string) (*Result, error) {
	bytesValue, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}
	result := gjson.Get(string(bytesValue), path)
	return &result, nil
}
