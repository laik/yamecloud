package gateway

import (
	"github.com/gin-gonic/gin"
	"github.com/yametech/yamecloud/pkg/action/api"
	"github.com/yametech/yamecloud/pkg/action/api/common"
	"github.com/yametech/yamecloud/pkg/action/service/tenant"
	"net/http"
)

type gatewayServer struct {
	name string
	*api.Server
	// action services
	*tenant.BaseDepartment
}

func (gw *gatewayServer) Name() string { return gw.name }

func NewGatewayServer(serviceName string, server *api.Server) *gatewayServer {
	gatewayServer := &gatewayServer{
		name:   serviceName,
		Server: server,
		// action service
		BaseDepartment: tenant.NewBaseDepartment(server.Interface),
	}

	server.Any("/*any", func(g *gin.Context) {
		if g.Request.RequestURI == "/user-login" {
			gatewayServer.userLogin(g)
			return
		}

		if g.Request.RequestURI == "/config" {
			gatewayServer.userConfig(g)
			return
		}
	})

	return gatewayServer
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (gw *gatewayServer) userConfig(g *gin.Context) {

	g.JSON(http.StatusOK, nil)
}

func (gw *gatewayServer) userLogin(g *gin.Context) {
	user := &User{}
	if err := g.ShouldBindJSON(user); err != nil {
		common.RequestParametersError(g, err)
		return
	}
	g.JSON(http.StatusOK, nil)
}
