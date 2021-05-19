package network

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yametech/yamecloud/pkg/action/api/common"
	"net/http"
)

// GetIP none
func (s *networkServer) GetIP(g *gin.Context) {
	name := g.Param("name")
	if name == "" {
		common.RequestParametersError(g, fmt.Errorf("params not obtain namespace=%s name=%s", "", name))
		return
	}
	item, err := s.IP.Get("", name)
	if err != nil {
		common.InternalServerError(g, err, err)
		return
	}
	g.JSON(http.StatusOK, item)
}

// ListIP list all subnet
func (s *networkServer) ListIP(g *gin.Context) {
	list, err := s.IP.List("", "")
	if err != nil {
		common.InternalServerError(g, "", err)
		return
	}
	g.JSON(http.StatusOK, list)
}

// DeleteIP none
func (s *networkServer) DeleteIP(g *gin.Context) {
	name := g.Param("name")
	if name == "" {
		common.RequestParametersError(g, fmt.Errorf("params not obtain namespace=%s name=%s", "", name))
		return
	}
	err := s.SubNet.Delete("", name)
	if err != nil {
		common.InternalServerError(g, err, err)
		return
	}
	g.JSON(http.StatusOK, nil)
}
