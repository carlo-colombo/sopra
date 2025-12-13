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

//go:embed statics/flight_table.html
var flightTableHTML string

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
	funcMap := template.FuncMap{
		"formatTime": func(t time.Time) string {
			return t.Local().Format("02/01/2006 15:04")
		},
	}

	tmpl, err := template.New("index").Funcs(funcMap).Parse(indexHTML)
	if err != nil {
		log.Fatalf("failed to parse template: %v", err)
	}

	_, err = tmpl.Parse(flightTableHTML)
	if err != nil {
		log.Fatalf("failed to parse flight_table template: %v", err)
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

func (s *Server) getOperatorInfo(icao string) (*model.OperatorInfo, error) {
	operatorJSON, err := s.db.GetOperator(icao)
	if err != nil {
		return nil, err
	}
	if operatorJSON == "" {
		return &model.OperatorInfo{Shortname: "N/A"}, nil
	}
	var operatorInfo model.OperatorInfo
	if err := json.Unmarshal([]byte(operatorJSON), &operatorInfo); err != nil {
		return nil, err
	}
	return &operatorInfo, nil
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

	operator, err := s.getOperatorInfo(flight.OperatorIcao)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		Operator:        operator.Shortname,
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

	// Collect all unique operator ICAOs
	icaoSet := make(map[string]struct{})
	if lastFlight != nil {
		icaoSet[lastFlight.OperatorIcao] = struct{}{}
	}
	for _, flight := range last10Flights {
		icaoSet[flight.OperatorIcao] = struct{}{}
	}
	for _, flight := range mostCommonFlights {
		icaoSet[flight.OperatorIcao] = struct{}{}
	}
	var icaos []string
	for icao := range icaoSet {
		icaos = append(icaos, icao)
	}

	// Bulk fetch operator info
	operatorMap, err := s.db.GetOperators(icaos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type FlightData struct {
		*model.FlightInfo
		Operator *model.OperatorInfo
		LastSeen time.Time
	}

	var lastFlightData *FlightData
	if lastFlight != nil {
		var operator model.OperatorInfo
		if opJSON, ok := operatorMap[lastFlight.OperatorIcao]; ok {
			if err := json.Unmarshal([]byte(opJSON), &operator); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			operator.Shortname = "N/A"
		}
		lastFlightData = &FlightData{
			FlightInfo: lastFlight,
			Operator:   &operator,
			LastSeen:   lastFlightSeen,
		}
	}

	var last10FlightsData []FlightData
	for i, flight := range last10Flights {
		var operator model.OperatorInfo
		if opJSON, ok := operatorMap[flight.OperatorIcao]; ok {
			if err := json.Unmarshal([]byte(opJSON), &operator); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			operator.Shortname = "N/A"
		}
		last10FlightsData = append(last10FlightsData, FlightData{
			FlightInfo: flight,
			Operator:   &operator,
			LastSeen:   last10FlightsSeen[i],
		})
	}

	type MostCommonFlightData struct {
		*model.FlightInfo
		Operator *model.OperatorInfo
	}

	var mostCommonFlightsData []MostCommonFlightData
	for _, flight := range mostCommonFlights {
		var operator model.OperatorInfo
		if opJSON, ok := operatorMap[flight.OperatorIcao]; ok {
			if err := json.Unmarshal([]byte(opJSON), &operator); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			operator.Shortname = "N/A"
		}
		mostCommonFlightsData = append(mostCommonFlightsData, MostCommonFlightData{
			FlightInfo: flight,
			Operator:   &operator,
		})
	}

	data := struct {
		LastFlight        interface{}
		Last10Flights     interface{}
		MostCommonFlights interface{}
	}{
		LastFlight: map[string]interface{}{
			"Flights": []FlightData{*lastFlightData},
			"Header":  "Last Seen",
			"Class":   "last-flight",
		},
		Last10Flights: map[string]interface{}{
			"Flights": last10FlightsData,
			"Header":  "Last Seen",
			"Class":   "last-10-flights",
		},
		MostCommonFlights: map[string]interface{}{
			"Flights": mostCommonFlightsData,
			"Header":  "Count",
			"Class":   "most-common-flights",
		},
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

	// Collect all unique operator ICAOs
	icaoSet := make(map[string]struct{})
	for _, flight := range flights {
		icaoSet[flight.OperatorIcao] = struct{}{}
	}
	var icaos []string
	for icao := range icaoSet {
		icaos = append(icaos, icao)
	}

	// Bulk fetch operator info
	operatorMap, err := s.db.GetOperators(icaos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		var operator model.OperatorInfo
		if opJSON, ok := operatorMap[flight.OperatorIcao]; ok {
			if err := json.Unmarshal([]byte(opJSON), &operator); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			operator.Shortname = "N/A"
		}
		response := FlightResponse{
			Flight:          flight.Ident,
			Operator:        operator.Shortname,
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
