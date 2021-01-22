package gateway

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yametech/yamecloud/common"
	"github.com/yametech/yamecloud/pkg/k8s"
	"github.com/yametech/yamecloud/pkg/micro/gateway"
	"github.com/yametech/yamecloud/pkg/permission"
	"github.com/yametech/yamecloud/pkg/uri"
)

//func jwt() gin.HandlerFunc{
//	//return
//}
var excludeMap = map[string]string{"/user-login": "POST"}

func Authorize(auth *Authorization) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		method := c.Request.Method
		//skip exclude url
		if isSkip(path, method) {
			c.Next()
			return
		}

		//get token
		token := c.Request.Header.Get(common.AuthorizationHeader)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "error"})
			return
		}
		//decode token
		decodeToken, err := (&gateway.Token{}).Decode(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "error"})
			return
		}

		//get user permission
		//privilegeMap, err := auth.getPermission(decodeToken.UserName)
		//if err != nil {
		//	c.JSON(http.StatusUnauthorized, gin.H{"error": "error"})
		//	return
		//}

		privilegeMap := map[k8s.ResourceType]permission.Type{k8s.Pod: 255}

		//check permission
		op, err := uri.NewURIParser().ParseOp(c.Request.URL.Path)
		if err != nil || op == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "error"})
			return
		}
		if !checkPermission(&privilegeMap, op) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "error"})
			return
		}
		//check if namespace allow access when op.Namespace not nil
		if op.Namespace != "" {
			isAllow, err := auth.allowNamespaceAccess(decodeToken.UserName, op.Namespace)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "error"})
				return
			}

			if !isAllow {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "namespace not allow access"})
				return
			}
		}
	}

}

func checkPermission(permissionMap *map[k8s.ResourceType]permission.Type, op *uri.URI) bool {
	//permissionValue := (*permissionMap)[op.Resource]
	//if permissionValue&1<<op.Type != 0 {
	//	return true
	//}
	return true
}

func isSkip(path string, method string) bool {
	requestMethod, exists := excludeMap[path]
	if !exists {
		return false
	}
	if method == requestMethod {
		return true
	}
	return false
}
