package main

import (
	"context"
	"fmt"
	"github.com/beatlabs/patron/async"
	"github.com/beatlabs/patron/async/kafka"
	"github.com/taxibeat/pigeon/internal/config"
	"github.com/taxibeat/pigeon/internal/ingestion/http"
	"github.com/taxibeat/pigeon/internal/messaging/messenger"
	"os"

	"github.com/beatlabs/patron"
	"github.com/beatlabs/patron/log"

	"github.com/joho/godotenv"
)

var (
	version = "0.1"
	name    string
	cfg     = &config.Configuration{}
)

func init() {
	name = "pigeon"

	err := patron.Setup(name, version)
	if err != nil {
		fmt.Printf("failed to set up logging: %v", err)
		os.Exit(1)
	}
	err = godotenv.Load("../../config/.env")
	if err != nil {
		log.Debugf("no .env file exists: %v", err)
	}

	h, err := config.NewConfig(cfg)
	if err != nil {
		log.Fatalf("Impossibe to retrieve configuration: %v", err)
	}

	h.Harvest(context.Background())

	if cfg.KafkaBroker.Get() == "" {
		log.Fatalf("No value defined for kafka broker")
	}

	if cfg.KafkaTopic.Get() == "" {
		log.Fatalf("No value defined for kafka topic")
	}

	if cfg.KafkaGroup.Get() == "" {
		log.Fatalf("No value defined for kafka group")
	}
}

func main() {

	log.Fatalf("Config: %+v", cfg)

	var oo []patron.OptionFunc

	sdr, err := messenger.New(cfg)
	if err != nil {
		log.Fatalf("failed to create new messenger: %v", err)
	}

	rndp := http.New(sdr)

	oo = append(oo, patron.Routes(rndp.Routes()))

	// Set up Kafka
	kafkaCf, err := kafka.New(name, cfg.KafkaTopic.Get(), cfg.KafkaGroup.Get(), []string{cfg.KafkaBroker.Get()})
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
