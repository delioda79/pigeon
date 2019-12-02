package messenger

import (
	"github.com/taxibeat/pigeon/internal/config"
	"github.com/taxibeat/pigeon/internal/messaging"
	"github.com/taxibeat/pigeon/internal/messaging/twilio"
)

type settings struct {
	twilio messaging.Sender
}

func newDefaultSettings(cfg *config.Configuration) (*settings, error) {
	tw, err := twilio.New(cfg)
	if err != nil {
		return nil, err
	}

	return &settings{
		twilio: tw,
	}, nil
}
