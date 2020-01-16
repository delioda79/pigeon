package http

import (
	"context"
	"errors"
	"github.com/taxibeat/pigeon/internal/config/sanitarium"

	"github.com/beatlabs/patron/sync"

	"github.com/beatlabs/patron/sync/http"
	"github.com/taxibeat/pigeon/internal/config"
)

// HTTP is a controller for sanitarium routes
type HTTP struct {
	cfg *config.Configuration
}

// New returns a new HTTP router
func New(cfg *config.Configuration) (*HTTP, error) {
	if cfg == nil {
		return nil, errors.New("configuration cannot be nil")
	}
	return &HTTP{cfg: cfg}, nil
}

// GetRoutes returns teh sanitarium routes
func (h HTTP) GetRoutes() []http.Route {
	var rr []http.Route
	rr = append(rr, http.NewGetRoute("/config", h.getConfig, false))
	return rr
}

func (h HTTP) getConfig(ctx context.Context, _ *sync.Request) (*sync.Response, error) {
	schema := h.getConfigSchema()

	return sync.NewResponse(schema), nil
}

func (h HTTP) getConfigSchema() sanitarium.Schema {
	return sanitarium.Schema{
		Name:        "pigeon",
		Icon:        "",
		Description: "Notification service",
		Static:      []sanitarium.SchemaItem{},
		Dynamic: []sanitarium.SchemaItem{
			{
				ID:          "twilio-time-critical-pool",
				Type:        "text",
				Label:       "TwilioTimeCriticalPool",
				Description: "The time critical pool for twilio",
				Value:       h.cfg.TwilioTimeCriticalPool.Get(),
			},
			{
				ID:          "twilio-time-non-critical-pool",
				Type:        "text",
				Label:       "TwilioTimeNonCriticalPool",
				Description: "The time non critical pool for twilio",
				Value:       h.cfg.TwilioNonTimeCriticalPool.Get(),
			},
			{
				ID:          "http-enabled",
				Type:        "switch",
				Label:       "HTTPEnabled",
				Description: "Enable HTTP ingestion",
				Value:       h.cfg.HTTPEnabled.Get(),
			},
			{
				ID:          "kafka-enabled",
				Type:        "switch",
				Label:       "KafkaEnabled",
				Description: "Enable Kafka ingestion",
				Value:       h.cfg.KafkaConsumerEnabled.Get(),
			},
		},
	}
}
