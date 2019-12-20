package http

import (
	"bytes"
	"context"
	"errors"
	"github.com/beatlabs/patron/encoding/json"
	"github.com/beatlabs/patron/sync"
	"github.com/stretchr/testify/assert"
	"github.com/taxibeat/pigeon/internal/config"
	"github.com/taxibeat/pigeon/internal/messaging"
	"github.com/taxibeat/pigeon/internal/messaging/messagingfakes"
	"testing"
)

func TestSend(t *testing.T) {
	sdr := &messagingfakes.FakeSender{}
	cfg := &config.Configuration{}
	ingrr := New(sdr, cfg)

	bts, err := json.Encode(messaging.Message{})
	assert.Nil(t, err)

	tdtd := []struct {
		ok  bool
		bts []byte
		err error
		rtr messaging.MessageResource
	}{
		{false, []byte{}, nil, messaging.MessageResource{}},
		{true, bts, nil, messaging.MessageResource{}},
		{false, bts, errors.New("wrong"), messaging.MessageResource{}},
	}

	for _, td := range tdtd {

		sdr.SendReturns(td.rtr, td.err)

		rsp, err := ingrr.send(context.Background(), sync.NewRequest(map[string]string{}, bytes.NewReader(td.bts), map[string]string{}, json.Decode))
		if td.ok {
			assert.Nil(t, err)
			assert.NotNil(t, rsp)
			rsp, ok := rsp.Payload.(messaging.MessageResource)
			assert.True(t, ok)
			assert.Equal(t, messaging.MessageResource{}, rsp)
		} else {
			assert.NotNil(t, err)
			assert.Nil(t, rsp)
		}
	}
}

func TestIngestion_Routes(t *testing.T) {
	sdr := &messagingfakes.FakeSender{}
	ingrr := New(sdr, &config.Configuration{})

	assert.Len(t, ingrr.Routes(), 1)
}
