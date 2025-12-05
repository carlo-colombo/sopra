package config

import (
	"log"
	"os"
)

// ConfigureLogger sets up the logger to write to stderr.
func ConfigureLogger() {
	log.SetOutput(os.Stderr)
}
