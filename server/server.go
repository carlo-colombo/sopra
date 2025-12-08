package server

import (
	"encoding/json"
	"fmt"
	"github.com/carlo-colombo/sopra/config"
	"github.com/carlo-colombo/sopra/model"
	"log"
	"net/http"
	"time"

	"github.com/carlo-colombo/sopra/database"
)

// FlightService defines the interface for the flight service.
type FlightService interface {
	GetFlightsInRadius(lat, lon, radius float64) ([]model.FlightInfo, error)
}

// Server holds the HTTP server and its dependencies.
type Server struct {
	service FlightService
	config  *config.Config
	db      *database.DB
}

// NewServer creates a new Server instance.
func NewServer(s FlightService, cfg *config.Config, db *database.DB) *Server {
	return &Server{
		service: s,
		config:  cfg,
		db:      db,
	}
}

// Start starts the HTTP server.
func (s *Server) Start() {
	http.HandleFunc("/flights", s.getFlightsHandler)
	http.HandleFunc("/last-flight", s.getLastFlightHandler)
	http.HandleFunc("/all-flights", s.getAllFlightsHandler)

	port := fmt.Sprintf(":%d", s.config.Port)
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func (s *Server) getLastFlightHandler(w http.ResponseWriter, r *http.Request) {
	flight, lastSeen, err := s.db.GetLatestFlight()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if flight == nil {
		http.Error(w, "No flight data available", http.StatusNotFound)
		return
	}

	response := struct {
		Flight          string    `json:"flight"`
		Operator        string    `json:"operator"`
		DestinationCity string    `json:"destination_city"`
		DestinationCode string    `json:"destination_code_iata"`
		SourceCity      string    `json:"source_city"`
		SourceCode      string    `json:"source_code_iata"`
		LastTimeSeen    time.Time `json:"last_time_seen"`
		AirplaneModel   string    `json:"airplane_model"`
	}{
		Flight:          flight.Ident,
		Operator:        flight.Operator,
		DestinationCity: flight.Destination.City,
		DestinationCode: flight.Destination.Code,
		SourceCity:      flight.Origin.City,
		SourceCode:      flight.Origin.Code,
		LastTimeSeen:    lastSeen,
		AirplaneModel:   flight.AircraftType,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) getAllFlightsHandler(w http.ResponseWriter, r *http.Request) {
	flights, lastSeens, err := s.db.GetAllFlights()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(flights) == 0 {
		http.Error(w, "No flight data available", http.StatusNotFound)
		return
	}

	type FlightResponse struct {
		Flight          string    `json:"flight"`
		Operator        string    `json:"operator"`
		DestinationCity string    `json:"destination_city"`
		DestinationCode string    `json:"destination_code_iata"`
		SourceCity      string    `json:"source_city"`
		SourceCode      string    `json:"source_code_iata"`
		LastTimeSeen    time.Time `json:"last_time_seen"`
		AirplaneModel   string    `json:"airplane_model"`
	}

	var responses []FlightResponse
	for i, flight := range flights {
		response := FlightResponse{
			Flight:          flight.Ident,
			Operator:        flight.Operator,
			DestinationCity: flight.Destination.City,
			DestinationCode: flight.Destination.Code,
			SourceCity:      flight.Origin.City,
			SourceCode:      flight.Origin.Code,
			LastTimeSeen:    lastSeens[i],
			AirplaneModel:   flight.AircraftType,
		}
		responses = append(responses, response)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responses); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
