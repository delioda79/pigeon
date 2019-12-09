package kafka

import (
	"context"
	"github.com/beatlabs/patron/errors"
	"github.com/stretchr/testify/assert"
	"github.com/taxibeat/pigeon/internal/config"
	"github.com/taxibeat/pigeon/internal/messaging/messagingfakes"
	"testing"
)

type KfkM struct {
	err       error
	typeError bool
}

func (m *KfkM) Context() context.Context {
	return context.Background()
}
func (m *KfkM) Decode(v interface{}) error {

	rsp, ok := v.(*message)
	if !ok || m.typeError == true {
		return errors.New("wrong type")
	}

	v = rsp

	return m.err
}
func (m *KfkM) Ack() error {
	return nil
}
func (m *KfkM) Nack() error {
	return nil
}

func TestNew(t *testing.T) {
	cfg := &config.Configuration{}
	cfg.KafkaTimeCriticalTopic.Set("a")

	snd := &messagingfakes.FakeSender{}
	r, e := New("cmp1", true, cfg, snd)
	assert.Nil(t, e)
	assert.NotNil(t, r)

	r, e = New("cmp1", false, cfg, snd)
	assert.Nil(t, e)
	assert.NotNil(t, r)
}

func TestIngestionConsumer_Process(t *testing.T) {
	snd := &messagingfakes.FakeSender{}
	cns := IngestionConsumer{critical: true, sender: snd}

	msg := &KfkM{err: errors.New("wrong")}

	err := cns.Process(msg)
	assert.EqualError(t, err, "wrong")

	msg = &KfkM{err: nil, typeError: false}
	err = cns.Process(msg)
	assert.Nil(t, err)

	cns.critical = false
	msg = &KfkM{err: nil, typeError: false}
	err = cns.Process(msg)
	assert.Nil(t, err)

	assert.Equal(t, snd.SendArgsForCall(0).Critical, true)
	assert.Equal(t, snd.SendArgsForCall(1).Critical, false)
}
