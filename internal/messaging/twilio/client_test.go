package twilio

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/taxibeat/pigeon/internal/config"
	"github.com/taxibeat/pigeon/internal/messaging"
	"github.com/taxibeat/pigeon/twilio"
	"github.com/taxibeat/pigeon/twilio/twiliofakes"
	"testing"
)

func TestSend(t *testing.T) {

	psms := &twiliofakes.FakeProgrammableSMS{}

	cfg := &config.Configuration{}
	cfg.TwilioToken.Set("token")
	cfg.TwilioSID.Set("TSID")
	cfg.RestURL.Set("REST")
	cfg.TwilioCallBack.Set("/path")
	cfg.TwilioTimeCriticalPool.Set("pool1")

	tw := Twilio{cl: psms, cfg: cfg}

	msg := messaging.Message{
		ID:       "a1234",
		Content:  "Hello there",
		Receiver: "receiver1",
		Type:     messaging.TimeCriticalSMS,
	}

	err := tw.Send(msg)
	assert.Nil(t, err)

	mc := psms.SendArgsForCall(0)

	assert.Equal(t, mc.Body, msg.Content)
	assert.Equal(t, mc.From, cfg.TwilioTimeCriticalPool.Get())
	assert.Equal(t, mc.To, msg.Receiver)
	assert.Equal(t, mc.StatusCallback, fmt.Sprintf("REST/path/%s", msg.ID))
}

func TestError(t *testing.T) {
	psms := &twiliofakes.FakeProgrammableSMS{}
	psms.SendReturns(twilio.MessageResource{}, errors.New("wrong"))

	cfg := &config.Configuration{}

	tw := Twilio{cl: psms, cfg: cfg}

	err := tw.Send(messaging.Message{Type: messaging.TimeCriticalSMS})
	assert.EqualError(t, err, "wrong")

	err = tw.Send(messaging.Message{})
	assert.EqualError(t, err, unknownType)
}
