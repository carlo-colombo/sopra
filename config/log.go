package config

import (
	"io"
	"log"
	"os"
)

// ConfigureLogger sets up the logger to write to a file or stdout.
func ConfigureLogger(print bool) {
	if print {
		log.SetOutput(os.Stdout)
	} else {
		file, err := os.OpenFile("sopra.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}
		mw := io.MultiWriter(os.Stdout, file)
		log.SetOutput(mw)
	}
}
