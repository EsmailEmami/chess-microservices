package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/esmailemami/chess/api-gateway/api/prometheus"
	"github.com/esmailemami/chess/api-gateway/api/util"
	"github.com/esmailemami/chess/shared/logging"
)

func ExecuteDuration(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ignore for websocket requests
		if util.IsWebSocketRequest(r) {
			next.ServeHTTP(w, r)
		} else {
			startTime := time.Now()
			next.ServeHTTP(w, r)
			duration := time.Since(startTime)

			if duration.Seconds() > 60 {
				logging.Warn(fmt.Sprintf("Request [%s] took %s to execute", r.RequestURI, duration))
			}

			// monitor the executation
			prometheus.RequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration.Seconds())
		}
	})
}
