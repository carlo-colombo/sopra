package main

import (
	"log"
	"sopra/client"
	"sopra/config"
	"sopra/server"
	"sopra/service"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	if cfg.OpenSkyClient.ID == "" || cfg.OpenSkyClient.Secret == "" {
		log.Fatal("OPENREDISKY_CLIENT_ID and OPENREDISKY_CLIENT_SECRET environment variables are required")
	}

	openskyClient := client.NewOpenSkyClient(cfg.OpenSkyClient.ID, cfg.OpenSkyClient.Secret)
	appService := service.NewService(openskyClient)

	httpServer := server.NewServer(appService, cfg)
	httpServer.Start()
}
