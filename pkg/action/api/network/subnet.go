package network

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yametech/yamecloud/pkg/action/api/common"
	"github.com/yametech/yamecloud/pkg/action/service"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"net/http"
)

// GetSubNet none
func (s *networkServer) GetSubNet(g *gin.Context) {
	//namespace := g.Param("namespace")
	name := g.Param("name")
	if name == "" {
		common.RequestParametersError(g, fmt.Errorf("params not obtain namespace=%s name=%s", "", name))
		return
	}
	item, err := s.SubNet.Get("", name)
	if err != nil {
		common.InternalServerError(g, err, err)
		return
	}
	g.JSON(http.StatusOK, item)
}

// ListSubNet list all subnet
func (s *networkServer) ListSubNet(g *gin.Context) {
	list, err := s.SubNet.List("", "")
	if err != nil {
		common.InternalServerError(g, "", err)
		return
	}
	g.JSON(http.StatusOK, list)
}

// ApplySubNet Update or Create SubNet
func (s *networkServer) ApplySubNet(g *gin.Context) {
	namespace := ""
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
	newUnstructuredExtend, isUpdate, err := s.SubNet.Apply(namespace, name, &service.UnstructuredExtend{Unstructured: _unstructured})
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

// DeleteSubNet Delete Service
func (s *networkServer) DeleteSubNet(g *gin.Context) {
	namespace := ""
	name := g.Param("name")
	if name == "" {
		common.RequestParametersError(g, fmt.Errorf("params not obtain namespace=%s name=%s", namespace, name))
		return
	}
	err := s.SubNet.Delete(namespace, name)
	if err != nil {
		common.InternalServerError(g, err, err)
		return
	}
	g.JSON(http.StatusOK, nil)
}

type updateSubNet struct {
	Data unstructured.Unstructured `json:"data"`
}

// UpdateSubNet Update Service
func (s *networkServer) UpdateSubNet(g *gin.Context) {
	//namespace := g.Param("namespace")
	name := g.Param("name")
	if name == "" {
		common.RequestParametersError(g, fmt.Errorf("params not obtain namespace=%s name=%s", "", name))
		return
	}

	raw, err := g.GetRawData()
	if err != nil {
		common.RequestParametersError(g, fmt.Errorf("get raw data error (%s)", err))
		return
	}

	updateSubNetData := &updateSubNet{}
	if err := json.Unmarshal(raw, updateSubNetData); err != nil {
		common.RequestParametersError(g, fmt.Errorf("unmarshal from json data error (%s)", err))
		return
	}

	newUnstructuredExtend, _, err := s.SubNet.Apply("", name, &service.UnstructuredExtend{Unstructured: &updateSubNetData.Data})
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
