package gateway

//
//import (
//	"fmt"
//	"github.com/gin-gonic/gin"
//	"github.com/yametech/yamecloud/common"
//	apiCommon "github.com/yametech/yamecloud/pkg/action/api/common"
//	"github.com/yametech/yamecloud/pkg/uri"
//	"net/http"
//)
//
//const (
//	IsSkip                = "isSkip"
//	AuthorizationUserName = "userName"
//	UserIdentification    = "userIdentification"
//)
//
//var uriParser = uri.NewURIParser()
//
//func IsNeedSkip(auth IAuthorization) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		path := c.Request.URL.Path
//		method := c.Request.Method
//		isNeedSkip, err := auth.IsNeedSkip(method, path)
//		if err != nil {
//			c.Abort()
//			c.JSON(http.StatusForbidden, gin.H{"error": "error"})
//			SetRequestCompletedFlag(c.Request)
//			return
//		}
//
//		if isNeedSkip {
//			c.Set(IsSkip, true)
//		}
//	}
//
//}
//
//func ValidateToken(auth IAuthorization) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		isSkip := c.GetBool(IsSkip)
//		if isSkip {
//			c.Next()
//			return
//		}
//		token := c.Request.Header.Get(common.AuthorizationHeader)
//		if token == "" {
//			apiCommon.Unauthorized(c, "")
//			SetRequestCompletedFlag(c.Request)
//			return
//		}
//		//decode token
//		decodeToken, err := auth.ValidateToken(token)
//		if err != nil {
//			apiCommon.Unauthorized(c, decodeToken)
//			SetRequestCompletedFlag(c.Request)
//			return
//		}
//		c.Set(AuthorizationUserName, decodeToken.UserName)
//		c.Request.Header.Add(AuthorizationUserName, decodeToken.UserName)
//
//	}
//}
//
//func IsAdmin(auth IAuthorization) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		isSkip := c.GetBool(IsSkip)
//		if isSkip {
//			c.Next()
//			return
//		}
//		username := c.GetString(AuthorizationUserName)
//		if username == "" {
//			apiCommon.Forbidden(c, "")
//			SetRequestCompletedFlag(c.Request)
//			return
//		}
//		isAdmin, err := auth.IsAdmin(username)
//		if err != nil {
//			apiCommon.Forbidden(c, "")
//			SetRequestCompletedFlag(c.Request)
//			return
//		}
//		if isAdmin {
//			c.Set(UserIdentification, Admin)
//			c.Request.Header.Add(UserIdentification, string(Admin))
//		}
//
//	}
//}
//
//func IsTenantOwner(auth IAuthorization) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		isSkip := c.GetBool(IsSkip)
//		if isSkip {
//			c.Next()
//			return
//		}
//		_, exists := c.Get(UserIdentification)
//		if exists {
//			c.Next()
//			return
//		}
//		username := c.GetString(AuthorizationUserName)
//		if username == "" {
//			apiCommon.Forbidden(c, "")
//			SetRequestCompletedFlag(c.Request)
//			return
//		}
//		isTenantOwner, err := auth.IsTenantOwner(username)
//		if err != nil {
//			apiCommon.Forbidden(c, "")
//			SetRequestCompletedFlag(c.Request)
//			return
//		}
//		if isTenantOwner {
//			c.Set(UserIdentification, TenantOwner)
//			c.Request.Header.Add(UserIdentification, string(TenantOwner))
//		}
//
//	}
//}
//
//func IsDepartmentOwner(auth IAuthorization) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		isSkip := c.GetBool(IsSkip)
//		if isSkip {
//			c.Next()
//			return
//		}
//		_, exists := c.Get(UserIdentification)
//		if exists {
//			c.Next()
//			return
//		}
//
//		username := c.GetString(AuthorizationUserName)
//		if username == "" {
//			apiCommon.Forbidden(c, "")
//			SetRequestCompletedFlag(c.Request)
//			return
//		}
//		isDepartmentOwner, err := auth.IsDepartmentOwner(username)
//		if err != nil {
//			apiCommon.Forbidden(c, "")
//			SetRequestCompletedFlag(c.Request)
//			return
//		}
//		if isDepartmentOwner {
//			c.Set(UserIdentification, DepartmentOwner)
//			c.Request.Header.Add(UserIdentification, string(DepartmentOwner))
//		}
//
//	}
//}
//
//func CheckNamespace(auth IAuthorization) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		isSkip := c.GetBool(IsSkip)
//		if isSkip {
//			c.Next()
//			return
//		}
//
//		username := c.GetString(AuthorizationUserName)
//		if username == "" {
//			apiCommon.Forbidden(c, "")
//			SetRequestCompletedFlag(c.Request)
//			return
//		}
//
//		isAdmin := false
//		isTenantOwner := false
//		isDepartmentOwner := false
//		userIdentification, exists := c.Get(UserIdentification)
//		if exists {
//			if userIdentification.(Identification) == Admin {
//				isAdmin = true
//			}
//			if userIdentification.(Identification) == TenantOwner {
//				isTenantOwner = true
//			}
//			if userIdentification.(Identification) == DepartmentOwner {
//				isDepartmentOwner = true
//			}
//		}
//
//		op, err := uriParser.ParseOp(c.Request.Method, c.Request.URL.Path)
//		if err != nil {
//			apiCommon.Forbidden(c, "")
//			SetRequestCompletedFlag(c.Request)
//			return
//		}
//		checkNamespace, err := auth.CheckNamespace(username, op.Namespace, isAdmin, isTenantOwner, isDepartmentOwner)
//		if err != nil {
//			apiCommon.Forbidden(c, "")
//			SetRequestCompletedFlag(c.Request)
//			return
//		}
//		if !checkNamespace {
//			apiCommon.Forbidden(c, "")
//			SetRequestCompletedFlag(c.Request)
//			return
//		}
//
//	}
//}
//
//func CheckPermission(auth IAuthorization) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		isSkip := c.GetBool(IsSkip)
//		if isSkip {
//			c.Next()
//			return
//		}
//
//		userIdentification, exists := c.Get(UserIdentification)
//		if exists && userIdentification.(Identification) != OrdinaryUser {
//			c.Next()
//			return
//		}
//		c.Set(UserIdentification, OrdinaryUser)
//
//		username := c.GetString(AuthorizationUserName)
//		if username == "" {
//			apiCommon.Forbidden(c, "")
//			SetRequestCompletedFlag(c.Request)
//			return
//		}
//		op, err := uriParser.ParseOp(c.Request.Method, c.Request.URL.Path)
//		if err != nil {
//			apiCommon.Forbidden(c, "")
//			SetRequestCompletedFlag(c.Request)
//			return
//		}
//		checkPermission, err := auth.CheckPermission(username, op)
//		fmt.Println(username, "access", c.Request.URL.Path, ", checkPermission=", checkPermission)
//		if err != nil {
//			apiCommon.Forbidden(c, "")
//			SetRequestCompletedFlag(c.Request)
//			return
//		}
//		if !checkPermission {
//			apiCommon.Forbidden(c, "")
//			SetRequestCompletedFlag(c.Request)
//			return
//		}
//
//	}
//}
//
//
//func IsWithGranted(auth IAuthorization) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		path := c.Request.URL.Path
//		// todo replace real grant path
//		if path == "/grantPath" {
//			username := c.GetString(AuthorizationUserName)
//			if username == "" {
//				apiCommon.Forbidden(c, "")
//				SetRequestCompletedFlag(c.Request)
//				return
//			}
//			isWithGranted, err := auth.IsWithGranted(username)
//			if err != nil {
//				apiCommon.Forbidden(c, "")
//				SetRequestCompletedFlag(c.Request)
//				return
//			}
//			if !isWithGranted {
//				apiCommon.Forbidden(c, "")
//				SetRequestCompletedFlag(c.Request)
//				return
//			}
//		}
//	}
//
//}
//
//func SetRequestCompletedFlag(r *http.Request) {
//	r.Header.Add(common.RequestCompletedKey, "true")
//}
