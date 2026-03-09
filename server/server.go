package server

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/carlo-colombo/sopra/config"
	"github.com/carlo-colombo/sopra/database"
	"github.com/carlo-colombo/sopra/model"
	"github.com/hako/durafmt"
	// "github.com/carlo-colombo/sopra/service" // Removed as no longer used
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

// formatTimeAgo returns a human-readable string indicating how long ago a time was.
func formatTimeAgo(t time.Time) string {
	d := time.Since(t)
	if d < time.Minute {
		return "just now"
	}
	s := durafmt.Parse(d).LimitFirstN(2).String()
	parts := strings.Split(s, " ")
	if len(parts) == 4 {
		return fmt.Sprintf("%s %s and %s %s ago", parts[0], parts[1], parts[2], parts[3])
	}
	return s + " ago"
}

// NewServer creates a new Server instance.
func NewServer(s FlightService, cfg *config.Config, db *database.DB) *Server {
	loc, err := time.LoadLocation(cfg.Timezone)
	if err != nil {
		log.Printf("failed to load location %s: %v. Falling back to Local", cfg.Timezone, err)
		loc = time.Local
	}

	funcMap := template.FuncMap{
		"formatTime": func(t time.Time) string {
			return t.In(loc).Format("02/01/2006 15:04")
		},
		"timeAgo": formatTimeAgo,
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

// formatNumberWithThousandsSeparator adds thousand separators to a float64 number,
// rounding to a whole number, and returns it as a string.
func formatNumberWithThousandsSeparator(n float64) string {
	// Round to the nearest whole number
	rounded := int64(n + 0.5) // Add 0.5 for proper rounding, then convert to int64

	// Convert to string
	s := strconv.FormatInt(rounded, 10)

	// Add commas
	nSpaces := (len(s) - 1) / 3
	// Pre-allocate memory for the result string (original length + number of commas)
	var result strings.Builder
	result.Grow(len(s) + nSpaces)

	for i, r := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result.WriteRune(',')
		}
		result.WriteRune(r)
	}
	return result.String()
}

