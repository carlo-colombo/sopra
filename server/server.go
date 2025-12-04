package server

import (
	"encoding/json"
	"fmt"
	"github.com/carlo-colombo/sopra/config"
	"github.com/carlo-colombo/sopra/model"
	"log"
	"net/http"
)

// FlightService defines the interface for the flight service.
type FlightService interface {
	GetFlightsInRadius(lat, lon, radius float64) ([]model.FlightInfo, error)
}

// Server holds the HTTP server and its dependencies.
type Server struct {
	service FlightService
	config  *config.Config
}

// NewServer creates a new Server instance.
func NewServer(s FlightService, cfg *config.Config) *Server {
	return &Server{
		service: s,
		config:  cfg,
	}
}

// Start starts the HTTP server.
func (s *Server) Start() {
	http.HandleFunc("/flights", s.getFlightsHandler)

	port := fmt.Sprintf(":%d", s.config.Port)
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func (s *Server) getFlightsHandler(w http.ResponseWriter, r *http.Request) {
	flights, err := s.service.GetFlightsInRadius(s.config.Service.Latitude, s.config.Service.Longitude, s.config.Service.Radius)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(flights); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
