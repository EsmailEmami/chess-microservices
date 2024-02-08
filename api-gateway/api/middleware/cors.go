package middleware

import (
	"net/http"
	"strings"
)

func CORS(next http.Handler) http.Handler {
	accessOrigins := []string{}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			origin     = r.Header.Get("Origin")
			safeOrigin = false
		)

		if strings.HasPrefix(origin, "http://localhost") || strings.HasPrefix(origin, "http://127.0.0.1") ||
			strings.HasPrefix(origin, "https://localhost") || strings.HasPrefix(origin, "https://127.0.0.1") {
			safeOrigin = true
		} else {
			for i := 0; i < len(accessOrigins); i++ {
				if strings.EqualFold(accessOrigins[i], origin) {
					safeOrigin = true
					break
				}
			}
		}

		if !safeOrigin {
			origin = ""
		}

		w.Header().Add("Access-Control-Allow-Origin", origin)
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Content-Length, Accept-Encoding, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		if r.Method == http.MethodOptions {
			return
		}

		next.ServeHTTP(w, r)
	})
}
