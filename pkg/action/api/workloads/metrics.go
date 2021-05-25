package workloads

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yametech/yamecloud/pkg/action/api/common"
	"net/http"
)

func (s *workloadServer) DefaultMetrics(g *gin.Context) {
	body, err := g.GetRawData()
	if err != nil {
		common.RequestParametersError(g, fmt.Errorf("params not obtain or params parse error: %s", err))
		return
	}

	params := make(map[string]string)
	params["start"] = g.Query("start")
	params["end"] = g.Query("end")
	params["step"] = g.Query("step")
	params["kubernetes_namespace"] = g.Query("kubernetes_namespace")

	bufRaw, err := s.Metrics.ProxyToPrometheus(params, body)
	if err != nil {
		common.InternalServerError(g, "", err)
		return
	}

	g.JSON(http.StatusOK, bufRaw)
}

func (s *workloadServer) ListNodeMetrics(g *gin.Context) {
	result, err := s.Metrics.NodeMetricsList()
	if err != nil {
		common.InternalServerError(g, err, fmt.Errorf("node metrics list error: %s", err))
		return
	}
	g.JSON(http.StatusOK, result)
}

func (s *workloadServer) GetPodMetrics(g *gin.Context) {
	namespace := g.Query("namespace")
	name := g.Query("name")
	result, err := s.Metrics.PodMetrics(namespace, name)
	if err != nil {
		common.InternalServerError(g, err, fmt.Errorf("get pod metrics error: %s,namespace&pod:  %s/%s ", name, namespace,
			err))
		return
	}
	g.JSON(http.StatusOK, result)
}

func (s *workloadServer) ListPodMetrics(g *gin.Context) {
	namespace := g.Query("namespace")
	result, err := s.Metrics.PodMetricsList(namespace)
	if err != nil {
		common.InternalServerError(g, err, fmt.Errorf("list pod metrics error: %s,namespace:%s", err, namespace))
		return
	}
	g.JSON(http.StatusOK, result)
}
