package prometheus

import "github.com/prometheus/client_golang/prometheus"

var (
	RequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "path", "status"},
	)

	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Histogram of latencies for HTTP requests.",
		},
		[]string{"method", "path"},
	)

	ClientErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_client_errors_total",
			Help: "Total number of client-side errors.",
		},
		[]string{"method", "path", "status"},
	)
	ServerErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_server_errors_total",
			Help: "Total number of server-side errors.",
		},
		[]string{"method", "path", "status"},
	)
)

func init() {
	prometheus.MustRegister(RequestsTotal)
	prometheus.MustRegister(RequestDuration)
	prometheus.MustRegister(ClientErrors)
	prometheus.MustRegister(ServerErrors)
}
