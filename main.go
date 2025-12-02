package main

import (
	"log"
	"github.com/carlo-colombo/sopra/client"
	"github.com/carlo-colombo/sopra/config"
	"github.com/carlo-colombo/sopra/server"
	"github.com/carlo-colombo/sopra/service"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	log.Printf(cfg.String()) // Print the loaded configuration

	if cfg.OpenSkyClient.ID == "" || cfg.OpenSkyClient.Secret == "" {
		log.Fatal("OPENSKY_CLIENT_ID and OPENSKY_CLIENT_SECRET environment variables are required")
	}

	openskyClient := client.NewOpenSkyClient(cfg.OpenSkyClient.ID, cfg.OpenSkyClient.Secret)
	appService := service.NewService(openskyClient)

	httpServer := server.NewServer(appService, cfg)
	httpServer.Start()
}
