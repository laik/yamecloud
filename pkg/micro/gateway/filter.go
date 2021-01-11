package gateway

import (
	"github.com/micro/micro/v2/plugin"
	"net/http"
	"strings"
)

var _token = &Token{}

func filter(self http.Handler) plugin.Handler {
	return func(redirect http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.EqualFold(r.URL.Path, "/user-login") {
				self.ServeHTTP(w, r)
				return
			}

			// Ignore uri
			if strings.Contains(r.URL.Path, "/workload/shell/pod") ||
				strings.Contains(r.URL.Path, "/webhook") {
				redirect.ServeHTTP(w, r)
				return
			}

			//tokenHeader := r.Header.Get("Authorization")
			//userToken, e := _token.Decode(tokenHeader)
			//if e != nil {
			//	w.WriteHeader(http.StatusUnauthorized)
			//	return
			//}
			//
			//r.Header.Set(common.HttpRequestUserHeaderKey, userToken.UserName)

			// user login return CaaS config
			if r.Method == http.MethodGet && r.URL.Path == "/config" {
				self.ServeHTTP(w, r)
				return
			}

			// all uri access authorization certainty
			self.ServeHTTP(w, r)

			// other uri redirect to backend service
			redirect.ServeHTTP(w, r)
		})
	}
}
