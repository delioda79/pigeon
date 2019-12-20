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
	KafkaTimeCriticalTopic    sync.String `env:"PIGEON_KAFKA_TOPIC_TIME_CRITICAL"`
	KafkaNonTimeCriticalTopic sync.String `env:"PIGEON_KAFKA_TOPIC_NON_TIME_CRITICAL"`
	TwilioTimeCriticalPool    sync.String `env:"PIGEON_TIME_CRITICAL_SMS_POOL" consul:"services/pigeon/twilio-time-critical-pool"`
	TwilioNonTimeCriticalPool sync.String `env:"PIGEON_NON_TIME_CRITICAL_SMS_POOL" consul:"services/pigeon/twilio-time-non-critical-pool"`
	TwilioCallBack            sync.String `env:"PIGEON_TWILIO_CALLBACK_PATH"`
	RestURL                   sync.String `env:"PIGEON_REST_URL"`
	TwilioSID                 sync.String `env:"PIGEON_TWILIO_ACCOUNT_SID"`
	TwilioToken               sync.String `env:"PIGEON_TWILIO_ACCOUNT_TOKEN"`
	HTTPEnabled               sync.Bool   `env:"PIGEON_HTTP_MESSAGE_CONSUMER_ENABLED" consul:"services/pigeon/http-enabled"`
	KafkaConsumerEnabled      sync.Bool   `env:"PIGEON_KAFKA_MESSAGE_CONSUMER_ENABLED" consul:"services/pigeon/kafka-enabled"`
}

// NewConfig instantiates a new configuration object
func NewConfig(cfg *Configuration) (harvester.Harvester, error) {
	if cfg == nil {
		return nil, errors.New("Empty configuration provided")
	}
	h, err := harvester.New(cfg).Create()

	return h, err
}
