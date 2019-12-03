package messenger

import (
	"github.com/stretchr/testify/assert"
	"github.com/taxibeat/pigeon/internal/config"
	"github.com/taxibeat/pigeon/internal/messaging"
	"github.com/taxibeat/pigeon/internal/messaging/messagingfakes"
	"testing"
)

func TestSend(t *testing.T) {
	cfg := &config.Configuration{}
	cfg.TwilioTimeCriticalPool.Set("first")

	msgRs := messaging.MessageResource{Status: "received", ProviderID: "1"}

	firstSender := &messagingfakes.FakeSender{}

	senders := &settings{
		twilio: firstSender,
	}

	mng := Service{cfg, senders}

	msgRs.Type = messaging.TimeCriticalSMS
	firstSender.SendReturns(msgRs, nil)
	rs, err := mng.Send(messaging.Message{Type: messaging.TimeCriticalSMS})
	assert.Nil(t, err)
	assert.Equal(t, rs, msgRs)

	rs, err = mng.Send(messaging.Message{Type: "unknownProvider"})
	assert.NotNil(t, err)
	assert.Equal(t, rs, messaging.MessageResource{})
}
