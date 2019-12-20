package main

import (
	"context"
	"fmt"
	"github.com/beatlabs/patron"
	"github.com/beatlabs/patron/log"
	http3 "github.com/beatlabs/patron/sync/http"
	"github.com/taxibeat/pigeon/internal/config"
	"github.com/taxibeat/pigeon/internal/config/sanitarium"
	sanithttp "github.com/taxibeat/pigeon/internal/config/sanitarium/http"
	"github.com/taxibeat/pigeon/internal/ingestion/http"
	"github.com/taxibeat/pigeon/internal/ingestion/kafka"
	"github.com/taxibeat/pigeon/internal/messaging/messenger"
	"os"

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

	err = sanitarium.NewConfig(context.Background(), cfg)
	if err != nil {
		log.Fatalf("Impossible to create configuration %v", err)
	}

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
	var rr []http3.Route

	sdr, err := messenger.New(cfg)
	if err != nil {
		log.Fatalf("failed to create new messenger: %v", err)
	}

	rndp := http.New(sdr, cfg)

	rr = append(rr, rndp.Routes()...)

	// Set up Kafka

	kfkTimeCrCmp, err := kafka.New(name, true, cfg, sdr)
	if err != nil {
		log.Fatalf("failed to create kafka async component: %v", err)
	}

	kfkNonTimeCrCmp, err := kafka.New(name, false, cfg, sdr)
	if err != nil {
		log.Fatalf("failed to create kafka async component: %v", err)
	}

	oo = append(oo, patron.Components(kfkTimeCrCmp, kfkNonTimeCrCmp))

	// Add sanitarium routes
	sanitariumSrv, err := sanithttp.New(cfg)
	if err != nil {
		log.Fatalf("Cannot create HTTP component for sanitarium")
	}
	rr = append(rr, sanitariumSrv.GetRoutes()...)

	oo = append(oo, patron.Routes(rr))

	srv, err := patron.New(name, version, oo...)
	if err != nil {
		log.Fatalf("failed to create service %v", err)
	}

	err = srv.Run()
	if err != nil {
		log.Fatalf("failed to run service %v", err)
	}
}
