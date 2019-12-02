package twilio

const (
	// BaseTwilioURL is the base URL for the Twilio API
	BaseTwilioURL = "https://api.twilio.com/2010-04-01/Accounts/"
	// MSGURL is the schema definition for the messgaing API
	MSGURL = "/Messages.json"
)

// MessageResource is the result of a message sending
type MessageResource struct {
	AccountSID          string `json:"AccountSID"`
	APIVersion          string `json:"api_version"`
	Body                string `json:"body"`
	DateCreated         string `json:"date_created"`
	DateUpdates         string `json:"date_updates"`
	DateSent            string `json:"date_sent"`
	Direction           string `json:"direction"`
	ErrorCode           int    `json:"error_code"`
	ErrorMessage        string `json:"error_message"`
	From                string `json:"from"`
	MessagingServiceSID string `json:"messaging_service_sid"`
	NumMedia            string `json:"num_media"`
	NumSegments         string `json:"num_segments"`
	Price               string `json:"price"`
	PriceUnit           string `json:"price_unit"`
	SID                 string `json:"sid"`
	Status              string `json:"status"`
	To                  string `json:"to"`
	URI                 string `json:"uri"`
}

// MessageCreate represent a message request
type MessageCreate struct {
	To                  string   `json:"To"`
	StatusCallback      string   `json:"StatusCallback,omitempty"`
	ApplicationSID      string   `json:"ApplicationSid,omitempty"`
	MaxPrice            float64  `json:"MaxPrice,omitempty"`
	ProvideFeedback     bool     `json:"ProvideFeedback,omitempty"`
	ValidityPeriod      int64    `json:"ValidityPeriod,omitempty"`
	MaxRate             string   `json:"MaxRate,omitempty"`
	ForceDelivery       bool     `json:"ForceDelivery,omitempty"`
	ProviderSID         string   `json:"ProviderSid,omitempty"`
	AddressRetention    string   `json:"AddressRetention,omitempty"`
	SmartEncoded        bool     `json:"SmartEncoded,omitempty"`
	PersistentAction    []string `json:"PersistentAction,omitempty"`
	TransientAction     []string `json:"TransientAction,omitempty"`
	Title               string   `json:"Title,omitempty"`
	InteractiveData     string   `json:"InteractiveData,omitempty"`
	ForceOptIn          bool     `json:"ForceOptIn,omitempty"`
	RichLinkData        string   `json:"RichLinkData,omitempty"`
	TrafficType         string   `json:"TrafficType,omitempty"`
	From                string   `json:"From,omitempty"`
	MessagingServiceSID string   `json:"MessagingServiceSid,omitempty"`
	Body                string   `json:"Body"`
	MediaURL            []string `json:"MediaUrl,omitempty"`
}

// ProgrammableSMS is an interface to interact with the Programmable SMS Twilio API
type ProgrammableSMS interface {
	Send(m MessageCreate) (MessageResource, error)
}
