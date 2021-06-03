package workloads

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yametech/yamecloud/pkg/action/api/common"
	"github.com/yametech/yamecloud/pkg/action/service"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"net/http"
)

func (s *workloadServer) GetNamespace(g *gin.Context) {
	namespace := g.Param("namespace")
	if namespace == "" {
		common.RequestParametersError(g, fmt.Errorf("params not obtain namespace name=%s", namespace))
		return
	}
	item, err := s.Namespace.Get("", namespace)
	if err != nil {
		common.InternalServerError(g, err, err)
		return
	}
	g.JSON(http.StatusOK, item)
}

func (s *workloadServer) ListNamespace(g *gin.Context) {
	list, err := s.Namespace.List("", "")
	if err != nil {
		common.InternalServerError(g, "", err)
		return
	}
	g.JSON(http.StatusOK, list)
}

func (s *workloadServer) ApplyNamespace(g *gin.Context) {
	raw, err := g.GetRawData()
	if err != nil {
		common.RequestParametersError(g, fmt.Errorf("get raw data error (%s)", err))
		return
	}

	_unstructured := &unstructured.Unstructured{}
	if err := _unstructured.UnmarshalJSON(raw); err != nil {
		common.RequestParametersError(g, fmt.Errorf("unmarshal from json data error (%s)", err))
		return
	}
	name := _unstructured.GetName()
	newUnstructuredExtend, isUpdate, err := s.Namespace.Apply(_unstructured.GetName(), name, &service.UnstructuredExtend{Unstructured: _unstructured})
	if err != nil {
		common.InternalServerError(g, newUnstructuredExtend, fmt.Errorf("apply object error (%s)", err))
		return
	}

	if isUpdate {
		g.JSON(
			http.StatusOK,
			[]service.UnstructuredExtend{
				*newUnstructuredExtend,
			})
	} else {
		g.JSON(http.StatusOK, newUnstructuredExtend)
	}
}

func (s *workloadServer) UpdateNamespace(g *gin.Context) {
	namespace := g.Param("namespace")
	name := g.Param("name")
	if namespace == "" || name == "" {
		common.RequestParametersError(g, fmt.Errorf("params not obtain namespace=%s name=%s", namespace, name))
		return
	}
	raw, err := g.GetRawData()
	if err != nil {
		common.RequestParametersError(g, fmt.Errorf("get raw data error (%s)", err))
		return
	}

	updateNetWorkAttachmentData := &unstructured.Unstructured{}
	if err := json.Unmarshal(raw, updateNetWorkAttachmentData); err != nil {
		common.RequestParametersError(g, fmt.Errorf("unmarshal from json data error (%s)", err))
		return
	}

	newUnstructuredExtend, _, err := s.Namespace.Apply(namespace, name, &service.UnstructuredExtend{Unstructured: updateNetWorkAttachmentData})
	if err != nil {
		common.InternalServerError(g, err, err)
		return
	}
	g.JSON(
		http.StatusOK,
		[]service.UnstructuredExtend{
			*newUnstructuredExtend,
		})
}

func (s *workloadServer) DeleteNamespace(g *gin.Context) {
	name := g.Param("namespaces")
	if name == "" {
		common.RequestParametersError(g, fmt.Errorf("params not obtain name=%s", name))
		return
	}
	err := s.Namespace.Delete("", name)
	if err != nil {
		common.InternalServerError(g, err, err)
		return
	}
	g.JSON(http.StatusOK, nil)
}

type annotate struct {
	Data annotateData `json:"data"`
}
type annotateData struct {
	Namespace         string   `json:"namespace"`
	Nodes             []string `json:"nodes"`
	StorageClasses    []string `json:"storageClasses"`
	NetworkAttachment string   `json:"networkAttachment"`
}

