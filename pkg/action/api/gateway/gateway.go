package gateway

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yametech/yamecloud/pkg/action/api"
	"github.com/yametech/yamecloud/pkg/action/api/common"
)

type gatewayServer struct {
	name string
	*api.Server
	// action services
	loginHandle *LoginHandle
}

func (gw *gatewayServer) Name() string { return gw.name }

func NewGatewayServer(serviceName string, server *api.Server) *gatewayServer {
	gatewayServer := &gatewayServer{
		name:   serviceName,
		Server: server,
		// action service
		loginHandle: NewLoginHandle(server.Interface),
	}
	server.POST("/user-login", gatewayServer.userLogin)
	server.GET("/config", gatewayServer.userConfig)

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
	userConfig, err := gw.loginHandle.Auth(user)
	if err != nil {
		g.JSON(http.StatusBadRequest, "incorrect username or password")
		return
	}
	if userConfig == nil {
		g.JSON(http.StatusBadRequest, "incorrect username or password")
		return
	}
	g.JSON(http.StatusOK, userConfig)
}
