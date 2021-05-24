package workloads

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yametech/yamecloud/pkg/action/api/common"
	"net/http"
)

func (s *workloadServer) NamespacesMetrics(g *gin.Context) {
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