func (s *workloadServer) AnnotateNamespaceAllowedNode(g *gin.Context) {
	rawData, err := g.GetRawData()
	if err != nil {
		common.RequestParametersError(g, fmt.Errorf("params not obtain form data error: %s", err))
		return
	}

	ad := annotate{}
	err = json.Unmarshal(rawData, &ad)
	if err != nil {
		common.RequestParametersError(g, fmt.Errorf("params not obtain form data (%s) error: %s", rawData, err))
		return
	}

	// check namespace exist
	namespace, err := s.Namespace.Get("", ad.Data.Namespace)
	if err != nil {
		common.RequestParametersError(g, fmt.Errorf("params not obtain form data (%s) check namespace error: %s", rawData, err))
		return
	}

	currentAnnotations := namespace.GetAnnotations()

	coordinates := make([]map[string]string, 0)
	for _, nodeName := range ad.Data.Nodes {
		node, err := s.Node.Get("", nodeName)
		if err != nil {
			common.RequestParametersError(g, fmt.Errorf("node not found: %s", nodeName))
			return
		}
		coordinate := make(map[string]string)
		nodeLabels := node.GetLabels()
		if value, exist := nodeLabels["nuwa.kubernetes.io/zone"]; !exist {
			common.RequestParametersError(g, fmt.Errorf(`node %s not flag nuwa.kubernetes.io/zone`, nodeName))
			return
		} else {
			coordinate["zone"] = value
		}

		if value, exist := nodeLabels["nuwa.kubernetes.io/rack"]; !exist {
			common.RequestParametersError(g, fmt.Errorf(`node %s not flag nuwa.kubernetes.io/rack`, nodeName))
			return
		} else {
			coordinate["rack"] = value
		}

		if value, exist := nodeLabels["nuwa.kubernetes.io/host"]; !exist {
			common.RequestParametersError(g, fmt.Errorf(`node %s not flag nuwa.kubernetes.io/host`, nodeName))
			return
		} else {
			coordinate["host"] = value
		}

		coordinate["name"] = nodeName
		coordinates = append(coordinates, coordinate)
	}

	if len(coordinates) > 0 {
		bs, _ := json.Marshal(distinctCoordinateListMap(coordinates))
		currentAnnotations["nuwa.kubernetes.io/default_resource_limit"] = string(bs)
	} else {
		delete(currentAnnotations, "nuwa.kubernetes.io/default_resource_limit")
	}

	namespace.SetAnnotations(currentAnnotations)

	newUnstructuredExtend, _, err := s.Namespace.Apply("", namespace.GetName(), namespace)
	if err != nil {
		common.InternalServerError(g, err, err)
		return
	}
	g.JSON(http.StatusOK, newUnstructuredExtend)
}

func distinctCoordinateListMap(list []map[string]string) []map[string]string {
	d := make(map[string]map[string]string)
	for _, item := range list {
		key := fmt.Sprintf("%s-%s-%s", item["zone"], item["rack"], item["host"])
		d[key] = map[string]string{
			"zone": item["zone"],
			"rack": item["rack"],
			"host": item["host"],
			"name": item["name"],
		}
	}

	result := make([]map[string]string, 0)
	for _, v := range d {
		result = append(result, v)
	}

	return result
}

func (s *workloadServer) AnnotateNamespaceNetworkAttach(g *gin.Context) {
	rawData, err := g.GetRawData()
	if err != nil {
		common.RequestParametersError(g, fmt.Errorf("params not obtain form data error: %s", err))
		return
	}

	ad := annotate{}
	err = json.Unmarshal(rawData, &ad)
	if err != nil {
		common.RequestParametersError(g, fmt.Errorf("params not obtain form data (%s) error: %s", rawData, err))
		return
	}

	// check namespace exist
	namespace, err := s.Namespace.Get("", ad.Data.Namespace)
	if err != nil {
		common.RequestParametersError(g, fmt.Errorf("params not obtain form data (%s) check namespace error: %s", rawData, err))
		return
	}

	currentAnnotations := namespace.GetAnnotations()

	if ad.Data.NetworkAttachment != "" {
		currentAnnotations["k8s.v1.cni.cncf.io/namespaces"] = ad.Data.NetworkAttachment

	} else {
		delete(currentAnnotations, "k8s.v1.cni.cncf.io/namespaces")
	}

	namespace.SetAnnotations(currentAnnotations)

	newUnstructuredObj, _, err := s.Namespace.Apply("", ad.Data.Namespace, namespace)
	if err != nil {
		common.InternalServerError(g, err, err)
		return
	}

	g.JSON(http.StatusOK, newUnstructuredObj)
}

func (s *workloadServer) AnnotateNamespaceAllowedStorageClass(g *gin.Context) {
	rawData, err := g.GetRawData()
	if err != nil {
		common.RequestParametersError(g, fmt.Errorf("params not obtain form data error: %s", err))
		return
	}

	ad := annotate{}
	err = json.Unmarshal(rawData, &ad)
	if err != nil {
		common.RequestParametersError(g, fmt.Errorf("params not obtain form data (%s) error: %s", rawData, err))
		return
	}

	// check namespace exist
	namespace, err := s.Namespace.Get("", ad.Data.Namespace)
	if err != nil {
		common.RequestParametersError(g, fmt.Errorf("params not obtain form data (%s) check namespace error: %s", rawData, err))
		return
	}
	currentAnnotations := namespace.GetAnnotations()
	ad.Data.StorageClasses = clearEmptyCharactersItem(ad.Data.StorageClasses)

	if len(ad.Data.StorageClasses) > 0 {
		bs, _ := json.Marshal(ad.Data.StorageClasses)
		currentAnnotations["fuxi.kubernetes.io/default_storage_limit"] = string(bs)
	} else {
		delete(currentAnnotations, "fuxi.kubernetes.io/default_storage_limit")
	}
	namespace.SetAnnotations(currentAnnotations)

	newUnstructuredObj, _, err := s.Namespace.Apply("", ad.Data.Namespace, namespace)
	if err != nil {
		common.InternalServerError(g, err, err)
		return
	}

	g.JSON(http.StatusOK, newUnstructuredObj)
}

func clearEmptyCharactersItem(slice []string) []string {
	result := make([]string, 0)
	for _, item := range slice {
		if len(item) == 0 || item == "" {
			continue
		}
		result = append(result, item)
	}
	return result
}
