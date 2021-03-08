package gateway

import (
	"fmt"
	"github.com/micro/micro/v2/plugin"
	"github.com/yametech/yamecloud/common"
	apiCommon "github.com/yametech/yamecloud/pkg/action/api/common"
	"github.com/yametech/yamecloud/pkg/uri"
	"net/http"
	"strings"
)

const (
	IsSkip                = "isSkip"
	AuthorizationUserName = "userName"
	UserIdentification    = "userIdentification"
	UnauthorizedMessage   = "user unauthorized"
	ForbiddenMessage      = "not allow to access"
)

var uriParser = uri.NewURIParser()

func SkipFilter(auth IAuthorization) plugin.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			method := r.Method
			isNeedSkip, err := auth.IsNeedSkip(method, path)
			if err != nil {
				apiCommon.ResponseJSONFromError(w, http.StatusForbidden, nil, err)
				return
			}

			if isNeedSkip {
				r.Header.Set(IsSkip, "true")
			}
			h.ServeHTTP(w, r)
		})
	}
}

func ValidateTokenFilter(auth IAuthorization) plugin.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			isSkip := r.Header.Get(IsSkip)
			if isSkip == "true" {
				h.ServeHTTP(w, r)
				return
			}
			token := r.Header.Get(common.AuthorizationHeader)
			if token == "" {
				//apiCommon.Unauthorized(c, "")
				//SetRequestCompletedFlag(c.Request)
				//w.WriteHeader(http.StatusUnauthorized)
				apiCommon.ResponseJSON(w, http.StatusUnauthorized, nil, UnauthorizedMessage)
				return
			}
			//decode token
			decodeToken, err := auth.ValidateToken(token)
			if err != nil {
				//apiCommon.Unauthorized(c, decodeToken)
				//SetRequestCompletedFlag(c.Request)
				apiCommon.ResponseJSONFromError(w, http.StatusUnauthorized, token, err)
				return
			}
			r.Header.Add(AuthorizationUserName, decodeToken.UserName)
			h.ServeHTTP(w, r)
		})
	}
}

func IdentificationFilter(auth IAuthorization) plugin.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			isSkip := r.Header.Get(IsSkip)
			if isSkip == "true" {
				h.ServeHTTP(w, r)
				return
			}
			username := r.Header.Get(AuthorizationUserName)
			if username == "" {
				//apiCommon.Forbidden(c, "")
				//SetRequestCompletedFlag(c.Request)
				apiCommon.ResponseJSON(w, http.StatusForbidden, nil, ForbiddenMessage)
				return
			}
			isAdmin, err := auth.IsAdmin(username)
			if err != nil {
				apiCommon.ResponseJSONFromError(w, http.StatusForbidden, username, err)
				return
			}
			if isAdmin {
				r.Header.Add(UserIdentification, string(Admin))
				h.ServeHTTP(w, r)
				return
			}
			isTenantOwner, err := auth.IsTenantOwner(username)
			if err != nil {
				apiCommon.ResponseJSONFromError(w, http.StatusForbidden, username, err)
				return
			}
			if isTenantOwner {
				r.Header.Add(UserIdentification, string(TenantOwner))
				h.ServeHTTP(w, r)
				return
			}

			isDepartmentOwner, err := auth.IsDepartmentOwner(username)
			if err != nil {
				apiCommon.ResponseJSONFromError(w, http.StatusForbidden, username, err)
				return
			}
			if isDepartmentOwner {
				r.Header.Add(UserIdentification, string(DepartmentOwner))
				h.ServeHTTP(w, r)
				return
			}
			r.Header.Add(UserIdentification, string(OrdinaryUser))
			h.ServeHTTP(w, r)
		})
	}
}

func NamespaceFilter(auth IAuthorization) plugin.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			isSkip := r.Header.Get(IsSkip)
			if isSkip == "true" {
				h.ServeHTTP(w, r)
				return
			}

			username := r.Header.Get(AuthorizationUserName)
			if username == "" {
				apiCommon.ResponseJSON(w, http.StatusForbidden, nil, ForbiddenMessage)
				return
			}

			isAdmin := false
			isTenantOwner := false
			isDepartmentOwner := false
			userIdentification := r.Header.Get(UserIdentification)
			if userIdentification == string(Admin) {
				isAdmin = true
			}
			if userIdentification == string(TenantOwner) {
				isTenantOwner = true
			}
			if userIdentification == string(DepartmentOwner) {
				isDepartmentOwner = true
			}

			op, err := uriParser.ParseOp(r.Method, r.URL.Path)
			if err != nil {
				fmt.Printf("username %s access url %s error: %v\n", username, r.URL.Path, err)
				apiCommon.ResponseJSONFromError(w, http.StatusForbidden, username, err)
				return
			}
			checkNamespace, err := auth.CheckNamespace(username, op.Namespace, isAdmin, isTenantOwner, isDepartmentOwner)
			if err != nil {
				apiCommon.ResponseJSONFromError(w, http.StatusForbidden, username, err)
				return
			}
			if !checkNamespace {
				apiCommon.ResponseJSON(w, http.StatusForbidden, username, ForbiddenMessage)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}

func PermissionFilter(auth IAuthorization) plugin.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			isSkip := r.Header.Get(IsSkip)
			if isSkip == "true" {
				h.ServeHTTP(w, r)
				return
			}

			username := r.Header.Get(AuthorizationUserName)
			if username == "" {
				apiCommon.ResponseJSON(w, http.StatusForbidden, nil, ForbiddenMessage)
				return
			}

			userIdentification := r.Header.Get(UserIdentification)
			if userIdentification == string(Admin) ||
				userIdentification == string(TenantOwner) ||
				userIdentification == string(DepartmentOwner) {
				h.ServeHTTP(w, r)
				return
			}

			op, err := uriParser.ParseOp(r.Method, r.URL.Path)
			if err != nil {
				fmt.Printf("username %s access url %s error: %v\n", username, r.URL.Path, err)
				apiCommon.ResponseJSONFromError(w, http.StatusForbidden, username, err)
				return
			}
			checkPermission, err := auth.CheckPermission(username, op)
			fmt.Println(username, "access", r.URL.Path, ", checkPermission=", checkPermission)
			if err != nil {
				apiCommon.ResponseJSONFromError(w, http.StatusForbidden, username, err)
				return
			}
			if !checkPermission {
				apiCommon.ResponseJSON(w, http.StatusForbidden, username, ForbiddenMessage)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}

func GrantCheckFilter(auth IAuthorization) plugin.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			// todo replace real grant path
			if path == "/grantPath" {
				username := r.Header.Get(AuthorizationUserName)
				if username == "" {
					apiCommon.ResponseJSON(w, http.StatusForbidden, nil, ForbiddenMessage)
					return
				}
				isWithGranted, err := auth.IsWithGranted(username)
				if err != nil {
					apiCommon.ResponseJSONFromError(w, http.StatusForbidden, username, err)
					return
				}
				if !isWithGranted {
					apiCommon.ResponseJSON(w, http.StatusForbidden, username, ForbiddenMessage)
					return
				}
			}
			h.ServeHTTP(w, r)
		})
	}
}

func ServerFilter(self http.Handler) plugin.Handler {
	return func(redirect http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if strings.EqualFold(r.URL.Path, "/user-login") {
				self.ServeHTTP(w, r)
				return
			}

			// user login return CaaS config
			if r.Method == http.MethodGet && r.URL.Path == "/config" {
				self.ServeHTTP(w, r)
				return
			}
			redirect.ServeHTTP(w, r)
		})
	}
}
