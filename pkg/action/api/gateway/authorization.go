package gateway

import (
	"encoding/json"
	"fmt"
	"github.com/yametech/yamecloud/common"
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/action/service/tenant"
	v1 "github.com/yametech/yamecloud/pkg/apis/yamecloud/v1"
	"github.com/yametech/yamecloud/pkg/k8s"
	"github.com/yametech/yamecloud/pkg/micro/gateway"
	"github.com/yametech/yamecloud/pkg/permission"
	"k8s.io/apimachinery/pkg/runtime"
	"time"
)

type Authorization struct {
	userServices     *tenant.BaseUser
	roleUserServices *tenant.BaseRoleUser
	roleServices     *tenant.BaseRole
	deptServices     *tenant.BaseDepartment
	tenantServices   *tenant.BaseTenant
}

func NewAuthorization(svcInterface service.Interface) *Authorization {
	auth := &Authorization{
		userServices:     tenant.NewBaseUser(svcInterface),
		roleUserServices: tenant.NewBaseRoleUser(svcInterface),
		roleServices:     tenant.NewBaseRole(svcInterface),
		deptServices:     tenant.NewBaseDepartment(svcInterface),
		tenantServices:   tenant.NewBaseTenant(svcInterface),
	}
	return auth
}

func (auth *Authorization) allowNamespaceAccess(userName string, namespace string) (bool, error) {
	obj, err := auth.userServices.Get("kube-system", userName)
	if err != nil {
		return false, err
	}
	ueRoles, err := obj.Get("spec.roles")
	if err != nil {
		return false, err
	}
	roles := ueRoles.([]string)
	for _, item := range roles {
		ueRole, err := auth.roleServices.Get("kube-system", item)
		if err != nil {
			break
		}
		ueNamespace, err := ueRole.Get("spec.namespace")
		if err != nil {
			break
		}
		namespaces := ueNamespace.([]string)
		if In(namespaces, namespace) {
			return true, nil
		}
	}

	return false, nil
}

func In(list []string, item string) bool {
	for _, _item := range list {
		if _item == item {
			return true
		}
	}
	return false
}

func (auth *Authorization) getUser(userName string) (*v1.BaseUser, error) {
	obj, err := auth.userServices.Get("kube-system", userName)
	if err != nil {
		return nil, err
	}
	baseUser := &v1.BaseUser{}
	err = runtimeObjectToInstanceObj(obj, baseUser)
	if err != nil {
		return nil, err
	}
	return baseUser, nil
}

func (auth *Authorization) getDept(deptName string) (*v1.BaseDepartment, error) {
	obj, err := auth.deptServices.Get("kube-system", deptName)
	if err != nil {
		return nil, err
	}
	baseDept := &v1.BaseDepartment{}
	err = runtimeObjectToInstanceObj(obj, baseDept)
	if err != nil {
		return nil, err
	}
	return baseDept, nil
}

func (auth *Authorization) getPermission(username string) (map[k8s.ResourceType]permission.Type, error) {
	user, err := auth.getUser(username)
	if err != nil {
		return nil, err
	}
	privilegeMap := make(map[k8s.ResourceType]permission.Type)
	for _, item := range user.Spec.Roles {
		ueRole, err := auth.roleServices.Get("kube-system", item)
		if err != nil {
			break
		}
		rolePrivilege, err := ueRole.Get("spec.privilege")
		if err != nil {
			break
		}
		rolePrivilegeMap := rolePrivilege.(map[k8s.ResourceType]permission.Type)
		for key, value := range rolePrivilegeMap {
			if value != 0 {
				privilegeMap[key] = privilegeMap[key] | value
			}
		}
	}
	return privilegeMap, nil

}

func (auth *Authorization) Config(tokenStr string) ([]byte, error) {
	return nil, nil
}

func (auth *Authorization) Auth(user *User) ([]byte, error) {
	userObj, err := auth.userServices.Get("kube-system", user.Username)
	if err != nil {
		return nil, err
	}
	password, err := userObj.Get("spec.password")
	if err != nil {
		return nil, fmt.Errorf("password not match")
	}
	//baseUser := &v1.BaseUser{}
	//runtimeObjectToInstanceObj(obj.Unstructured, baseUser)
	if password != user.Password {
		return nil, fmt.Errorf("password not match")
	}

	expireTime := time.Now().Add(time.Hour * 24).Unix()
	tokenStr, err := (&gateway.Token{}).Encode(common.MicroSaltUserHeader, user.Username, expireTime)
	if err != nil {
		return nil, err
	}
	// user AllowedNamespaces
	specDeptId, err := userObj.Get("spec.department_id")
	if err != nil {
		return nil, err
	}
	deptId := specDeptId.(string)
	deptObj, err := auth.deptServices.Get("kube-system", deptId)
	if err != nil {
		return nil, err
	}
	deptSpecNamespace, err := deptObj.Get("spec.namespace")
	if err != nil {
		return nil, err
	}
	var namespace []string
	if deptSpecNamespace != nil {
		namespace = deptSpecNamespace.([]string)
	}

	deptSpecDefaultNamespace, err := deptObj.Get("spec.default_namespace")
	if err != nil {
		return nil, err
	}
	var defaultNamespace string
	if deptSpecDefaultNamespace != nil {
		defaultNamespace = deptSpecDefaultNamespace.(string)
	}

	return []byte(
		NewUserConfig(
			user.Username,
			tokenStr,
			namespace,
			defaultNamespace,
		).String(),
	), nil
}

func runtimeObjectToInstanceObj(robj runtime.Object, targeObj interface{}) error {
	bytesData, err := json.Marshal(robj)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytesData, targeObj)
}
