package http

import (
	"context"
	"fmt"
	"github.com/beatlabs/patron/sync"
	phttp "github.com/beatlabs/patron/sync/http"
	"github.com/taxibeat/pigeon/internal/config"
	"github.com/taxibeat/pigeon/internal/ingestion"
	"github.com/taxibeat/pigeon/internal/messaging"
	http "net/http"
)

const (
	base = "/v1/message"
)

// Ingestion represents a handler for messaging requests ingestion
type Ingestion struct {
	sdr messaging.Sender
	cfg *config.Configuration
}

func (ng *Ingestion) send(ctx context.Context, request *sync.Request) (*sync.Response, error) {
	msg := &messaging.Message{}

	err := request.Decode(msg)
	if err != nil {
		ingestion.ObserveCount("sms", "http", false, false, false, true)
		return nil, phttp.NewErrorWithCodeAndPayload(400, err.Error())
	}

	rs, err := ng.sdr.Send(*msg)
	if err != nil {
		ingestion.ObserveCount("sms", "http", false, msg.Critical, true, true)
		return nil, phttp.NewErrorWithCodeAndPayload(400, err.Error())
	}

	ingestion.ObserveCount("sms", "http", true, msg.Critical, true, true)
	return sync.NewResponse(rs), nil
}

// Routes returns the available routes
func (ng *Ingestion) Routes() []phttp.Route {
	return []phttp.Route{
		phttp.NewPostRoute(fmt.Sprintf("%s%s", base, "/send"), ng.send, true, ng.blockRoute),
	}
}

func (ng *Ingestion) blockRoute(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !ng.cfg.HTTPEnabled.Get() {
			w.WriteHeader(503)
			w.Write([]byte("Service Unavailable"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

// New returns a new Ingestion routes handler
func New(sdr messaging.Sender, cfg *config.Configuration) Ingestion {
	return Ingestion{sdr: sdr, cfg: cfg}
}