// Start starts the HTTP server.
func (s *Server) Start() error {
	http.HandleFunc("/", s.getStatsHandler)
	http.HandleFunc("/flights", s.getFlightsHandler)
	http.HandleFunc("/last-flight", s.getLastFlightHandler)
	http.HandleFunc("/all-flights", s.getAllFlightsHandler)

	port := fmt.Sprintf(":%d", s.config.Port)
	if err := http.ListenAndServe(port, nil); err != nil {
		return fmt.Errorf("server failed to start: %w", err)
	}
	return nil
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
		Flight              string    `json:"flight"`
		Operator            string    `json:"operator"`
		DestinationCity     string    `json:"destination_city"`
		DestinationCodeIata string    `json:"destination_code_iata"`
		DestinationCodeIcao string    `json:"destination_code_icao"`
		SourceCity          string    `json:"source_city"`
		SourceCodeIata      string    `json:"source_code_iata"`
		SourceCodeIcao      string    `json:"source_code_icao"`
		LastTimeSeen        time.Time `json:"last_time_seen"`
		LastSeenAgo         string    `json:"last_seen_ago"`
		AirplaneModel       string    `json:"airplane_model"`
		Distance            float64   `json:"distance_m"` // Reverted to float64
		CO2KG               float64   `json:"co2_kg"`     // Reverted to float64
	}{
		Flight:              flight.Ident,
		Operator:            operator.Shortname,
		DestinationCity:     flight.Destination.City,
		DestinationCodeIata: flight.Destination.CodeIata,
		DestinationCodeIcao: flight.Destination.CodeIcao,
		SourceCity:          flight.Origin.City,
		SourceCodeIata:      flight.Origin.CodeIata,
		SourceCodeIcao:      flight.Origin.CodeIcao,
		LastTimeSeen:        lastSeen,
		LastSeenAgo:         formatTimeAgo(lastSeen),
		AirplaneModel:       flight.AircraftType,
		Distance:            flight.Distance, // Assign raw float64
		CO2KG:               flight.CO2KG,    // Assign raw float64
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
		// Populate Display fields for HTML template
		lastFlight.DistanceDisplay = formatNumberWithThousandsSeparator(lastFlight.Distance / 1000)
		lastFlight.CO2KGDisplay = fmt.Sprintf("%.0f", lastFlight.CO2KG)

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
		// Populate Display fields for HTML template
		flight.DistanceDisplay = formatNumberWithThousandsSeparator(flight.Distance / 1000)
		flight.CO2KGDisplay = fmt.Sprintf("%.0f", flight.CO2KG)

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
		// Populate Display fields for HTML template
		flight.DistanceDisplay = formatNumberWithThousandsSeparator(flight.Distance / 1000)
		flight.CO2KGDisplay = fmt.Sprintf("%.0f", flight.CO2KG)

		mostCommonFlightsData = append(mostCommonFlightsData, MostCommonFlightData{
			FlightInfo: flight,
			Operator:   &operator,
		})
	}

	topDestinations, err := s.db.GetTopDestinations()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	topSources, err := s.db.GetTopSources()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type StatWithPerc struct {
		model.AirportStat
		Percentage float64
	}

	getStatsWithPerc := func(stats []model.AirportStat) []StatWithPerc {
		max := 0
		for _, s := range stats {
			if s.Count > max {
				max = s.Count
			}
		}
		var res []StatWithPerc
		for _, s := range stats {
			perc := 0.0
			if max > 0 {
				perc = float64(s.Count) / float64(max) * 100
			}
			res = append(res, StatWithPerc{s, perc})
		}
		return res
	}

	destStats := getStatsWithPerc(topDestinations)
	srcStats := getStatsWithPerc(topSources)

	data := struct {
		LastFlight        interface{}
		Last10Flights     interface{}
		MostCommonFlights interface{}
		TopDestinations   []StatWithPerc
		TopSources        []StatWithPerc
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
		TopDestinations: destStats,
		TopSources:      srcStats,
	}

	if err := s.template.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) getAllFlightsHandler(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 0
	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, "invalid limit parameter", http.StatusBadRequest)
			return
		}
	}

	flights, lastSeens, err := s.db.GetAllFlights(limit)
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
		Flight              string    `json:"flight"`
		Operator            string    `json:"operator"`
		DestinationCity     string    `json:"destination_city"`
		DestinationCodeIata string    `json:"destination_code_iata"`
		DestinationCodeIcao string    `json:"destination_code_icao"`
		SourceCity          string    `json:"source_city"`
		SourceCodeIata      string    `json:"source_code_iata"`
		SourceCodeIcao      string    `json:"source_code_icao"`
		LastTimeSeen        time.Time `json:"last_time_seen"`
		LastSeenAgo         string    `json:"last_seen_ago"`
		AirplaneModel       string    `json:"airplane_model"`
		Distance            float64   `json:"distance_m"` // Reverted to float64
		CO2KG               float64   `json:"co2_kg"`     // Reverted to float64
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
			Flight:              flight.Ident,
			Operator:            operator.Shortname,
			DestinationCity:     flight.Destination.City,
			DestinationCodeIata: flight.Destination.CodeIata,
			DestinationCodeIcao: flight.Destination.CodeIcao,
			SourceCity:          flight.Origin.City,
			SourceCodeIata:      flight.Origin.CodeIata,
			SourceCodeIcao:      flight.Origin.CodeIcao,
			LastTimeSeen:        lastSeens[i],
			LastSeenAgo:         formatTimeAgo(lastSeens[i]),
			AirplaneModel:       flight.AircraftType,
			Distance:            flight.Distance, // Assign raw float64
			CO2KG:               flight.CO2KG,    // Assign raw float64
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
