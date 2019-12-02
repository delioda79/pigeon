package messenger

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

var (
	messageCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "pigeon",
			Subsystem: "messenger",
			Name:      "handled_messages",
			Help:      "Counts every time a message gets receiver by the messenger",
		},
		[]string{"type", "provider", "success"},
	)

	messageProcessingLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "pigeon",
			Subsystem: "messenger",
			Name:      "latency_seconds",
			Help:      "Latency of ingested messages processing",
		},
		[]string{"type", "provider"},
	)
)

func init() {
	prometheus.MustRegister(messageCounter)
	prometheus.MustRegister(messageProcessingLatency)
}

// ObserveCount is responsible to count distance matrix calls
func ObserveCount(mType, provider string, success bool) {
	messageCounter.WithLabelValues(mType, provider, strconv.FormatBool(success)).Inc()
}

// ObserveLatency is responsible to observe latency.
func ObserveLatency(mType, provider string, latency time.Duration) {
	messageProcessingLatency.WithLabelValues(mType, provider).Observe(latency.Seconds())
}
