package messenger

import (
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
	var timeCritical bool

	provider = twilioProvider

	if m.Critical {
		timeCritical = true
	}

	start := time.Now()
	rs, e = s.senders.twilio.Send(m)

	ObserveLatency("sms", provider, timeCritical, time.Since(start))
	ObserveCount("sms", provider, true, timeCritical)

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
