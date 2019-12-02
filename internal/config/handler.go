package config

import (
	"github.com/beatlabs/harvester"
	"github.com/beatlabs/harvester/sync"
	"github.com/beatlabs/patron/errors"
)

// Configuration holds all the configuration for harvester
type Configuration struct {
	KafkaBroker               sync.String `env:"PIGEON_KAFKA_BROKER"`
	KafkaGroup                sync.String `env:"PIGEON_KAFKA_GROUP"`
	KafkaTopic                sync.String `env:"PIGEON_KAFKA_TOPIC"`
	TwilioTimeCriticalPool    sync.String `env:"PIGEON_TIME_CRITICAL_SMS_POOL"`
	TwilioNonTimeCriticalPool sync.String `env:"PIGEON_NON_TIME_CRITICAL_SMS_POOL"`
	TwilioCallBack            sync.String `env:"PIGEON_TWILIO_CALLBACK_PATH"`
	RestURL                   sync.String `env:"PIGEON_REST_URL"`
	TwilioSID                 sync.String `env:"PIGEON_TWILIO_ACCOUNT_SID"`
	TwilioToken               sync.String `env:"PIGEON_TWILIO_ACCOUNT_TOKEN"`
}

// NewConfig instantiates a new configuration object
func NewConfig(cfg *Configuration) (harvester.Harvester, error) {
	if cfg == nil {
		return nil, errors.New("Empty configuration provided")
	}
	h, err := harvester.New(cfg).Create()

	return h, err
}
