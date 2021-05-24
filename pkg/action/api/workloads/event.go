package workloads

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yametech/yamecloud/pkg/action/api/common"
	"net/http"
	"strconv"
)

func (s *workloadServer) GetEvent(g *gin.Context) {
	namespace := g.Param("namespace")
	name := g.Param("name")
	if namespace == "" || name == "" {
		common.RequestParametersError(g, fmt.Errorf("params not obtain namespace=%s name=%s", namespace, name))
		return
	}
	item, err := s.Event.Get(namespace, name)
	if err != nil {
		common.InternalServerError(g, err, err)
		return
	}
	g.JSON(http.StatusOK, item)
}

func (s *workloadServer) ListEvent(g *gin.Context) {
	limit := g.Param("limit")
	namespace := g.Param("namespace")
	flag := g.ClientIP()

	limitNum := int64(10000)
	var err error
	if limit != "" {
		limitNum, err = strconv.ParseInt(limit, 64, 10)
		if err != nil {
			common.RequestParametersError(g, fmt.Errorf("params not obtain namespace=%s", namespace))
			return
		}
	}
	list, err := s.Event.ListLimit(namespace, flag, 0, limitNum, "")
	if err != nil {
		common.InternalServerError(g, "", err)
		return
	}
	g.JSON(http.StatusOK, list)
}
