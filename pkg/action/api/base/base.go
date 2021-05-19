package base

import (
	"fmt"
	"github.com/yametech/yamecloud/pkg/action/api"
	"github.com/yametech/yamecloud/pkg/action/service/tenant"
)

type baseServer struct {
	name string
	*api.Server
	// action services
	*tenant.BaseTenant
	*tenant.BaseDepartment
	*tenant.BaseUser
	*tenant.BaseRole
	*tenant.BaseRoleUser
}

func (s *baseServer) Name() string { return s.name }

func NewBaseServer(serviceName string, server *api.Server) *baseServer {
	baseServer := &baseServer{
		name:           serviceName,
		Server:         server,
		BaseTenant:     tenant.NewBaseTenant(server.Interface),
		BaseDepartment: tenant.NewBaseDepartment(server.Interface),
		BaseUser:       tenant.NewBaseUser(server.Interface),
		BaseRole:       tenant.NewBaseRole(server.Interface),
		BaseRoleUser:   tenant.NewBaseRoleUser(server.Interface),
	}
	group := baseServer.Group(fmt.Sprintf("/%s", serviceName))

	// BaseTenant
	{
		group.GET("/apis/yamecloud.io/v1/basetenants", baseServer.ListBaseTenant)
		group.GET("/apis/yamecloud.io/v1/basetenants/:name", baseServer.GetBaseTenant)
		group.POST("/apis/yamecloud.io/v1/basetenants", baseServer.ApplyBaseTenant)
		group.DELETE("/apis/yamecloud.io/v1/basetenants/:name", baseServer.DeleteBaseTenant)
	}

	// BaseDepartment
	{
		group.GET("/apis/yamecloud.io/v1/basedepartments", baseServer.ListBaseDepartment)
		group.GET("/apis/yamecloud.io/v1/basedepartments/:name", baseServer.GetBaseDepartment)
		group.POST("/apis/yamecloud.io/v1/basedepartments", baseServer.ApplyBaseDepartment)
		group.DELETE("/apis/yamecloud.io/v1/basedepartments/:name", baseServer.DeleteBaseDepartment)
	}

	// BaseRole
	{
		group.GET("/apis/yamecloud.io/v1/baseroles", baseServer.ListBaseRole)
		group.GET("/apis/yamecloud.io/v1/baseroles/:name", baseServer.GetBaseRole)
		group.POST("/apis/yamecloud.io/v1/baseroles", baseServer.ApplyBaseRole)
		group.PUT("/apis/yamecloud.io/v1/baseroles/:name", baseServer.UpdateBaseRole)
		group.DELETE("/apis/yamecloud.io/v1/baseroles/:name", baseServer.DeleteBaseRole)
	}

	// BaseUser
	{
		group.GET("/apis/yamecloud.io/v1/baseusers", baseServer.ListBaseUser)
		group.GET("/apis/yamecloud.io/v1/baseusers/:name", baseServer.GetBaseUser)
		group.POST("/apis/yamecloud.io/v1/baseusers", baseServer.ApplyBaseUser)
		group.PUT("/apis/yamecloud.io/v1/baseusers/:name", baseServer.UpdateBaseUser)
		group.DELETE("/apis/yamecloud.io/v1/baseusers/:name", baseServer.DeleteBaseUser)
	}

	// BaseRoleUser
	{
		group.GET("/apis/yamecloud.io/v1/baseroleusers", baseServer.ListBaseRoleUser)
		group.GET("/apis/yamecloud.io/v1/baseroleusers/:name", baseServer.GetBaseRoleUser)
		group.POST("/apis/yamecloud.io/v1/baseroleusers", baseServer.ApplyBaseRoleUser)
		group.DELETE("/apis/yamecloud.io/v1/baseroleusers/:name", baseServer.DeleteBaseRoleUser)
	}

	//permission
	{
		group.GET("/permission_tree", baseServer.treePermission)
	}

	return baseServer
}
