package main

import (
	"fmt"
	"github.com/taxibeat/peristeri/internal/env"
	"os"

	"github.com/beatlabs/patron"
	"github.com/beatlabs/patron/log"

	"github.com/beatlabs/patron/async"

	"github.com/beatlabs/patron/async/kafka"
	"github.com/joho/godotenv"
)

var (
	version = "0.1"
	kafkaBroker, kafkaGroup, kafkaTopic,
	name string
)

func init() {
	name = "peristeri"

	err := patron.Setup(name, version)
	if err != nil {
		fmt.Printf("failed to set up logging: %v", err)
		os.Exit(1)
	}
	err = godotenv.Load("../../config/.env")
	if err != nil {
		log.Debugf("no .env file exists: %v", err)
	}

	kafkaBroker = env.MustGetEnv("PERISTERI_KAFKA_BROKER")
	kafkaGroup = env.MustGetEnv("PERISTERI_KAFKA_GROUP")
	kafkaTopic = env.MustGetEnv("PERISTERI_KAFKA_TOPIC")
}

func main() {

	var oo []patron.OptionFunc

	// Set up Kafka
	kafkaCf, err := kafka.New(name, kafkaTopic, kafkaGroup, []string{kafkaBroker})
	if err != nil {
		log.Fatalf("failed to create kafka consumer factory: %v", err)
	}

	kafkaCmp, err := async.New("RENAME", func(m async.Message) error { return nil }, kafkaCf)
	if err != nil {
		log.Fatalf("failed to create kafka async component: %v", err)
	}

	oo = append(oo, patron.Components(kafkaCmp))

	srv, err := patron.New(name, version, oo...)
	if err != nil {
		log.Fatalf("failed to create service %v", err)
	}

	err = srv.Run()
	if err != nil {
		log.Fatalf("failed to run service %v", err)
	}

}
