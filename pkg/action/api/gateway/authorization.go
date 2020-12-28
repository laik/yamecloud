package gateway

import (
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/action/service/tenant"
	"github.com/yametech/yamecloud/pkg/k8s"
	"github.com/yametech/yamecloud/pkg/micro/gateway"
	"github.com/yametech/yamecloud/pkg/uri"
)

type Identification string

const (
	Admin           Identification = "admin"
	TenantOwner     Identification = "tenantOwner"
	DepartmentOwner Identification = "tenantOwner"
	OrdinaryUser    Identification = "ordinaryUser"
)

type IAuthorization interface {
	IsNeedSkip(method, path string) (bool, error)
	ValidateToken(token string) (*gateway.CustomClaims, error)
	IsAdmin(userName string) (bool, error)
	IsTenantOwner(userName string) (bool, error)
	IsDepartmentOwner(userName string) (bool, error)
	IsWithGranted(userName string) (bool, error)
	CheckPermission(userName string, op *uri.Op) (bool, error)
	CheckNamespace(userName, namespace string) (bool, error)
}

var _ IAuthorization = (*Authorization)(nil)

var adminList = []string{"admin"}
var excludeMap = map[string]string{"/user-login": "POST"}

type Authorization struct {
	userServices     *tenant.BaseUser
	roleUserServices *tenant.BaseRoleUser
	roleServices     *tenant.BaseRole
	deptServices     *tenant.BaseDepartment
	tenantServices   *tenant.BaseTenant
}

func NewAuthorization(svcInterface service.Interface) *Authorization {
	return &Authorization{
		userServices:     tenant.NewBaseUser(svcInterface),
		roleUserServices: tenant.NewBaseRoleUser(svcInterface),
		roleServices:     tenant.NewBaseRole(svcInterface),
		deptServices:     tenant.NewBaseDepartment(svcInterface),
		tenantServices:   tenant.NewBaseTenant(svcInterface),
	}
}

func (auth *Authorization) IsNeedSkip(method, path string) (bool, error) {
	value := excludeMap[path]
	if method == value {
		return true, nil
	}
	return false, nil
}

func (auth *Authorization) ValidateToken(token string) (*gateway.CustomClaims, error) {
	decodeToken, err := (&gateway.Token{}).Decode(token)
	if err != nil {
		return nil, err
	}
	return decodeToken, nil
}

//check whether a user is an admin
func (auth *Authorization) IsAdmin(userName string) (bool, error) {
	for _, _item := range adminList {
		if userName == _item {
			return true, nil
		}
	}
	return false, nil
}

//check whether a user is a tenant owner
func (auth *Authorization) IsTenantOwner(userName string) (bool, error) {
	userObj, err := auth.userServices.Get("kube-system", userName)
	if err != nil {
		return false, err
	}
	isTenantValue, err := userObj.Get("spec.is_tenant")
	if err != nil {
		return false, err
	}
	isTenant := isTenantValue.(bool)
	return isTenant, nil
}

//check whether a user is a department owner
func (auth *Authorization) IsDepartmentOwner(userName string) (bool, error) {
	deptObjList, err := auth.deptServices.List("kube-system", "spec.department_owner="+userName)
	if err != nil {
		return false, err
	}
	if len(deptObjList.Items) > 0 {
		return true, nil
	}

	return false, nil

}

//check whether a user is with granted
func (auth *Authorization) IsWithGranted(userName string) (bool, error) {
	userObj, err := auth.userServices.Get("kube-system", userName)
	if err != nil {
		return false, err
	}
	specIsWithGranted, err := userObj.Get("spec.is_with_granted")
	if err != nil {
		return false, err
	}
	isWithGranted := specIsWithGranted.(bool)
	return isWithGranted, nil
}

//check whether a user has specified uri permission
func (auth *Authorization) CheckPermission(userName string, op *uri.Op) (bool, error) {
	roles, err := auth.getUserRoles(userName)
	if err != nil {
		return false, err
	}

	for _, role := range roles {
		roleObj, err := auth.roleServices.Get("kube-system", role)
		if err != nil {
			break
		}
		specPrivilege, err := roleObj.Get("spec.privilege")
		if err != nil {
			break
		}
		privilege := specPrivilege.(map[k8s.ResourceType]uint16)
		resourcePrivilege := privilege[op.Resource]
		if resourcePrivilege > 0 && (resourcePrivilege&1<<op.Type) > 0 {
			return true, nil
		}

	}

	return false, nil
}

//check whether a user allow access specified namespace
func (auth *Authorization) CheckNamespace(userName, namespace string) (bool, error) {
	roles, err := auth.getUserRoles(userName)
	if err != nil {
		return false, err
	}

	for _, role := range roles {
		roleObj, err := auth.roleServices.Get("kube-system", role)
		if err != nil {
			break
		}
		specNamespace, err := roleObj.Get("spec.namespace")
		if err != nil {
			break
		}
		namespaces := specNamespace.([]string)
		if In(namespaces, namespace) {
			return true, nil
		}
	}

	return false, nil
}

func (auth *Authorization) getUserRoles(username string) ([]string, error) {
	userObj, err := auth.userServices.Get("kube-system", username)
	if err != nil {
		return nil, err
	}
	specRoles, err := userObj.Get("spec.roles")
	if err != nil {
		return nil, err
	}
	roles := specRoles.([]string)
	return roles, nil
}

func In(list []string, item string) bool {
	for _, _item := range list {
		if _item == item {
			return true
		}
	}
	return false
}
