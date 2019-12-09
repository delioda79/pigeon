package ingestion

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
)

var (
	messageCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "pigeon",
			Subsystem: "ingestion",
			Name:      "message_request",
			Help:      "Counts every time a message gets receiver by the messenger",
		},
		[]string{"type", "provider", "success", "time_critical", "valid_message", "valid_criticality"},
	)
)

func init() {
	prometheus.MustRegister(messageCounter)
}

// ObserveCount is responsible to count distance matrix calls
func ObserveCount(mType, provider string, success, timeCritical, validMessage, validCriticality bool) {
	messageCounter.WithLabelValues(mType, provider, strconv.FormatBool(success), strconv.FormatBool(timeCritical), strconv.FormatBool(validMessage), strconv.FormatBool(validCriticality)).Inc()
}
