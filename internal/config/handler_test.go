package config

import (
	"context"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewConfiguration(t *testing.T) {
	os.Setenv("PIGEON_KAFKA_BROKER", "testbroker")
	os.Setenv("PIGEON_KAFKA_GROUP", "testgroup")
	os.Setenv("PIGEON_KAFKA_TOPIC_TIME_CRITICAL", "testtopic")
	cfg := &Configuration{}
	h, e := NewConfig(cfg)

	assert.Nil(t, e)
	assert.Equal(t, "", cfg.KafkaBroker.Get())
	assert.Equal(t, "", cfg.KafkaGroup.Get())
	assert.Equal(t, "", cfg.KafkaTimeCriticalTopic.Get())
	assert.Equal(t, "", cfg.TwilioTimeCriticalPool.Get())

	h.Harvest(context.Background())

	assert.Equal(t, "testbroker", cfg.KafkaBroker.Get())
	assert.Equal(t, "testgroup", cfg.KafkaGroup.Get())
	assert.Equal(t, "testtopic", cfg.KafkaTimeCriticalTopic.Get())
	assert.Equal(t, "", cfg.TwilioTimeCriticalPool.Get())
}

func TestWrongConfig(t *testing.T) {
	_, e := NewConfig(nil)

	assert.NotNil(t, e)
}
