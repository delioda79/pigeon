package messenger

import (
	"fmt"
	"github.com/beatlabs/patron/errors"
	"github.com/taxibeat/pigeon/internal/config"
	"github.com/taxibeat/pigeon/internal/messaging"
	"time"
)

const (
	unknownProvider = "unknownProvider"
	twilioProvider  = "twilio"
)

// Service is a messenger which finds the correct provider and uses it to send messages
type Service struct {
	cfg     *config.Configuration
	senders *settings
}

// Send uses teh correct provider to send messages
func (s Service) Send(m messaging.Message) (rs messaging.MessageResource, e error) {
	var provider string

	start := time.Now()

	switch m.Type {
	case messaging.TimeCriticalSMS:
		provider = twilioProvider
		rs, e = s.senders.twilio.Send(m)
	default:
		ObserveCount(unknownProvider, unknownProvider, false)
		return messaging.MessageResource{}, errors.New(fmt.Sprintf("Unknown type %s", m.Type))
	}

	ObserveLatency(m.Type, provider, time.Since(start))
	ObserveCount(m.Type, provider, true)

	return rs, e
}

// New instantiates a new messenger
func New(cfg *config.Configuration) (*Service, error) {

	senders, err := newDefaultSettings(cfg)
	if err != nil {
		return nil, err
	}
	mng := &Service{
		cfg:     cfg,
		senders: senders,
	}

	return mng, nil
}
