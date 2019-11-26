package config

import (
	"github.com/beatlabs/harvester"
	"github.com/beatlabs/harvester/sync"
	"github.com/beatlabs/patron/errors"
)

type Configuration struct {
	KafkaBroker  sync.String `env:"PIGEON_KAFKA_BROKER"`
	KafkaGroup   sync.String `env:"PIGEON_KAFKA_GROUP"`
	KafkaTopic   sync.String `env:"PIGEON_KAFKA_TOPIC"`
	NecessarySMS sync.String `env:"PIGEON_NECESSARY_SMS"`
}

func NewConfig(cfg *Configuration) (harvester.Harvester, error) {
	if cfg == nil {
		return nil, errors.New("Empty configuration provided")
	}
	h, err := harvester.New(cfg).Create()

	return h, err
}
