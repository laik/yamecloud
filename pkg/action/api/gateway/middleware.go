package gateway

import (
	"github.com/gin-gonic/gin"
	"github.com/yametech/yamecloud/pkg/k8s"
	"github.com/yametech/yamecloud/pkg/micro/gateway"
	"github.com/yametech/yamecloud/pkg/permission"
	"github.com/yametech/yamecloud/pkg/uri"
	"net/http"
)

//func jwt() gin.HandlerFunc{
//	//return
//}

var excludeMap = map[string]string{"/user-login": "POST"}

func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		method := c.Request.Method
		//skip exclude url
		if isSkip(path, method) {
			c.Next()
			return
		}

		//get token
		token := c.Request.Header.Get("Authorization")
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

		getPermission(decodeToken.UserName)
		//check permission
		op, err := uri.NewUriParser().ParseOp(c.Request.URL.Path)
		if err != nil || op == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "error"})
			return
		}

	}

}

func getPermission(name string) map[string]permission.Type {
	return map[k8s.ResourceType]permission.Type{k8s.Pod: 255}

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
