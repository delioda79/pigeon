package http

import (
	"github.com/stretchr/testify/assert"
	"github.com/taxibeat/pigeon/internal/config"
	"testing"
)

func TestConfig(t *testing.T) {
	cfg := &config.Configuration{}
	rtr, err := New(cfg)
	assert.Nil(t, err)

	sch := rtr.getConfigSchema()

	assert.Empty(t, sch.Static)
	assert.Len(t, sch.Dynamic, 4)
}
