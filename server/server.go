package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/carlo-colombo/sopra/config"
	"github.com/carlo-colombo/sopra/service"
)

// Server holds the HTTP server and its dependencies.
type Server struct {
	service *service.Service
	config  *config.Config
}

// NewServer creates a new Server instance.
func NewServer(s *service.Service, cfg *config.Config) *Server {
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
