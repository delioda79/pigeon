package twilio

import (
	"fmt"
	"github.com/beatlabs/patron/errors"
	phttp "github.com/beatlabs/patron/trace/http"
	"github.com/taxibeat/pigeon/internal/config"
	"github.com/taxibeat/pigeon/internal/messaging"
	"github.com/taxibeat/pigeon/twilio"
	"github.com/taxibeat/pigeon/twilio/progrsms"
)

const (
	unknownType = "unknown message type"
)

// Twilio is the messaging provider for teh twilio sms service
type Twilio struct {
	cfg *config.Configuration
	cl  twilio.ProgrammableSMS
}

// Send sends a message through twilio
func (tp *Twilio) Send(m messaging.Message) error {

	mc := twilio.MessageCreate{
		To:             m.Receiver,
		Body:           m.Content,
		StatusCallback: fmt.Sprintf("%s%s/%s", tp.cfg.RestURL.Get(), tp.cfg.TwilioCallBack.Get(), m.ID),
	}

	switch m.Type {
	case messaging.TimeCriticalSMS:
		mc.From = tp.cfg.TwilioTimeCriticalPool.Get()
		_, err := tp.cl.Send(mc)
		if err != nil {
			return err
		}
		return nil
	}

	return errors.New(unknownType)
}

// New creates a new twilio provider
func New(cfg *config.Configuration) (*Twilio, error) {

	pcl, err := phttp.New()
	if err != nil {
		return nil, err
	}

	tcl := messaging.NewHTTPClient(pcl)

	cl, _ := progrsms.New(tcl, cfg.TwilioSID.Get(), cfg.TwilioToken.Get(), "")

	return &Twilio{cfg: cfg, cl: cl}, nil
}
