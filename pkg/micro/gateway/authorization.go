package gateway

import (
	"fmt"
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/action/service/tenant"
	"github.com/yametech/yamecloud/pkg/k8s"
	"github.com/yametech/yamecloud/pkg/uri"
	"github.com/yametech/yamecloud/pkg/utils"
	"net/http"
	"strings"
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
	ValidateToken(token string) (*CustomClaims, error)
	IsAdmin(userName string) (bool, error)
	IsTenantOwner(userName string) (bool, error)
	IsDepartmentOwner(userName string) (bool, error)
	IsWithGranted(userName string) (bool, error)
	CheckPermission(userName string, op *uri.URI) (bool, error)
	CheckNamespace(userName, namespace string, isAdmin, isTenantOwner, isDepartmentOwner bool) (bool, error)
}

var _ IAuthorization = (*Authorization)(nil)

var adminList = []string{"admin"}
var excludeMap = map[string]string{
	"/user-login":           http.MethodPost,
	"/workload/watch":       http.MethodGet,
	"/base/permission_tree": http.MethodGet,
}

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
	if strings.Contains(path, "/workload/shell/pod") ||
		strings.Contains(path, "/webhook") {
		return true, nil
	}
	return false, nil
}

func (auth *Authorization) ValidateToken(token string) (*CustomClaims, error) {
	decodeToken, err := (&Token{}).Decode(token)
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
	isTenantValue, err := userObj.Get("spec.is_tenant_owner")
	if err != nil {
		return false, err
	}
	if isTenantValue == nil {
		return false, nil
	}
	isTenant := isTenantValue.(bool)
	return isTenant, nil
}

//check whether a user is a department owner
func (auth *Authorization) IsDepartmentOwner(userName string) (bool, error) {
	userObj, err := auth.userServices.Get("kube-system", userName)
	if err != nil {
		return false, err
	}
	specTenantId, err := userObj.Get("spec.tenant_id")
	if err != nil {
		return false, err
	}
	if specTenantId == nil || specTenantId == "" {
		return false, nil
	}
	tenantId := specTenantId.(string)
	selector := fmt.Sprintf("tenant.yamecloud.io=%s", tenantId)
	deptObjList, err := auth.deptServices.List("kube-system", selector)
	if err != nil {
		return false, err
	}
	if len(deptObjList.Items) == 0 {
		return false, nil
	}
	for _, item := range deptObjList.Items {
		specOwner := utils.Get(item.Object, "spec.owner")
		if specOwner == "" {
			continue
		}
		owner := specOwner.(string)
		if owner == userName {
			return true, nil
		}
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
func (auth *Authorization) CheckPermission(userName string, uri *uri.URI) (bool, error) {
	roles, err := auth.getUserRoles(userName)
	if err != nil {
		return false, err
	}
	resource := uri.Resource
	if resource == k8s.StatefulSet {
		if uri.Group == "nuwa.nip.io" && uri.Version == "v1" {
			resource = k8s.StatefulSet1
		}
	}
	for _, role := range roles {
		roleObj, err := auth.roleServices.Get("kube-system", role)
		if err != nil {
			continue
		}
		specPrivilege, err := roleObj.Get("spec.privilege")
		if err != nil || specPrivilege == nil {
			continue
		}
		privilege := specPrivilege.(map[k8s.ResourceType]interface{})
		resourcePrivilege, exist := privilege[resource]
		if !exist || len(resourcePrivilege.([]interface{})) == 0 {
			continue
		}
		for _, item := range resourcePrivilege.([]interface{}) {
			if item.(string) == uri.Op {
				return true, nil
			}
		}
	}

	return false, nil
}

//check whether a user allow access specified namespace
func (auth *Authorization) CheckNamespace(userName, namespace string, isAdmin, isTenantOwner, isDepartmentOwner bool) (bool, error) {

	if isAdmin {
		return true, nil
	}
	allowNamespaces, err := auth.AllowNamespaces(userName, isAdmin, isTenantOwner, isDepartmentOwner)
	if err != nil {
		return false, err
	}

	if In(allowNamespaces, namespace) {
		return true, nil
	}
	return false, nil
}

func (auth *Authorization) AllowNamespaces(username string, isAdmin, isTenantOwner, isDepartmentOwner bool) ([]string, error) {
	if isAdmin {
		return nil, nil
	}

	userObj, _ := auth.userServices.Get("kube-system", username)
	// user AllowedNamespaces
	allowNamespaces := make([]string, 0)
	if isTenantOwner {
		specTenantId, err := userObj.Get("spec.tenant_id")
		if err != nil {
			return nil, fmt.Errorf("get tenant_id error")
		}
		if specTenantId == "" {
			return nil, fmt.Errorf("tenant do not exist")
		}
		tenantObj, err := auth.tenantServices.Get("kube-system", specTenantId.(string))
		if err != nil {
			return nil, fmt.Errorf("query tenant error")
		}
		tenantNamespaces, err := tenantObj.Get("spec.namespaces")
		if err != nil {
			return nil, fmt.Errorf("get tenant namespaces error")
		}
		if tenantNamespaces != nil && len(tenantNamespaces.([]interface{})) > 0 {
			for _, item := range tenantNamespaces.([]interface{}) {
				allowNamespaces = append(allowNamespaces, item.(string))

			}
		}
		return allowNamespaces, nil
	}
	if isDepartmentOwner {
		specTenantId, err := userObj.Get("spec.tenant_id")
		if err != nil {
			return nil, err
		}
		if specTenantId == "" {
			return nil, nil
		}
		tenantId := specTenantId.(string)
		selector := fmt.Sprintf("tenant.yamecloud.io=%s", tenantId)
		deptObjList, err := auth.deptServices.List("kube-system", selector)
		if err != nil {
			return nil, err
		}
		if len(deptObjList.Items) > 0 {
			for _, item := range deptObjList.Items {
				specOwner := utils.Get(item.Object, "spec.owner")
				if specOwner == "" {
					continue
				}
				owner := specOwner.(string)
				if owner == username {
					deptNamespaces := utils.Get(item.Object, "spec.namespaces")
					for _, item1 := range deptNamespaces.([]interface{}) {
						if !In(allowNamespaces, item1.(string)) {
							allowNamespaces = append(allowNamespaces, item1.(string))
						}
					}

				}
			}
		}
	}
	specRoles, err := userObj.Get("spec.roles")
	if err != nil {
		return nil, err
	}
	if specRoles != nil {
		for _, item := range specRoles.([]interface{}) {
			roleObj, err := auth.roleServices.Get("kube-system", item.(string))
			if err != nil {
				continue
			}
			specNamespaces, err := roleObj.Get("spec.namespaces")
			if err != nil {
				continue
			}
			for _, item1 := range specNamespaces.([]interface{}) {
				if !In(allowNamespaces, item1.(string)) {
					allowNamespaces = append(allowNamespaces, item1.(string))
				}
			}

		}
	}
	return allowNamespaces, nil
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
	if specRoles == nil || len(specRoles.([]interface{})) == 0 {
		return nil, nil
	}
	var roles []string
	for _, item := range specRoles.([]interface{}) {
		roles = append(roles, item.(string))
	}
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
