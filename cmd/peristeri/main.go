package main

import (
	"fmt"
	"os"

	"github.com/beatlabs/patron"
	"github.com/beatlabs/patron/log"

	"context"
	"github.com/beatlabs/patron/sync"
	sync_http "github.com/beatlabs/patron/sync/http"
	"net/http"

	"github.com/beatlabs/patron/async"

	"github.com/beatlabs/patron/async/kafka"

	"github.com/beatlabs/patron/async/amqp"
)

var (
	version = "dev"
)

func main() {
	name := "peristeri"

	err := patron.Setup(name, version)
	if err != nil {
		fmt.Printf("failed to set up logging: %v", err)
		os.Exit(1)
	}

	var oo []patron.OptionFunc

	// Set up HTTP routes
	routes := make([]sync_http.Route, 0)
	// Append a GET route
	routes = append(routes, sync_http.NewRoute("/", http.MethodGet, func(ctx context.Context, req *sync.Request) (*sync.Response, error) {
		return sync.NewResponse("Get data"), nil
	}, true, nil))

	oo = append(oo, patron.Routes(routes))

	// Set up Kafka
	kafkaCf, err := kafka.New(name, "json.Type", "TOPIC", "GROUP", []string{"BROKER"})
	if err != nil {
		log.Fatalf("failed to create kafka consumer factory: %v", err)
	}

	kafkaCmp, err := async.New("RENAME", nil, kafkaCf)
	if err != nil {
		log.Fatalf("failed to create kafka async component: %v", err)
	}

	oo = append(oo, patron.Components(kafkaCmp))

	// Set up Amqp
	amqpCf, err := amqp.New("URL", "QUEUE", "EXCHANGE")
	if err != nil {
		log.Fatalf("failed to create amqp consumer factory: %v", err)
	}

	amqpCmp, err := async.New("RENAME", nil, amqpCf)
	if err != nil {
		log.Fatalf("failed to create kafka async component: %v", err)
	}

	oo = append(oo, patron.Components(amqpCmp))

	srv, err := patron.New(name, version, oo...)
	if err != nil {
		log.Fatalf("failed to create service %v", err)
	}

	err = srv.Run()
	if err != nil {
		log.Fatalf("failed to run service %v", err)
	}

}
