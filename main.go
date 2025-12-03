package main

import (
	"github.com/carlo-colombo/sopra/client"
	"github.com/carlo-colombo/sopra/config"
	"github.com/carlo-colombo/sopra/server"
	"github.com/carlo-colombo/sopra/service"
	"log"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	log.Printf("%s", cfg.String()) // Print the loaded configuration

	if cfg.OpenSkyClient.ID == "" || cfg.OpenSkyClient.Secret == "" {
		log.Fatal("OPENSKY_CLIENT_ID and OPENSKY_CLIENT_SECRET environment variables are required")
	}

	if cfg.FlightAware.APIKey == "" {
		log.Fatal("FLIGHTAWARE_API_KEY environment variable is required")
	}

	openskyClient := client.NewOpenSkyClient(cfg.OpenSkyClient.ID, cfg.OpenSkyClient.Secret)
	flightawareClient := client.NewFlightAwareClient(cfg.FlightAware.APIKey)
	appService := service.NewService(openskyClient, flightawareClient)

	httpServer := server.NewServer(appService, cfg)
	httpServer.Start()
}
