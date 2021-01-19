package tekton

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yametech/yamecloud/pkg/action/api/common"
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/utils"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"net/http"
	"strconv"
	"time"
)

// Get Pipeline
func (s *tektonServer) GetPipeline(g *gin.Context) {
	namespace := g.Param("namespace")
	name := g.Param("name")
	if namespace == "" || name == "" {
		common.RequestParametersError(g, fmt.Errorf("params not obtain namespace=%s name=%s", namespace, name))
		return
	}
	item, err := s.Pipeline.Get(namespace, name)
	if err != nil {
		common.InternalServerError(g, err, err)
		return
	}
	g.JSON(http.StatusOK, item)
}

// Subscribe Pipeline
func (s *tektonServer) ListPipeline(g *gin.Context) {
	list, err := s.Pipeline.List(g.Param("namespace"), "")
	if err != nil {
		common.InternalServerError(g, "", err)
		return
	}
	g.JSON(http.StatusOK, list)
}

// Update or Create Pipeline
func (s *tektonServer) ApplyPipeline(g *gin.Context) {
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
	newUnstructuredExtend, isUpdate, err := s.Pipeline.Apply(namespace, name, &service.UnstructuredExtend{Unstructured: _unstructured})
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

// Delete Pipeline
func (s *tektonServer) DeletePipeline(g *gin.Context) {
	namespace := g.Param("namespace")
	name := g.Param("name")
	if namespace == "" || name == "" {
		common.RequestParametersError(g, fmt.Errorf("params not obtain namespace=%s name=%s", namespace, name))
		return
	}
	err := s.Pipeline.Delete(namespace, name)
	if err != nil {
		common.InternalServerError(g, err, err)
		return
	}
	g.JSON(http.StatusOK, nil)
}

// Run Pipeline
func (s *tektonServer) RunPipeline(g *gin.Context) {
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
	_unstructured := &unstructured.Unstructured{}
	if err := _unstructured.UnmarshalJSON(raw); err != nil {
		common.RequestParametersError(g, fmt.Errorf("unmarshal from json data error (%s)", err))
		return
	}

	//create tektonGraph
	pipelineObj, err := s.Pipeline.Get(namespace, name)
	if err != nil {
		common.InternalServerError(g, name, fmt.Errorf("get object error (%s)", err))
		return
	}
	var width int64 = 1000
	var height int64 = 1000
	graphJson := ""
	graphName := pipelineObj.GetAnnotations()[GraphAnnotationKey]
	graphObj, err := s.TektonGraph.Get(namespace, graphName)
	if err != nil {
		common.InternalServerError(g, graphName, fmt.Errorf("get object error (%s)", err))
		return
	}
	if graphObj != nil {
		specWidth, err := graphObj.Get("spec.width")
		if err == nil {
			width = specWidth.(int64)
		}
		specHeight, err := graphObj.Get("spec.height")
		if err == nil {
			height = specHeight.(int64)
		}
		graphData, err := graphObj.Get("spec.data")
		if err == nil {
			graphJson = graphData.(string)
		} else {
			bytes, err := json.Marshal(GetDefaultInitData())
			if err == nil {
				graphJson = string(bytes)
			}
		}

	}
	runGraphName := "run-" + name + strconv.FormatInt(time.Now().Unix(), 10)
	runGraphObj := &service.UnstructuredExtend{Unstructured: &unstructured.Unstructured{Object: map[string]interface{}{}}}
	runGraphObj.SetKind(TektonGraphKind)
	runGraphObj.SetAPIVersion(TektonGraphApiVersion)
	runGraphObj.SetNamespace(namespace)
	runGraphObj.SetName(runGraphName)
	runGraphObj.SetLabels(map[string]string{"namespace": pipelineObj.GetLabels()["namespace"]})

	utils.Set(runGraphObj.Object, "spec", map[string]interface{}{
		"data":   graphJson,
		"width":  width,
		"height": height,
	})
	runGraphObjBack, _, err := s.TektonGraph.Apply(namespace, runGraphName, runGraphObj)
	if err != nil {
		common.InternalServerError(g, runGraphObj, fmt.Errorf("apply object error (%s)", err))
		return
	}
	//create pipelineRun
	pipelineRunObj := &service.UnstructuredExtend{Unstructured: _unstructured}
	pipelineRunObj.SetNamespace(namespace)
	pipelineRunObj.SetAnnotations(map[string]string{
		RunGraphAnnotationKey: runGraphName,
	})

	appliedPipelineRunObj, _, err := s.PipelineRun.Apply(namespace, pipelineRunObj.GetName(), pipelineRunObj)
	if err != nil {
		common.InternalServerError(g, pipelineRunObj, fmt.Errorf("apply object error (%s)", err))
		return
	}
	utils.Set(runGraphObjBack.Object, "metadata.ownerReferences", []map[string]interface{}{
		{
			"apiVersion":         appliedPipelineRunObj.GetAPIVersion(),
			"kind":               appliedPipelineRunObj.GetKind(),
			"name":               pipelineObj.GetName(),
			"uid":                appliedPipelineRunObj.GetUID(),
			"controller":         false,
			"blockOwnerDeletion": false,
		},
	})

	//update tektonGraph
	appliedRunGraphObj, _, err := s.TektonGraph.Apply(namespace, runGraphObjBack.GetName(), runGraphObjBack)
	if err != nil {
		common.InternalServerError(g, runGraphObjBack, fmt.Errorf("apply object error (%s)", err))
		return
	}
	_ = appliedRunGraphObj

	g.JSON(http.StatusOK, appliedPipelineRunObj)

}
