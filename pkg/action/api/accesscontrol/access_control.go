package accesscontrol

import (
	"fmt"
	"github.com/yametech/yamecloud/pkg/action/api"
	"github.com/yametech/yamecloud/pkg/action/service/dac"
)

type accessControlServer struct {
	name string
	*api.Server
	// action services
	*dac.ClusterRole
	*dac.ServiceAccount
	*dac.RoleBinding
	*dac.ClusterRoleBinding
	*dac.Role
	*dac.PSP
}

func (s *accessControlServer) Name() string {
	return s.name
}

func NewApiServer(serviceName string, server *api.Server) *accessControlServer {
	apiServer := &accessControlServer{
		name:   serviceName,
		Server: server,
		// action service
		ClusterRole:        dac.NewClusterRole(server.Interface),
		ServiceAccount:     dac.NewServiceAccount(server.Interface),
		RoleBinding:        dac.NewRoleBinding(server.Interface),
		ClusterRoleBinding: dac.NewClusterRoleBinding(server.Interface),
		Role:               dac.NewRole(server.Interface),
		PSP:                dac.NewPSP(server.Interface),
	}
	group := apiServer.Group(fmt.Sprintf("/%s", serviceName))

	// access control
	//roles
	{
		group.GET("/apis/rbac.authorization.k8s.io/v1/roles", apiServer.ListRole)
		group.GET("/apis/rbac.authorization.k8s.io/v1/roles/:name", apiServer.GetRole)
		group.POST("/apis/rbac.authorization.k8s.io/v1/roles", apiServer.ApplyRole)
	}

	// clusterRole cluster level
	{
		group.GET("/apis/rbac.authorization.k8s.io/v1/clusterroles", apiServer.ListClusterRole)
		group.GET("/apis/rbac.authorization.k8s.io/v1/clusterroles/:name", apiServer.GetClusterRole)
		group.POST("/apis/rbac.authorization.k8s.io/v1/clusterroles", apiServer.ApplyClusterRole)
	}
	// RoleBinding
	{
		group.GET("/apis/rbac.authorization.k8s.io/v1/rolebindings", apiServer.ListRoleBinding)
		group.GET("/apis/rbac.authorization.k8s.io/v1/namespaces/:namespace/rolebindings/:name", apiServer.GetRoleBinding)
		group.POST("/apis/rbac.authorization.k8s.io/v1/namespaces/:namespace/rolebindings", apiServer.ApplyRoleBinding)
	}
	// ClusterRolebind
	{
		group.GET("/apis/rbac.authorization.k8s.io/v1/clusterrolebindings", apiServer.ListClusterRoleBind)
		group.GET("/apis/rbac.authorization.k8s.io/v1/namespaces/:namespace/clusterrolebindings/:name ", apiServer.GetClusterRoleBind)
		group.POST("/apis/rbac.authorization.k8s.io/v1/namespaces/:namespace/clusterrolebindings", apiServer.ApplyClusterRoleBind)
	}

	// service account cluster level
	{
		group.GET("/api/v1/serviceaccounts", apiServer.ListServiceAccount)
		group.GET("/api/v1/namespaces/:namespace/serviceaccounts/:name", apiServer.GetServiceAccount)
		group.POST("/serviceaccount/patch/:method", apiServer.ApplyServiceAccount)
	}

	// #podsecuritypolicies
	{
		group.GET("/apis/policy/v1beta1/podsecuritypolicies", apiServer.ListPSP)
		group.GET("/apis/policy/v1beta1/namespaces/:namespace/podsecuritypolicies/:name", apiServer.GetPSP)
		group.POST("/apis/policy/v1beta1/namespaces/:namespace/podsecuritypolicies", apiServer.ApplyPSP)
	}

	return apiServer
}
