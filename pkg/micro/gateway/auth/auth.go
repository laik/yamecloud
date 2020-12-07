package auth

import (
	"github.com/micro/micro/plugin"
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
func JWTAuthWrapper(handler http.Handler) plugin.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(handler.ServeHTTP)
	}
}
