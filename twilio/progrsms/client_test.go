package progrsms

import (
	"encoding/json"
	"fmt"
	"github.com/beatlabs/patron/errors"
	"github.com/stretchr/testify/assert"
	"github.com/taxibeat/pigeon/twilio"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestLocalCall(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		bd, err := ioutil.ReadAll(req.Body)
		assert.Nil(t, err)
		vls, _ := url.ParseQuery(string(bd))
		assert.Equal(
			t,
			url.Values{
				"Body":       []string{"A story"},
				"From":       []string{"+15005550006"},
				"To":         []string{"+447988346883"},
				"ForceOptIn": []string{"true"},
			},
			vls,
		)

		rsp := twilio.MessageResource{Body: vls.Get("Body"), From: vls.Get("From"), To: vls.Get("To")}
		bts, err := json.Marshal(rsp)
		assert.Nil(t, err)
		rw.Write(bts)
	}))

	td := []struct {
		msg twilio.MessageCreate
		err error
	}{
		{msg: twilio.MessageCreate{From: "+15005550006", Body: "A story", To: "+447988346883", ForceOptIn: true}, err: nil},
		{msg: twilio.MessageCreate{Body: "A story", To: "+447988346883", ForceOptIn: true}, err: errors.New(requiredFromErrorMessage)},
		{msg: twilio.MessageCreate{From: "+15005550006", To: "+447988346883", ForceOptIn: true}, err: errors.New(requiredBodyErrorMessage)},
	}

	cl, err := New(nil, "AC7fd971c3663f0feef37019f5b359b97a", "1733407bdf0a95790012b6b88ffb9c5f", fmt.Sprintf("%s/", ts.URL))
	assert.Nil(t, err)

	for _, d := range td {
		rsp, err := cl.Send(d.msg)
		if d.err != nil {
			assert.EqualError(t, d.err, err.Error())
		} else {
			assert.Nil(t, err)
			assert.Equal(t, twilio.MessageResource{Body: d.msg.Body, From: d.msg.From, To: d.msg.To}, rsp)
		}
	}
}
