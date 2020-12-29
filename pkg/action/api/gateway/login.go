package gateway

import (
	"encoding/json"
	"fmt"
	"github.com/yametech/yamecloud/common"
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/action/service/tenant"
	"github.com/yametech/yamecloud/pkg/micro/gateway"
	"k8s.io/apimachinery/pkg/runtime"
	"time"
)

type LoginHandle struct {
	userService       *tenant.BaseUser
	roleUserService   *tenant.BaseRoleUser
	departmentService *tenant.BaseDepartment
	tenant            *tenant.BaseTenant
}

func NewLoginHandle(svcInterface service.Interface) *LoginHandle {
	auth := &LoginHandle{
		userService:       tenant.NewBaseUser(svcInterface),
		roleUserService:   tenant.NewBaseRoleUser(svcInterface),
		departmentService: tenant.NewBaseDepartment(svcInterface),
		tenant:            tenant.NewBaseTenant(svcInterface),
	}
	return auth
}

func (lh *LoginHandle) Auth(user *User) ([]byte, error) {
	userObj, err := lh.userService.Get("kube-system", user.Username)
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
	deptObj, err := lh.departmentService.Get("kube-system", deptId)
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
