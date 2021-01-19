package gateway

import (
	"github.com/yametech/yamecloud/pkg/action/service"
)

type LoginHandle struct {
	auth *Authorization
}

func NewLoginHandle(svcInterface service.Interface) *LoginHandle {
	lh := &LoginHandle{
		auth: NewAuthorization(svcInterface),
	}
	return lh
}
