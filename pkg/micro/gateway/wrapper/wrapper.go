package wrapper

import (
	"github.com/gin-gonic/gin"
	"github.com/yametech/yamecloud/common"
	"github.com/yametech/yamecloud/pkg/micro/gateway/token"
	"net/http"
)

//func(w http.ResponseWriter, r *http.Request) {
//	//var tokenHeader string
//	if r.URL.Path == "/user-login" {
//		loginHandler.ServeHTTP(w, r)
//		return
//	}
//
//	if strings.Contains(r.URL.Path, "/workload/shell/pod") {
//		h.ServeHTTP(w, r)
//		return
//	}
//
//	if strings.Contains(r.URL.Path, "/webhook") {
//		h.ServeHTTP(w, r)
//		return
//	}
//
//	tokenHeader := r.Header.Get("Authorization")
//	userFromToken, e := token.Decode(tokenHeader)
//	if e != nil {
//		w.WriteHeader(http.StatusUnauthorized)
//		return
//	}
//
//	r.Header.Set(common.HttpRequestUserHeaderKey, userFromToken.UserName)
//	//Config
//	if r.Method == http.MethodGet && r.URL.Path == "/config" {
//		loginHandler.ServeHTTP(w, r)
//		return
//	}
//	h.ServeHTTP(w, r)
//}

// JWTAuthWrapper

func JWTAuthWrapperMiddleware(g *gin.Context) {
	userToken, err := (&token.Token{}).Decode(g.GetHeader(common.AuthorizationHeader))
	if err != nil {
		g.JSON(http.StatusUnauthorized, nil)
		return
	}
	g.Header(common.HttpRequestUserHeaderKey, userToken.UserName)
	g.Next()
}

// JWTAuthWrapper
func CorsWrapper(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, PATCH, DELETE")
		w.Header().Add("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers,X-Access-Token,XKey,Authorization")
		w.Header().Add("Access-Control-Allow-Origin", "*")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
		}
		handler.ServeHTTP(w, r)
	})
}
