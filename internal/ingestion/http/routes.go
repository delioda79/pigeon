package http

import (
	"context"
	"fmt"
	"github.com/beatlabs/patron/sync"
	"github.com/beatlabs/patron/sync/http"
	"github.com/taxibeat/pigeon/internal/messaging"
)

const (
	base = "/message"
)

// Ingestion represents a handler for messaging requests ingestion
type Ingestion struct {
	sdr messaging.Sender
}

func (ng *Ingestion) send(ctx context.Context, request *sync.Request) (*sync.Response, error) {
	msg := &messaging.Message{}

	err := request.Decode(msg)
	if err != nil {
		return nil, http.NewErrorWithCodeAndPayload(400, err)
	}

	rs, err := ng.sdr.Send(*msg)
	if err != nil {
		return nil, http.NewErrorWithCodeAndPayload(400, err)
	}

	return sync.NewResponse(rs), nil
}

// Routes returns the available routes
func (ng *Ingestion) Routes() []http.Route {
	return []http.Route{
		http.NewPostRoute(fmt.Sprintf("%s%s", base, "/send"), ng.send, true),
	}
}

// New returns a new Ingestion routes handler
func New(sdr messaging.Sender) Ingestion {
	return Ingestion{sdr}
}