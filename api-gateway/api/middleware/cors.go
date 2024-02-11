package middleware

import (
	"net/http"
	"strings"

	"github.com/spf13/viper"
)

func CORS(next http.Handler) http.Handler {
	accessOrigins := viper.GetStringSlice("access_origins")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			origin     = r.Header.Get("Origin")
			safeOrigin = false
		)

		for i := 0; i < len(accessOrigins); i++ {
			if strings.EqualFold(accessOrigins[i], origin) {
				safeOrigin = true
				break
			}
		}

		if safeOrigin {
			w.Header().Add("Access-Control-Allow-Credentials", "true")
			w.Header().Add("Access-Control-Allow-Origin", origin)
		} else {
			w.Header().Add("Access-Control-Allow-Credentials", "false")
			w.Header().Add("Access-Control-Allow-Origin", strings.Join(accessOrigins, ", "))
		}

		w.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Content-Length, Accept-Encoding, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		if r.Method == http.MethodOptions {
			return
		}

		next.ServeHTTP(w, r)
	})
}
