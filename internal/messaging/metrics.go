package messaging

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

var (
	messageCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "pigeon",
			Subsystem: "messaging",
			Name:      "http_requests",
			Help:      "Counts every time a request is performed",
		},
		[]string{"provider", "url", "method", "status_code"},
	)

	messageProcessingLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "pigeon",
			Subsystem: "messaging",
			Name:      "latency_seconds",
			Help:      "Latency of requests",
		},
		[]string{"provider", "url", "method", "status_code"},
	)
)

func init() {
	prometheus.MustRegister(messageCounter)
	prometheus.MustRegister(messageProcessingLatency)
}

// ObserveCount is responsible to count distance matrix calls
func ObserveCount(provider, url, method string, statusCode int) {
	messageCounter.WithLabelValues(provider, url, method, strconv.FormatInt(int64(statusCode), 10)).Inc()
}

// ObserveLatency is responsible to observe latency.
func ObserveLatency(provider, url, method string, statusCode int, latency time.Duration) {
	messageProcessingLatency.WithLabelValues(provider, url, method, strconv.FormatInt(int64(statusCode), 10)).Observe(latency.Seconds())
}
