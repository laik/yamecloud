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

func (w *workloadServer) GetReplicaSet(g *gin.Context) {
	namespace := g.Param("namespace")
	name := g.Param("name")
	if namespace == "" || name == "" {
		common.RequestParametersError(g, fmt.Errorf("params not obtain namespace=%s name=%s", namespace, name))
		return
	}
	item, err := w.ReplicaSet.Get(namespace, name)
	if err != nil {
		common.InternalServerError(g, err, err)
		return
	}
	g.JSON(http.StatusOK, item)
}

func (w *workloadServer) ListReplicaSet(g *gin.Context) {
	list, err := w.ReplicaSet.List("", "")
	if err != nil {
		common.InternalServerError(g, "", err)
		return
	}
	g.JSON(http.StatusOK, list)
}

func (w *workloadServer) ApplyReplicaSet(g *gin.Context) {
	namespace := g.Param("namespace")
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
	newUnstructuredExtend, isUpdate, err := w.ReplicaSet.Apply(namespace, name, &service.UnstructuredExtend{Unstructured: _unstructured})
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

func (w *workloadServer) UpdateReplicaSet(g *gin.Context) {
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

	updateNetWorkAttachmentData := &unstructured.Unstructured{}
	if err := json.Unmarshal(raw, updateNetWorkAttachmentData); err != nil {
		common.RequestParametersError(g, fmt.Errorf("unmarshal from json data error (%s)", err))
		return
	}

	newUnstructuredExtend, _, err := w.ReplicaSet.Apply("", name, &service.UnstructuredExtend{Unstructured: updateNetWorkAttachmentData})
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

func (w *workloadServer) DeleteReplicaSet(g *gin.Context) {
	namespace := g.Param("namespace")
	name := g.Param("name")
	if namespace == "" || name == "" {
		common.RequestParametersError(g, fmt.Errorf("params not obtain namespace=%s name=%s", namespace, name))
		return
	}
	err := w.ReplicaSet.Delete(namespace, name)
	if err != nil {
		common.InternalServerError(g, err, err)
		return
	}
	g.JSON(http.StatusOK, nil)
}
