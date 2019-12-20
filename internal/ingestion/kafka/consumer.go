package kafka

import (
	"errors"
	"github.com/beatlabs/patron/async"
	"github.com/beatlabs/patron/async/kafka"
	"github.com/beatlabs/patron/log"
	"github.com/taxibeat/pigeon/internal/config"
	"github.com/taxibeat/pigeon/internal/ingestion"
	"github.com/taxibeat/pigeon/internal/messaging"
)

const (
	wrongCriticalityError = "wrong criticality"
	notImplErr            = "not implemented"
)

type message struct {
	ID        string
	Recipient string
	Content   string
}

// IngestionConsumer is a consumer for notification messages
type IngestionConsumer struct {
	critical bool
	sender   messaging.Sender
	cfg      *config.Configuration
}

// Process is the kafka message processor
func (i *IngestionConsumer) Process(m async.Message) error {
	log.Debugf("Received message %v", m)

	if !i.cfg.KafkaConsumerEnabled.Get() {
		return errors.New("kafka is disabled")
	}

	req := &message{}
	if err := m.Decode(req); err != nil {
		log.Debugf("Impossible to decode %v", m)
		ingestion.ObserveCount("sms", "kafka", false, false, false, false)
		return err
	}

	msg := messaging.Message{
		ID:        req.ID,
		Critical:  i.critical,
		Recipient: req.Recipient,
		Content:   req.Content,
	}

	_, err := i.sender.Send(msg)

	log.Debugf("Message %v sent with result %v", msg, err != nil)

	ingestion.ObserveCount("sms", "http", err == nil, msg.Critical, true, true)
	return err
}

// New returns a new consumer
func New(name string, critical bool, cfg *config.Configuration, snd messaging.Sender) (*async.Component, error) {
	var topic string

	if critical {
		topic = cfg.KafkaTimeCriticalTopic.Get()
	} else {
		topic = cfg.KafkaNonTimeCriticalTopic.Get()
	}

	cns := IngestionConsumer{sender: snd, cfg: cfg}

	kafkaCf, err := kafka.New(name, topic, cfg.KafkaGroup.Get(), []string{cfg.KafkaBroker.Get()})
	if err != nil {
		log.Fatalf("failed to create kafka consumer factory: %v", err)
	}

	kafkaCmp, err := async.New("RENAME", cns.Process, kafkaCf)
	if err != nil {
		log.Fatalf("failed to create kafka async component: %v", err)
	}

	return kafkaCmp, err
}
