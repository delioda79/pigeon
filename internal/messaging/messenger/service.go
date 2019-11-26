package messenger

import (
	"fmt"
	"github.com/beatlabs/patron/errors"
	"github.com/taxibeat/pigeon/internal/config"
	"github.com/taxibeat/pigeon/internal/messaging"
	"time"
)

const (
	twilio  = "twilio"
	unknown = "unknown"
)

type Service struct {
	cfg     *config.Configuration
	senders map[string]messaging.Sender
}

func (s Service) Send(m messaging.Message) error {
	var e error
	var provider string

	start := time.Now()

	switch m.Type {
	case messaging.Necessary:
		provider = s.cfg.NecessarySMS.Get()
		snd, ok := s.senders[s.cfg.NecessarySMS.Get()]
		if !ok {
			ObserveCount(messaging.Necessary, provider, false)
			return errors.New(fmt.Sprintf("Unknown provider %s", s.cfg.NecessarySMS.Get()))
		}
		e = snd.Send(m)
	default:
		ObserveCount(unknown, unknown, false)
		return errors.New(fmt.Sprintf("Unknown type %s", m.Type))
	}

	ObserveLatency(m.Type, provider, time.Since(start))
	ObserveCount(m.Type, provider, true)

	return e
}

func New(cfg *config.Configuration) *Service {

	senders := map[string]messaging.Sender{}

	mng := &Service{
		cfg:     cfg,
		senders: senders,
	}

	return mng
}
