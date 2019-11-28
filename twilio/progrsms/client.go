package progrsms

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/beatlabs/patron/errors"
	phttp "github.com/beatlabs/patron/trace/http"
	"github.com/taxibeat/pigeon/twilio"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	requiredFromErrorMessage = "From and Messaging Service ID cannot be both empty"
	requiredBodyErrorMessage = "Body and Media URL cannot be both empty"
)

type Client struct {
	hc         phttp.Client
	accountSID string
	authToken  string
	url        string
}

func (cl *Client) Send(m twilio.MessageCreate) (twilio.MessageResource, error) {

	if m.From == "" && m.MessagingServiceSID == "" {
		return twilio.MessageResource{}, errors.New(requiredFromErrorMessage)
	}

	if m.Body == "" && len(m.MediaURL) == 0 {
		return twilio.MessageResource{}, errors.New(requiredBodyErrorMessage)
	}

	msgData, err := cl.setValues(m)
	if err != nil {
		return twilio.MessageResource{}, err
	}

	msgDataReader := *strings.NewReader(msgData.Encode())

	req, err := http.NewRequest("POST", cl.url, &msgDataReader)
	if err != nil {
		return twilio.MessageResource{}, err
	}
	cl.setBasicRequestDetails(req)

	resp, err := cl.hc.Do(context.Background(), req)
	if err != nil {
		return twilio.MessageResource{}, err
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data *twilio.MessageResource
		decoder := json.NewDecoder(resp.Body)
		if decoder.Decode(&data); err != nil {
			return twilio.MessageResource{}, errors.Errorf("Impossible unmarshal response: %v", err)
		}
		return *data, nil
	}

	bts, _ := ioutil.ReadAll(resp.Body)

	return twilio.MessageResource{}, errors.Errorf("Request returned status %d and message : %v", resp.StatusCode, string(bts))
}

func (cl *Client) setValues(m twilio.MessageCreate) (url.Values, error) {
	bts, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	var vl map[string]interface{}
	if err := json.Unmarshal(bts, &vl); err != nil {
		return nil, err
	}

	msgData := url.Values{}
	for k, v := range vl {
		msgData.Set(k, fmt.Sprintf("%v", v))
	}
	return msgData, nil
}

func (cl *Client) setBasicRequestDetails(req *http.Request) {
	req.SetBasicAuth(
		cl.accountSID,
		cl.authToken,
	)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
}

func New(cl phttp.Client, sid, token, baseURL string) (*Client, error) {
	var err error
	if cl == nil {
		cl, err = phttp.New()
		if err != nil {
			return nil, err
		}
	}

	if baseURL == "" {
		baseURL = twilio.BaseTwilioURL
	}

	return &Client{cl, sid, token, fmt.Sprintf("%s%s%s", baseURL, sid, twilio.MSGURL)}, nil
}
