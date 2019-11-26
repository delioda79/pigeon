package env

import (
	"github.com/beatlabs/patron/log"
	"os"
)

// MustGetEnv gets a value from the environment or logs fatal
func MustGetEnv(key string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("Missing configuration %s", key)
	}
	return v
}
