package gateway

import (
	"github.com/gin-gonic/gin"
	"github.com/yametech/yamecloud/common"
	"github.com/yametech/yamecloud/pkg/uri"
	"net/http"
)

const (
	IsSkip                = "isSkip"
	AuthorizationUserName = "userName"
	UserIdentification    = "userIdentification"
)

func IsNeedSkip(auth IAuthorization) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		method := c.Request.Method
		isNeedSkip, err := auth.IsNeedSkip(method, path)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusUnauthorized, gin.H{"error": "error"})
			return
		}

		if isNeedSkip {
			c.Set(IsSkip, true)
		}
	}

}

func ValidateToken(auth IAuthorization) gin.HandlerFunc {
	return func(c *gin.Context) {
		isSkip := c.GetBool(IsSkip)
		if isSkip {
			c.Next()
			return
		}
		token := c.Request.Header.Get(common.AuthorizationHeader)
		if token == "" {
			c.Abort()
			c.JSON(http.StatusUnauthorized, gin.H{"error": "error"})
			return
		}
		//decode token
		decodeToken, err := auth.ValidateToken(token)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusUnauthorized, gin.H{"error": "error"})
			return
		}
		c.Set(AuthorizationUserName, decodeToken.UserName)

	}
}

func IsAdmin(auth IAuthorization) gin.HandlerFunc {
	return func(c *gin.Context) {
		isSkip := c.GetBool(IsSkip)
		if isSkip {
			c.Next()
			return
		}
		username := c.GetString(AuthorizationUserName)
		if username == "" {
			c.Abort()
			c.JSON(http.StatusUnauthorized, gin.H{"error": "error"})
			return
		}
		isAdmin, err := auth.IsAdmin(username)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusUnauthorized, gin.H{"error": "error"})
			return
		}
		if isAdmin {
			c.Set(UserIdentification, Admin)
		}

	}
}

func IsTenantOwner(auth IAuthorization) gin.HandlerFunc {
	return func(c *gin.Context) {
		isSkip := c.GetBool(IsSkip)
		if isSkip {
			c.Next()
			return
		}
		_, exists := c.Get(UserIdentification)
		if exists {
			c.Next()
			return
		}
		username := c.GetString(AuthorizationUserName)
		if username == "" {
			c.Abort()
			c.JSON(http.StatusUnauthorized, gin.H{"error": "error"})
			return
		}
		isTenantOwner, err := auth.IsTenantOwner(username)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "error"})
			return
		}
		if isTenantOwner {
			c.Set(UserIdentification, TenantOwner)
		}

	}
}

func IsDepartmentOwner(auth IAuthorization) gin.HandlerFunc {
	return func(c *gin.Context) {
		isSkip := c.GetBool(IsSkip)
		if isSkip {
			c.Next()
			return
		}
		_, exists := c.Get(UserIdentification)
		if exists {
			c.Next()
			return
		}

		username := c.GetString(AuthorizationUserName)
		if username == "" {
			c.Abort()
			c.JSON(http.StatusUnauthorized, gin.H{"error": "error"})
			return
		}
		isDepartmentOwner, err := auth.IsDepartmentOwner(username)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusUnauthorized, gin.H{"error": "error"})
			return
		}
		if isDepartmentOwner {
			c.Set(UserIdentification, DepartmentOwner)
		}

	}
}

func CheckUser(auth IAuthorization) gin.HandlerFunc {
	return func(c *gin.Context) {
		isSkip := c.GetBool(IsSkip)
		if isSkip {
			c.Next()
			return
		}

		_, exists := c.Get(UserIdentification)
		if exists {
			c.Next()
			return
		}
		c.Set(UserIdentification, OrdinaryUser)

		username := c.GetString(AuthorizationUserName)
		if username == "" {
			c.Abort()
			c.JSON(http.StatusUnauthorized, gin.H{"error": "error"})
			return
		}
		op, err := uri.NewUriParser().ParseOp(c.Request.URL.Path)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusUnauthorized, gin.H{"error": "error"})
			return
		}
		checkNamespace, err := auth.CheckNamespace(username, op.Namespace)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusUnauthorized, gin.H{"error": "error"})
			return
		}
		if !checkNamespace {
			c.Abort()
			c.JSON(http.StatusUnauthorized, gin.H{"error": "error"})
			return
		}
		checkPermission, err := auth.CheckPermission(username, op)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusUnauthorized, gin.H{"error": "error"})
			return
		}
		if !checkPermission {
			c.Abort()
			c.JSON(http.StatusUnauthorized, gin.H{"error": "error"})
			return
		}

	}
}

func IsWithGranted(auth IAuthorization) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		// todo replace real grant path
		if path == "/grantPath" {
			username := c.GetString(AuthorizationUserName)
			if username == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "error"})
				return
			}
			isWithGranted, err := auth.IsWithGranted(username)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "error"})
				return
			}
			if !isWithGranted {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "error"})
				return
			}
		}
	}

}