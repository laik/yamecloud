package tekton

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yametech/yamecloud/pkg/action/api/common"
	"github.com/yametech/yamecloud/pkg/action/service"
	v1 "github.com/yametech/yamecloud/pkg/apis/yamecloud/v1"
	"github.com/yametech/yamecloud/pkg/utils"
	"net/http"
)

const (
	TektonGraphKind       = "TektonGraph"
	TektonGraphApiVersion = "fuxi.nip.io/v1"
	DefaultNodeId         = "1-1"
	DefaultTaskName       = "node-1-1"
)

// Get TektonGraph
func (s *tektonServer) GetTektonGraph(g *gin.Context) {
	namespace := g.Param("namespace")
	name := g.Param("name")
	if namespace == "" || name == "" {
		common.RequestParametersError(g, fmt.Errorf("params not obtain namespace=%s name=%s", namespace, name))
		return
	}
	item, err := s.TektonGraph.Get(namespace, name)
	if err != nil {
		common.InternalServerError(g, err, err)
		return
	}
	g.JSON(http.StatusOK, item)
}

// Subscribe TektonGraph
func (s *tektonServer) ListTektonGraph(g *gin.Context) {
	list, err := s.TektonGraph.List(g.Param("namespace"), "")
	if err != nil {
		common.InternalServerError(g, "", err)
		return
	}
	g.JSON(http.StatusOK, list)
}

// Update or Create TektonGraph
func (s *tektonServer) ApplyTektonGraph(g *gin.Context) {
	namespace := g.Param("namespace")
	raw, err := g.GetRawData()
	if err != nil {
		common.RequestParametersError(g, fmt.Errorf("get raw data error (%s)", err))
		return
	}

	obj := &v1.TektonGraph{}
	err = json.Unmarshal(raw, obj)
	if err != nil {
		common.RequestParametersError(g, fmt.Errorf("get raw data error (%s)", err))
		return
	}
	_unstructured, err := utils.ObjectToUnstructured(obj)
	if err := _unstructured.UnmarshalJSON(raw); err != nil {
		common.RequestParametersError(g, fmt.Errorf("unmarshal from json data error (%s)", err))
		return
	}
	name := _unstructured.GetName()
	newUnstructuredExtend, isUpdate, err := s.TektonGraph.Apply(namespace, name, &service.UnstructuredExtend{Unstructured: _unstructured})
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

// Delete TektonGraph
func (s *tektonServer) DeleteTektonGraph(g *gin.Context) {
	namespace := g.Param("namespace")
	name := g.Param("name")
	if namespace == "" || name == "" {
		common.RequestParametersError(g, fmt.Errorf("params not obtain namespace=%s name=%s", namespace, name))
		return
	}
	err := s.TektonGraph.Delete(namespace, name)
	if err != nil {
		common.InternalServerError(g, err, err)
		return
	}
	g.JSON(http.StatusOK, nil)
}

func GetDefaultInitData() map[string]interface{} {
	return map[string]interface{}{
		"nodes": []map[string]interface{}{
			{
				"id":       DefaultNodeId,
				"x":        20,
				"y":        20,
				"role":     0,
				"taskName": DefaultTaskName,
				"anchorPoints": [2][2]float32{
					{0, 0.5},
					{1, 0.5},
				},
			},
		},
	}
}
