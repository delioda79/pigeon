package twilio

import (
	"github.com/taxibeat/pigeon/internal/config"
	"github.com/taxibeat/pigeon/twilio"
)

type Twilio struct {
	cfg *config.Configuration
	cl  twilio.ProgrammableSMS
}
