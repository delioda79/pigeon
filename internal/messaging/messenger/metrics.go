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
		[]string{"type", "provider", "success", "time_critical"},
	)

	messageProcessingLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "pigeon",
			Subsystem: "messenger",
			Name:      "latency_seconds",
			Help:      "Latency of ingested messages processing",
		},
		[]string{"type", "provider", "success", "time_critical"},
	)
)

func init() {
	prometheus.MustRegister(messageCounter)
	prometheus.MustRegister(messageProcessingLatency)
}

// ObserveCount is responsible to count distance matrix calls
func ObserveCount(mType, provider string, success, timeCritical bool) {
	messageCounter.WithLabelValues(mType, provider, strconv.FormatBool(success), strconv.FormatBool(timeCritical)).Inc()
}

// ObserveLatency is responsible to observe latency.
func ObserveLatency(mType, provider string, success, timeCritical bool, latency time.Duration) {
	messageProcessingLatency.WithLabelValues(mType, provider, strconv.FormatBool(success), strconv.FormatBool(timeCritical)).Observe(latency.Seconds())
}
