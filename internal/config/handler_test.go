package config

import (
	"context"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewCOnfiguration(t *testing.T) {
	os.Setenv("PIGEON_KAFKA_BROKER", "testbroker")
	os.Setenv("PIGEON_KAFKA_GROUP", "testgroup")
	os.Setenv("PIGEON_KAFKA_TOPIC", "testtopic")
	cfg := &Configuration{}
	h, e := NewConfig(cfg)

	assert.Nil(t, e)
	assert.Equal(t, "", cfg.KafkaBroker.Get())
	assert.Equal(t, "", cfg.KafkaGroup.Get())
	assert.Equal(t, "", cfg.KafkaTopic.Get())
	assert.Equal(t, "", cfg.NecessarySMS.Get())

	h.Harvest(context.Background())

	assert.Equal(t, "testbroker", cfg.KafkaBroker.Get())
	assert.Equal(t, "testgroup", cfg.KafkaGroup.Get())
	assert.Equal(t, "testtopic", cfg.KafkaTopic.Get())
	assert.Equal(t, "", cfg.NecessarySMS.Get())
}

func TestWrongConfig(t *testing.T) {
	_, e := NewConfig(nil)

	assert.NotNil(t, e)
}
