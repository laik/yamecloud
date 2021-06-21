package accesscontrol

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yametech/yamecloud/pkg/action/api/common"
	"github.com/yametech/yamecloud/pkg/action/service"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"net/http"
)

// List Role
func (s *accessControlServer) ListRole(g *gin.Context) {
	list, err := s.Role.List("", "")
	if err != nil {
		common.InternalServerError(g, "", err)
		return
	}
	g.JSON(http.StatusOK, list)
}

// Get Role
func (s *accessControlServer) GetRole(g *gin.Context) {
	namespace := g.Param("namespace")
	name := g.Param("name")

	if name == "" || namespace == "" {
		common.RequestParametersError(g, fmt.Errorf("params not obtain namespace=%s name=%s", namespace, name))
		return
	}

	item, err := s.Role.Get(namespace, name)
	if err != nil {
		common.InternalServerError(g, "", err)
		return
	}
	g.JSON(http.StatusOK, item)
}

// Update or Create Role
func (s *accessControlServer) ApplyRole(g *gin.Context) {
	name := g.Param("name")
	if name == "" {
		common.RequestParametersError(g, fmt.Errorf("params not obtain name=%s", name))
		return
	}
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
	namespace := _unstructured.GetNamespace()
	if namespace == "" {
		common.RequestParametersError(g, fmt.Errorf("data namespace is empty"))
		return
	}

	newUnstructuredExtend, isUpdate, err := s.Role.Apply(namespace, name, &service.UnstructuredExtend{Unstructured: _unstructured})
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
