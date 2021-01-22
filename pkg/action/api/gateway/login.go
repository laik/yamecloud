package gateway

import (
	"fmt"
	"github.com/yametech/yamecloud/common"
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/action/service/tenant"
	"github.com/yametech/yamecloud/pkg/micro/gateway"
	"time"
)

type LoginHandle struct {
	userServices     *tenant.BaseUser
	roleUserServices *tenant.BaseRoleUser
	roleServices     *tenant.BaseRole
	deptServices     *tenant.BaseDepartment
	tenantServices   *tenant.BaseTenant
}

func (lh *LoginHandle) Auth(user *User) ([]byte, error) {
	userObj, err := lh.userServices.Get("kube-system", user.Username)
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
	fmt.Println(tokenStr)
	// user AllowedNamespaces
	specDeptId, err := userObj.Get("spec.department_id")
	if err != nil {
		return nil, err
	}
	deptId := specDeptId.(string)
	deptObj, err := lh.deptServices.Get("kube-system", deptId)
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

func NewLoginHandle(svcInterface service.Interface) *LoginHandle {
	lh := &LoginHandle{
		userServices:     tenant.NewBaseUser(svcInterface),
		roleUserServices: tenant.NewBaseRoleUser(svcInterface),
		roleServices:     tenant.NewBaseRole(svcInterface),
		deptServices:     tenant.NewBaseDepartment(svcInterface),
		tenantServices:   tenant.NewBaseTenant(svcInterface),
	}
	return lh
}
