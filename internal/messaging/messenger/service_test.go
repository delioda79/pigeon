package messenger

import (
	"errors"
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

	msgRs.Critical = true
	firstSender.SendReturns(msgRs, nil)
	rs, err := mng.Send(messaging.Message{Critical: true})
	assert.Nil(t, err)
	assert.Equal(t, rs, msgRs)
	rsp := firstSender.SendArgsForCall(0)
	assert.Equal(t, rsp.Critical, true)

	msgRs.Critical = false
	firstSender.SendReturns(msgRs, nil)
	rs, err = mng.Send(messaging.Message{Critical: false})
	assert.Nil(t, err)
	assert.Equal(t, rs, msgRs)
	rsp = firstSender.SendArgsForCall(1)
	assert.Equal(t, rsp.Critical, false)

	firstSender.SendReturns(messaging.MessageResource{}, errors.New("an error"))

	rs, err = mng.Send(messaging.Message{Critical: false})
	assert.NotNil(t, err)
	assert.Equal(t, rs, messaging.MessageResource{})
}
