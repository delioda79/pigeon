package twilio

import (
	"fmt"
	phttp "github.com/beatlabs/patron/trace/http"
	"github.com/taxibeat/pigeon/internal/config"
	"github.com/taxibeat/pigeon/internal/messaging"
	"github.com/taxibeat/pigeon/twilio"
	"github.com/taxibeat/pigeon/twilio/progrsms"
)

const (
	notImplemented = "not implemented"
	providerName   = "twilio"
)

// Twilio is the messaging provider for teh twilio sms service
type Twilio struct {
	cfg *config.Configuration
	cl  twilio.ProgrammableSMS
}

// Send sends a message through twilio
func (tp *Twilio) Send(m messaging.Message) (messaging.MessageResource, error) {

	mc := twilio.MessageCreate{
		To:   m.Recipient,
		Body: m.Content,
	}

	if m.ID != "" {
		mc.StatusCallback = fmt.Sprintf("%s%s/%s", tp.cfg.RestURL.Get(), tp.cfg.TwilioCallBack.Get(), m.ID)
	}

	if !m.Critical {
		mc.From = tp.cfg.TwilioNonTimeCriticalPool.Get()
	} else {
		mc.From = tp.cfg.TwilioTimeCriticalPool.Get()
	}

	rsp, err := tp.cl.Send(mc)
	if err != nil {
		return messaging.MessageResource{}, err
	}
	return messaging.MessageResource{Message: m, Status: rsp.Status, ProviderID: rsp.SID}, nil
}

// New creates a new twilio provider
func New(cfg *config.Configuration) (*Twilio, error) {

	pcl, err := phttp.New()
	if err != nil {
		return nil, err
	}

	tcl := messaging.NewHTTPClient(pcl, providerName)

	cl, _ := progrsms.New(tcl, cfg.TwilioSID.Get(), cfg.TwilioToken.Get(), "")

	return &Twilio{cfg: cfg, cl: cl}, nil
}
