package network

import (
	"github.com/gin-gonic/gin"
	"github.com/yametech/yamecloud/pkg/action/api/common"
	"net/http"
)

func (s *networkServer) ListAttachment(g *gin.Context) {
	list, err := s.IP.List(g.Param("namespace"), g.Param("selector"))
	if err != nil {
		common.InternalServerError(g, "", err)
		return
	}
	g.JSON(http.StatusOK, list)
}
