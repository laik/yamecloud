package tekton

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yametech/yamecloud/pkg/action/api/common"
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/utils"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"net/http"
)

const (
	PipelineRunKind       = "PipelineRun"
	PipelineRunApiVersion = "tekton.dev/v1alpha1"
	GraphAnnotationKey    = "yamecloud.io/tektongraphs"
	RunGraphAnnotationKey = "yamecloud.io/run-tektongraphs"
)

// Get PipelineRun
func (s *tektonServer) GetPipelineRun(g *gin.Context) {
	namespace := g.Param("namespace")
	name := g.Param("name")
	if namespace == "" || name == "" {
		common.RequestParametersError(g, fmt.Errorf("params not obtain namespace=%s name=%s", namespace, name))
		return
	}
	item, err := s.PipelineRun.Get(namespace, name)
	if err != nil {
		common.InternalServerError(g, err, err)
		return
	}
	g.JSON(http.StatusOK, item)
}

// Subscribe PipelineRun
func (s *tektonServer) ListPipelineRun(g *gin.Context) {
	list, err := s.PipelineRun.List(g.Param("namespace"), "")
	if err != nil {
		common.InternalServerError(g, "", err)
		return
	}
	g.JSON(http.StatusOK, list)
}

// Update or Create PipelineRun
func (s *tektonServer) ApplyPipelineRun(g *gin.Context) {
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
	newUnstructuredExtend, isUpdate, err := s.PipelineRun.Apply(namespace, name, &service.UnstructuredExtend{Unstructured: _unstructured})
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

// Delete PipelineRun
func (s *tektonServer) DeletePipelineRun(g *gin.Context) {
	namespace := g.Param("namespace")
	name := g.Param("name")
	if namespace == "" || name == "" {
		common.RequestParametersError(g, fmt.Errorf("params not obtain namespace=%s name=%s", namespace, name))
		return
	}
	err := s.PipelineRun.Delete(namespace, name)
	if err != nil {
		common.InternalServerError(g, err, err)
		return
	}
	g.JSON(http.StatusOK, nil)
}

// Rerun PipelineRun
func (s *tektonServer) RerunPipelineRun(g *gin.Context) {
	namespace := g.Param("namespace")
	name := g.Param("name")
	if namespace == "" || name == "" {
		common.RequestParametersError(g, fmt.Errorf("params not obtain namespace=%s name=%s", namespace, name))
		return
	}
	pipelineRunObj, err := s.PipelineRun.Get(namespace, name)
	if err != nil {
		common.InternalServerError(g, name, fmt.Errorf("get object error (%s)", err))
		return
	}

	oldPipelineRunAnnotations := pipelineRunObj.GetAnnotations()
	tektonGraphName, exist := oldPipelineRunAnnotations[RunGraphAnnotationKey]
	if !exist {
		common.RequestParametersError(g, fmt.Errorf("pipeline run not grpah, pipeline run namespace: %s name: %s", namespace, name))
		return
	}

	oldGraph, err := s.TektonGraph.Get(namespace, tektonGraphName)
	if err != nil {
		common.InternalServerError(g, err, fmt.Errorf("get object error (%s)", err))
		return
	}

	//delete old pipelineRun
	err = s.PipelineRun.Delete(namespace, name)
	if err != nil {
		common.InternalServerError(g, err, fmt.Errorf("delete object error (%s)", err))
		return
	}

	newPipelineRunObj := &service.UnstructuredExtend{Unstructured: &unstructured.Unstructured{Object: map[string]interface{}{}}}
	newPipelineRunObj.SetKind(pipelineRunObj.GetKind())
	newPipelineRunObj.SetAPIVersion(pipelineRunObj.GetAPIVersion())
	newPipelineRunObj.SetName(pipelineRunObj.GetName())
	newPipelineRunObj.SetNamespace(pipelineRunObj.GetNamespace())
	newPipelineRunObj.SetAnnotations(oldPipelineRunAnnotations)
	newPipelineRunObj.SetLabels(pipelineRunObj.GetLabels())

	utils.Set(newPipelineRunObj.Object, "spec", map[string]interface{}{
		"pipelineRef":         utils.Get(pipelineRunObj.Object, "spec.pipelineRef"),
		"pipelineSpec":        utils.Get(pipelineRunObj.Object, "spec.pipelineSpec"),
		"resources":           utils.Get(pipelineRunObj.Object, "spec.resources"),
		"params":              utils.Get(pipelineRunObj.Object, "spec.params"),
		"serviceAccountName":  utils.Get(pipelineRunObj.Object, "spec.serviceAccountName"),
		"serviceAccountNames": utils.Get(pipelineRunObj.Object, "spec.serviceAccountNames"),
		"timeout":             utils.Get(pipelineRunObj.Object, "spec.timeout"),
		"podTemplate":         utils.Get(pipelineRunObj.Object, "spec.podTemplate"),
		"workspaces":          utils.Get(pipelineRunObj.Object, "spec.workspaces"),
	})

	//create new pipelineRun
	newPipelineRun, _, err := s.PipelineRun.Apply(namespace, name, newPipelineRunObj)
	if err != nil {
		common.InternalServerError(g, newPipelineRunObj, fmt.Errorf("apply object error (%s)", err))
		return
	}

	tektonGraph := &service.UnstructuredExtend{Unstructured: &unstructured.Unstructured{Object: map[string]interface{}{}}}
	tektonGraph.SetName(tektonGraphName)
	tektonGraph.SetNamespace(oldGraph.GetNamespace())
	tektonGraph.SetKind(oldGraph.GetKind())
	tektonGraph.SetAPIVersion(oldGraph.GetAPIVersion())
	utils.Set(tektonGraph.Object, "spec", oldGraph.Object["spec"])
	utils.Set(tektonGraph.Object, "metadata.ownerReferences", []map[string]interface{}{
		{
			"apiVersion":         newPipelineRun.GetAPIVersion(),
			"kind":               newPipelineRun.GetKind(),
			"name":               newPipelineRun.GetName(),
			"uid":                newPipelineRun.GetUID(),
			"controller":         true,
			"blockOwnerDeletion": true,
		},
	})

	//create new pipelineRun graph
	_, _, err = s.TektonGraph.Apply(namespace, tektonGraphName, tektonGraph)
	if err != nil {
		common.InternalServerError(g, tektonGraph, fmt.Errorf("apply object error (%s)", err))
		return
	}

	g.JSON(http.StatusOK, newPipelineRun)
}
