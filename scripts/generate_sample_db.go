package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/carlo-colombo/sopra/database"
	"github.com/carlo-colombo/sopra/model"
)

func main() {
	dbPath := filepath.Join("sample", "sopra.db")

	// Ensure sample directory exists
	if err := os.MkdirAll("sample", 0755); err != nil {
		log.Fatalf("failed to create sample directory: %v", err)
	}

	// Remove existing db if it exists
	os.Remove(dbPath)

	// Initialize the db (this runs migrations and seeds)
	db, err := database.NewDB(dbPath)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()

	// Clear the default seed data from migrations
	if err := db.ClearFlightLog(); err != nil {
		log.Fatalf("failed to clear flight log: %v", err)
	}

	operators := []struct {
		icao string
		json string
	}{
		{"SWR", `{"name":"Swiss International Air Lines","shortname":"Swiss","country":"Switzerland"}`},
		{"DLH", `{"name":"Lufthansa","shortname":"Lufthansa","country":"Germany"}`},
		{"AFR", `{"name":"Air France","shortname":"Air France","country":"France"}`},
		{"BAW", `{"name":"British Airways","shortname":"British Airways","country":"United Kingdom"}`},
		{"UAE", `{"name":"Emirates","shortname":"Emirates","country":"United Arab Emirates"}`},
		{"QFA", `{"name":"Qantas","shortname":"Qantas","country":"Australia"}`},
		{"SIA", `{"name":"Singapore Airlines","shortname":"Singapore Airlines","country":"Singapore"}`},
		{"ACA", `{"name":"Air Canada","shortname":"Air Canada","country":"Canada"}`},
		{"JAL", `{"name":"Japan Airlines","shortname":"Japan Airlines","country":"Japan"}`},
		{"KLM", `{"name":"KLM Royal Dutch Airlines","shortname":"KLM","country":"Netherlands"}`},
	}

	for _, op := range operators {
		if err := db.LogOperator(op.icao, op.json); err != nil {
			log.Printf("failed to log operator %s: %v", op.icao, err)
		}
	}

	airports := []model.AirportDetail{
		{CodeIata: "ZRH", Name: "Zurich Airport", City: "Zurich"},
		{CodeIata: "JFK", Name: "John F. Kennedy International Airport", City: "New York"},
		{CodeIata: "FRA", Name: "Frankfurt Airport", City: "Frankfurt"},
		{CodeIata: "LAX", Name: "Los Angeles International Airport", City: "Los Angeles"},
		{CodeIata: "CDG", Name: "Charles de Gaulle Airport", City: "Paris"},
		{CodeIata: "SFO", Name: "San Francisco International Airport", City: "San Francisco"},
		{CodeIata: "LHR", Name: "Heathrow Airport", City: "London"},
		{CodeIata: "ORD", Name: "O'Hare International Airport", City: "Chicago"},
		{CodeIata: "DXB", Name: "Dubai International Airport", City: "Dubai"},
		{CodeIata: "MIA", Name: "Miami International Airport", City: "Miami"},
		{CodeIata: "SYD", Name: "Sydney Kingsford Smith Airport", City: "Sydney"},
		{CodeIata: "DFW", Name: "Dallas/Fort Worth International Airport", City: "Dallas"},
		{CodeIata: "SIN", Name: "Singapore Changi Airport", City: "Singapore"},
		{CodeIata: "EWR", Name: "Newark Liberty International Airport", City: "Newark"},
		{CodeIata: "YYZ", Name: "Toronto Pearson International Airport", City: "Toronto"},
		{CodeIata: "HKG", Name: "Hong Kong International Airport", City: "Hong Kong"},
		{CodeIata: "HND", Name: "Haneda Airport", City: "Tokyo"},
		{CodeIata: "SEA", Name: "Seattle-Tacoma International Airport", City: "Seattle"},
		{CodeIata: "AMS", Name: "Amsterdam Airport Schiphol", City: "Amsterdam"},
		{CodeIata: "MEX", Name: "Mexico City International Airport", City: "Mexico City"},
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Create a pool of ~40 unique flight identifiers to ensure collisions
	flightIdents := make([]string, 40)
	for i := 0; i < 40; i++ {
		op := operators[r.Intn(len(operators))]
		flightIdents[i] = fmt.Sprintf("%s%d", op.icao, 100+r.Intn(900))
	}

	for i := 0; i < 100; i++ {
		ident := flightIdents[r.Intn(len(flightIdents))]
		opIcao := ident[:3]

		originIdx := r.Intn(len(airports))
		destIdx := r.Intn(len(airports))
		for destIdx == originIdx {
			destIdx = r.Intn(len(airports))
		}

		aircraftTypes := []string{"A320", "A333", "A359", "A388", "B738", "B748", "B77W", "B789", "GLF6"}

		f := &model.FlightInfo{
			Ident:         ident,
			OperatorIcao:  opIcao,
			Origin:        airports[originIdx],
			Destination:   airports[destIdx],
			Latitude:      47.0 + r.Float64(),
			Longitude:     8.0 + r.Float64(),
			AircraftType:  aircraftTypes[r.Intn(len(aircraftTypes))],
			RouteDistance: 500 + r.Intn(5000),
		}

		if err := db.LogFlight(f.Ident, f); err != nil {
			log.Printf("failed to log flight %s: %v", f.Ident, err)
		}
	}

	log.Printf("Sample database created at %s with 100 flight logs", dbPath)
}
