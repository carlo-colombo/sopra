package server

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/carlo-colombo/sopra/config"
	"github.com/carlo-colombo/sopra/database"
	"github.com/carlo-colombo/sopra/model"
)

//go:embed statics/index.html
var indexHTML string

// FlightService defines the interface for the flight service.
type FlightService interface {
	GetFlightsInRadius(lat, lon, radius float64) ([]model.FlightInfo, error)
}

// Server holds the HTTP server and its dependencies.
type Server struct {
	service  FlightService
	config   *config.Config
	db       *database.DB
	template *template.Template
}

// NewServer creates a new Server instance.
func NewServer(s FlightService, cfg *config.Config, db *database.DB) *Server {
	tmpl, err := template.New("index").Parse(indexHTML)
	if err != nil {
		log.Fatalf("failed to parse template: %v", err)
	}
	return &Server{
		service:  s,
		config:   cfg,
		db:       db,
		template: tmpl,
	}
}

// Start starts the HTTP server.
func (s *Server) Start() {
	http.HandleFunc("/", s.getStatsHandler)
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

func (s *Server) getStatsHandler(w http.ResponseWriter, r *http.Request) {
	lastFlight, lastFlightSeen, err := s.db.GetLatestFlight()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	last10Flights, last10FlightsSeen, err := s.db.GetLast10Flights()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	mostCommonFlights, err := s.db.GetMostCommonFlights()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type FlightData struct {
		*model.FlightInfo
		LastSeen time.Time
	}

	var lastFlightData *FlightData
	if lastFlight != nil {
		lastFlightData = &FlightData{
			FlightInfo: lastFlight,
			LastSeen:   lastFlightSeen,
		}
	}

	var last10FlightsData []FlightData
	for i, flight := range last10Flights {
		last10FlightsData = append(last10FlightsData, FlightData{
			FlightInfo: flight,
			LastSeen:   last10FlightsSeen[i],
		})
	}

	data := struct {
		LastFlight        *FlightData
		Last10Flights     []FlightData
		MostCommonFlights []*model.FlightInfo
	}{
		LastFlight:        lastFlightData,
		Last10Flights:     last10FlightsData,
		MostCommonFlights: mostCommonFlights,
	}

	if err := s.template.Execute(w, data); err != nil {
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
