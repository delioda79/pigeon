package sanitarium

import (
	"context"
	"github.com/beatlabs/harvester"
	"github.com/beatlabs/harvester/monitor/consul"
	"github.com/beatlabs/harvester/sync"
	"github.com/taxibeat/pigeon/internal/config"
	"time"
)

// Schema is the sanitarium schema
type Schema struct {
	Name        string       `json:"name"`
	Icon        string       `json:"icon"`
	Description string       `json:"description"`
	Static      []SchemaItem `json:"static"`
	Dynamic     []SchemaItem `json:"dynamic"`
}

// SchemaItem is an item for the sanitarium schema
type SchemaItem struct {
	ID          string       `json:"id"`
	Type        string       `json:"type"`
	Label       string       `json:"label"`
	Description string       `json:"description"`
	Default     interface{}  `json:"default,omitempty"`
	Value       interface{}  `json:"value,omitempty"`
	Options     []OptionItem `json:"options,omitempty"`
}

// OptionItem represents a schema item option
type OptionItem struct {
	Value interface{} `json:"value"`
	Label string      `json:"label"`
}

type consulSeedConfig struct {
	Address        sync.String `seed:"localhost:8500" env:"SEED_CONSUL_HTTP_ADDR"`
	Datacenter     sync.String `seed:"dc1" env:"SEED_CONSUL_DATACENTER"`
	Token          sync.String `seed:"" env:"SEED_CONSUL_TOKEN"`
	TimeoutSeconds sync.Int64  `seed:"3" env:"SEED_CONSUL_TIMEOUT"`
}

type consulWatchConfig struct {
	Address        sync.String `seed:"localhost:8500" env:"WATCH_CONSUL_HTTP_ADDR"`
	Datacenter     sync.String `seed:"dc1" env:"WATCH_CONSUL_DATACENTER"`
	Token          sync.String `seed:"" env:"WATCH_CONSUL_TOKEN"`
	TimeoutSeconds sync.Int64  `seed:"300" env:"WATCH_CONSUL_TIMEOUT"`
}

// NewConfig harvests a configuration
func NewConfig(ctx context.Context, cfg *config.Configuration) error {
	csc, err := newConsulSeedConfig(ctx)
	if err != nil {
		return err
	}
	cwc, err := newConsulWatchConfig(ctx)
	if err != nil {
		return err
	}

	watchItems := buildWatchItems()
	h, err := harvester.New(cfg).
		WithConsulSeed(csc.Address.Get(), csc.Datacenter.Get(), csc.Token.Get(), time.Duration(csc.TimeoutSeconds.Get())*time.Second).
		WithConsulMonitor(cwc.Address.Get(), cwc.Datacenter.Get(), cwc.Token.Get(), time.Duration(cwc.TimeoutSeconds.Get())*time.Second, watchItems...).
		Create()

	if err != nil {
		return err
	}

	err = h.Harvest(ctx)
	if err != nil {
		return err
	}

	return nil
}

func newConsulSeedConfig(ctx context.Context) (*consulSeedConfig, error) {
	cfg := &consulSeedConfig{}
	h, err := harvester.New(cfg).
		Create()

	if err != nil {
		return nil, err
	}

	err = h.Harvest(ctx)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func newConsulWatchConfig(ctx context.Context) (*consulWatchConfig, error) {
	cfg := &consulWatchConfig{}
	h, err := harvester.New(cfg).
		Create()

	if err != nil {
		return nil, err
	}

	err = h.Harvest(ctx)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func buildWatchItems() []consul.Item {
	return []consul.Item{
		consul.NewKeyItem("services/pigeon/twilio-time-critical-pool"),
		consul.NewKeyItem("services/pigeon/twilio-time-non-critical-pool"),
		consul.NewKeyItem("services/pigeon/http-enabled"),
		consul.NewKeyItem("services/pigeon/kafka-enabled"),
	}
}
