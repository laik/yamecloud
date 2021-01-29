package gateway

import (
	"github.com/gin-gonic/gin"
	"github.com/yametech/yamecloud/pkg/action/api"
	"github.com/yametech/yamecloud/pkg/action/api/common"
	"net/http"
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
	//auth := NewAuthorization(server.Interface)
	//server.Use(
	//	IsNeedSkip(auth),
	//	ValidateToken(auth),
	//	IsAdmin(auth),
	//	IsTenantOwner(auth),
	//	IsDepartmentOwner(auth),
	//	CheckNamespace(auth),
	//	CheckPermission(auth),
	//	IsWithGranted(auth),
	//)
	server.POST("/user-login", gatewayServer.userLogin)
	server.GET("/config", gatewayServer.userConfig)
	//server.Any("/*any", func(g *gin.Context) {
	//	if g.Request.RequestURI == "/user-login" {
	//		gatewayServer.userLogin(g)
	//		return
	//	}
	//
	//	if g.Request.RequestURI == "/config" {
	//		gatewayServer.userConfig(g)
	//		return
	//	}
	//	g.Next()
	//})

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
		g.JSON(http.StatusBadRequest, gin.H{"msg": "账号或密码错误"})
		return
	}
	if userConfig == nil {
		g.JSON(http.StatusBadRequest, gin.H{"msg": "账号或密码错误"})
		return
	}
	g.JSON(http.StatusOK, userConfig)
}
