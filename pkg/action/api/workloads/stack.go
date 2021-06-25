package workloads

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yametech/yamecloud/pkg/action/api/common"
	"github.com/yametech/yamecloud/pkg/action/service"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"net/http"
)

func (w *workloadServer) Stack(g *gin.Context) {
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

	apiResources, err := w.APIResources.ListAPIResources()
	if err != nil {
		common.InternalServerError(g, "", fmt.Errorf("list api resources error (%s)", err))
		return
	}

	uniqueKey := fmt.Sprintf("%s/%s", _unstructured.GetAPIVersion(), _unstructured.GetKind())

	gvr, ok := apiResources[uniqueKey]
	if gvr == nil || !ok {
		common.RequestParametersError(g, fmt.Errorf("unknown resource apiVersion %s kind %s", _unstructured.GetAPIVersion(), _unstructured.GetKind()))
		return
	}

	namespace := _unstructured.GetNamespace()
	name := _unstructured.GetName()

	newUnstructuredExtend, isUpdate, err := w.APIResources.ApplyGVR(namespace, name, gvr, &service.UnstructuredExtend{Unstructured: _unstructured})

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
