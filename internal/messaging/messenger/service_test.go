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

	firstSender := &messagingfakes.FakeSender{}
	firstSender.SendReturns(nil)

	senders := &settings{
		twilio: firstSender,
	}

	mng := Service{cfg, senders}

	err := mng.Send(messaging.Message{Type: messaging.TimeCriticalSMS})
	assert.Nil(t, err)

	err = mng.Send(messaging.Message{Type: "unknownProvider"})
	assert.NotNil(t, err)

	cfg.TwilioTimeCriticalPool.Set("first")
	err = mng.Send(messaging.Message{Type: messaging.TimeCriticalSMS})
	assert.Nil(t, err)
}
