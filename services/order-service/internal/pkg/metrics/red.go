package metrics

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const namespace = "course"

func init() {
	prometheus.MustRegister(httpRequestTotal)
	prometheus.MustRegister(httpRequestDurationSeconds)
}

// HTTP RED metrics
var (
	httpRequestTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "order_service_http_requests_total",
		Help:      "Total number of HTTP requests",
	}, []string{"method", "path", "code"})

	httpRequestDurationSeconds = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "order_service_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "code"})
)

func IncHttpRequestsTotal(method, path string, code int) {
	httpRequestTotal.WithLabelValues(method, path, strconv.Itoa(code)).Inc()
}

func ObserveHttpRequestDuration(method, path string, code int, elapsed time.Duration) {
	httpRequestDurationSeconds.WithLabelValues(method, path, strconv.Itoa(code)).
		Observe(elapsed.Seconds())
}
