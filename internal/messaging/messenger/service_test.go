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
	cfg.NecessarySMS.Set("first")

	firstSender := &messagingfakes.FakeSender{}
	firstSender.SendReturns(nil)
	secondSender := &messagingfakes.FakeSender{}
	secondSender.SendReturns(nil)

	senders := map[string]messaging.Sender{
		"first":  firstSender,
		"second": secondSender,
	}

	mng := Service{cfg, senders}

	err := mng.Send(messaging.Message{Type: messaging.Necessary})
	assert.Nil(t, err)

	err = mng.Send(messaging.Message{Type: "unknown"})
	assert.NotNil(t, err)

	cfg.NecessarySMS.Set("second")
	err = mng.Send(messaging.Message{Type: messaging.Necessary})
	assert.Nil(t, err)

	cfg.NecessarySMS.Set("third")
	err = mng.Send(messaging.Message{Type: messaging.Necessary})
	assert.NotNil(t, err)

}
