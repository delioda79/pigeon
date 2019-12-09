package main

import (
	"context"
	"fmt"
	"github.com/taxibeat/pigeon/internal/config"
	"github.com/taxibeat/pigeon/internal/ingestion/http"
	"github.com/taxibeat/pigeon/internal/ingestion/kafka"
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

	if cfg.KafkaTimeCriticalTopic.Get() == "" {
		log.Fatalf("No value defined for kafka topic")
	}

	if cfg.KafkaGroup.Get() == "" {
		log.Fatalf("No value defined for kafka group")
	}
}

func main() {

	var oo []patron.OptionFunc

	sdr, err := messenger.New(cfg)
	if err != nil {
		log.Fatalf("failed to create new messenger: %v", err)
	}

	if cfg.HTTPEnabled.Get() {
		rndp := http.New(sdr)

		oo = append(oo, patron.Routes(rndp.Routes()))
	}

	// Set up Kafka
	if cfg.KafkaConsumerEnabled.Get() {
		kfkTimeCrCmp, err := kafka.New(name, true, cfg, sdr)
		if err != nil {
			log.Fatalf("failed to create kafka async component: %v", err)
		}

		kfkNonTimeCrCmp, err := kafka.New(name, false, cfg, sdr)
		if err != nil {
			log.Fatalf("failed to create kafka async component: %v", err)
		}

		oo = append(oo, patron.Components(kfkTimeCrCmp, kfkNonTimeCrCmp))
	}

	srv, err := patron.New(name, version, oo...)
	if err != nil {
		log.Fatalf("failed to create service %v", err)
	}

	err = srv.Run()
	if err != nil {
		log.Fatalf("failed to run service %v", err)
	}
}
