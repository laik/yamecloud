package editer

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yametech/yamecloud/pkg/action/api"
	"github.com/yametech/yamecloud/pkg/action/api/common"
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/action/service/editer"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"net/http"
	"strings"
)

type editerServer struct {
	name string
	*api.Server
	// action types
	*editer.APIResources
}

func (s *editerServer) Name() string { return s.name }

func NewEditerServer(serviceName string, server *api.Server) *editerServer {
	editerServer := &editerServer{
		name:   serviceName,
		Server: server,

		APIResources: editer.NewAPIResources(server),
	}
	group := editerServer.Group(fmt.Sprintf("/%s", serviceName))

	{
		group.POST("/stack", editerServer.Stack)
	}

	_ = group
	return editerServer
}

func (s *editerServer) Stack(g *gin.Context) {
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

	apiResources, err := s.APIResources.ListAPIResources()
	if err != nil {
		common.InternalServerError(g, "", fmt.Errorf("list api resources error (%s)", err))
		return
	}

	uniqueKey := fmt.Sprintf("%s/%s", _unstructured.GetAPIVersion(), _unstructured.GetKind())
	if strings.Count(uniqueKey, "/") == 1 {
		uniqueKey = fmt.Sprintf("%s%s", "/", uniqueKey)
	}

	gvr, ok := apiResources[uniqueKey]
	if gvr == nil || !ok {
		common.RequestParametersError(g, fmt.Errorf("unknown resource apiVersion %s kind %s", _unstructured.GetAPIVersion(), _unstructured.GetKind()))
		return
	}

	namespace := _unstructured.GetNamespace()
	name := _unstructured.GetName()

	newUnstructuredExtend, isUpdate, err := s.APIResources.ApplyGVR(namespace, name, gvr, &service.UnstructuredExtend{Unstructured: _unstructured})

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
