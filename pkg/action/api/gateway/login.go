package gateway

import (
	"fmt"
	"time"

	"github.com/yametech/yamecloud/common"
	"github.com/yametech/yamecloud/pkg/action/service"
	"github.com/yametech/yamecloud/pkg/action/service/tenant"
	"github.com/yametech/yamecloud/pkg/micro/gateway"
	"github.com/yametech/yamecloud/pkg/utils"
)

const PASSWORD_NOT_MATH = "password not match"

type LoginHandle struct {
	userServices     *tenant.BaseUser
	roleUserServices *tenant.BaseRoleUser
	roleServices     *tenant.BaseRole
	deptServices     *tenant.BaseDepartment
	tenantServices   *tenant.BaseTenant
	auth             *gateway.Authorization
}

func (lh *LoginHandle) getUserConfig(user *User, tokenStr string) (*userConfig, error) {
	//check whether or not an admin
	isAdmin, _ := lh.auth.IsAdmin(user.Username)

	//check whether or not a tenant owner
	isTenantOwner, err := lh.auth.IsTenantOwner(user.Username)
	if err != nil {
		return nil, err
	}
	//check whether or not a department owner
	isDepartmentOwner, err := lh.auth.IsDepartmentOwner(user.Username)
	if err != nil {
		return nil, err
	}

	allowNamespaces, err := lh.auth.AllowNamespaces(user.Username, isAdmin, isTenantOwner, isDepartmentOwner)
	if err != nil {
		return nil, err
	}

	return NewUserConfig(
		user.Username,
		tokenStr,
		allowNamespaces,
		isAdmin,
		isTenantOwner,
	), nil
}

func (lh *LoginHandle) Auth(user *User) (*userConfig, error) {
	userObj, err := lh.userServices.Get("", user.Username)
	if err != nil {
		return nil, err
	}

	password, err := userObj.Get("spec.password")
	if err != nil {
		return nil, fmt.Errorf(PASSWORD_NOT_MATH)
	}

	if utils.Sha1(user.Password) != password.(string) {
		return nil, fmt.Errorf(PASSWORD_NOT_MATH)
	}

	expireTime := time.Now().Add(time.Hour * 24).Unix()
	tokenStr, err := (&gateway.Token{}).Encode(common.MicroSaltUserHeader, user.Username, expireTime)
	if err != nil {
		return nil, err
	}

	return lh.getUserConfig(user, tokenStr)
}

func NewLoginHandle(svcInterface service.Interface) *LoginHandle {
	lh := &LoginHandle{
		userServices:     tenant.NewBaseUser(svcInterface),
		roleUserServices: tenant.NewBaseRoleUser(svcInterface),
		roleServices:     tenant.NewBaseRole(svcInterface),
		deptServices:     tenant.NewBaseDepartment(svcInterface),
		tenantServices:   tenant.NewBaseTenant(svcInterface),
		auth:             gateway.NewAuthorization(svcInterface),
	}
	return lh
}
