package loadshedding

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

const namespace = "load_shedding"

var (
	once sync.Once

	estimatedLimit = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "estimated_limit",
		Help:      "Total tokens allowed by selected limit algorithm",
	}, []string{"method"})
)

func init() {
	once.Do(func() {
		prometheus.MustRegister(estimatedLimit)
	})
}

func estimateLimit(handler string, limit int) {
	estimatedLimit.WithLabelValues(handler).Set(float64(limit))
}
