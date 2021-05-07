package base

import (
	"github.com/gin-gonic/gin"
	"github.com/yametech/yamecloud/pkg/permission"
	"net/http"
)

// Tree Permission
func (s *baseServer) treePermission(g *gin.Context) {
	g.JSON(http.StatusOK, permission.Tree)
}
