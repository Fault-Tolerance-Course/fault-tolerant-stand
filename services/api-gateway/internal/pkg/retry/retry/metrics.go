package retry

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

const namespace = "homework"

var (
	once sync.Once

	attemptsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "grpc_client_attempts_total",
			Help:      "Total number of attempts performed",
		},
		[]string{"method"},
	)

	retriesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "grpc_client_retries_total",
			Help:      "Total number of retries performed",
		},
		[]string{"method"},
	)

	throttledTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "grpc_client_throttled_total",
			Help:      "Total number of retries stopped by throttler",
		},
		[]string{"method"},
	)

	successTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "grpc_client_success_total",
			Help:      "Total number of successful RPCs",
		},
		[]string{"method"},
	)
)

func init() {
	once.Do(func() {
		prometheus.MustRegister(retriesTotal)
		prometheus.MustRegister(throttledTotal)
		prometheus.MustRegister(successTotal)
		prometheus.MustRegister(attemptsTotal)
	})
}
