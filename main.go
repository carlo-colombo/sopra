package main

import (
	"encoding/json"
	"fmt"
	"github.com/carlo-colombo/sopra/client"
	"github.com/carlo-colombo/sopra/config"
	"github.com/carlo-colombo/sopra/server"
	"github.com/carlo-colombo/sopra/service"
	"github.com/spf13/pflag"
	"log"
	"os"
)

func main() {
	pflag.Bool("print", false, "Print the result and logs to stdout")
	pflag.Parse()

	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	config.ConfigureLogger(cfg.Print)

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

	if cfg.Print {
		flights, err := appService.GetFlightsInRadius(cfg.Service.Latitude, cfg.Service.Longitude, cfg.Service.Radius)
		if err != nil {
			log.Printf("Error getting flights: %v", err)
			// Print an empty JSON array of FlightInfo or a JSON error object
			jsonError, marshalErr := json.MarshalIndent(map[string]interface{}{
				"error":   fmt.Sprintf("Failed to retrieve flights: %v", err),
				"flights": []interface{}{},
			}, "", "  ")
			if marshalErr != nil {
				log.Fatalf("Error marshalling error response to JSON: %v", marshalErr)
			}
			fmt.Println(string(jsonError))
			os.Exit(1) // Exit with an error code to indicate failure
		}
		jsonFlights, err := json.MarshalIndent(flights, "", "  ")
		if err != nil {
			log.Fatalf("Error marshalling flights to JSON: %v", err)
		}
		fmt.Println(string(jsonFlights))
		os.Exit(0)
	}

	httpServer := server.NewServer(appService, cfg)
	httpServer.Start()
}
