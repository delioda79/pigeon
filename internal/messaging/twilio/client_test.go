package twilio

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/taxibeat/pigeon/internal/config"
	"github.com/taxibeat/pigeon/internal/messaging"
	"github.com/taxibeat/pigeon/twilio"
	"github.com/taxibeat/pigeon/twilio/progrsms"
	"github.com/taxibeat/pigeon/twilio/twiliofakes"
	"net/http"
	"net/http/httptest"
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
	cfg.TwilioNonTimeCriticalPool.Set("pool1")

	tw := Twilio{cl: psms, cfg: cfg}

	msg := messaging.Message{
		ID:        "a1234",
		Content:   "Hello there",
		Recipient: "receiver1",
		Critical:  true,
	}

	rs, err := tw.Send(msg)
	assert.Nil(t, err)
	assert.Equal(t, messaging.MessageResource{msg, "", ""}, rs)

	mc := psms.SendArgsForCall(0)

	assert.Equal(t, mc.Body, msg.Content)
	assert.Equal(t, mc.From, cfg.TwilioTimeCriticalPool.Get())
	assert.Equal(t, mc.To, msg.Recipient)
	assert.Equal(t, mc.StatusCallback, fmt.Sprintf("REST/path/%s", msg.ID))

	msg.ID = ""
	rs, err = tw.Send(msg)
	mc = psms.SendArgsForCall(1)
	assert.Nil(t, err)
	assert.Equal(t, messaging.MessageResource{msg, "", ""}, rs)
	assert.Equal(t, mc.StatusCallback, "")
}

func TestError(t *testing.T) {
	psms := &twiliofakes.FakeProgrammableSMS{}
	psms.SendReturns(twilio.MessageResource{}, errors.New(notImplemented))

	cfg := &config.Configuration{}

	tw := Twilio{cl: psms, cfg: cfg}

	_, err := tw.Send(messaging.Message{Critical: true})
	assert.EqualError(t, err, notImplemented)

	_, err = tw.Send(messaging.Message{Critical: true})
	assert.EqualError(t, err, notImplemented)

	_, err = tw.Send(messaging.Message{})
	assert.EqualError(t, err, notImplemented)
}

func TestLocalTwilioServer(t *testing.T) {

	msg := messaging.Message{
		Content:   "Hello there",
		Recipient: "receiver1",
		Critical:  true,
	}

	dd := []struct {
		ID       string
		callback bool
	}{
		{ID: "a1234", callback: true},
		{ID: "", callback: false},
	}

	for _, d := range dd {
		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			err := req.ParseForm()
			assert.NoError(t, err)
			if d.ID != "" {
				sc, ok := req.Form["StatusCallback"]
				assert.True(t, ok)
				assert.NotEmpty(t, sc)
			} else {
				_, ok := req.Form["StatusCallback"]
				assert.False(t, ok)
			}
		}))

		sid := "asid"
		psms, err := progrsms.New(nil, sid, "atoken", fmt.Sprintf("%s/%s%s", srv.URL, sid, twilio.MSGURL))
		assert.Nil(t, err)
		cfg := &config.Configuration{}
		cfg.TwilioToken.Set("token")
		cfg.TwilioSID.Set("TSID")
		cfg.RestURL.Set("REST")
		cfg.TwilioCallBack.Set("/path")
		cfg.TwilioTimeCriticalPool.Set("pool1")
		cfg.TwilioNonTimeCriticalPool.Set("pool1")

		tw, err := New(cfg)
		assert.NoError(t, err)
		tw.cl = psms

		msg.ID = d.ID
		_, err = tw.Send(msg)
		assert.Nil(t, err)
	}
}
