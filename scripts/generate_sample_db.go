package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/carlo-colombo/sopra/database"
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

	// Add some operator data to make the demo look better
	operators := map[string]string{
		"SWR": `{"name":"Swiss International Air Lines","shortname":"Swiss","country":"Switzerland"}`,
		"DLH": `{"name":"Lufthansa","shortname":"Lufthansa","country":"Germany"}`,
		"AFR": `{"name":"Air France","shortname":"Air France","country":"France"}`,
		"BAW": `{"name":"British Airways","shortname":"British Airways","country":"United Kingdom"}`,
		"UAE": `{"name":"Emirates","shortname":"Emirates","country":"United Arab Emirates"}`,
		"QFA": `{"name":"Qantas","shortname":"Qantas","country":"Australia"}`,
		"SIA": `{"name":"Singapore Airlines","shortname":"Singapore Airlines","country":"Singapore"}`,
		"ACA": `{"name":"Air Canada","shortname":"Air Canada","country":"Canada"}`,
		"JAL": `{"name":"Japan Airlines","shortname":"Japan Airlines","country":"Japan"}`,
		"KLM": `{"name":"KLM Royal Dutch Airlines","shortname":"KLM","country":"Netherlands"}`,
	}

	for icao, jsonValue := range operators {
		if err := db.LogOperator(icao, jsonValue); err != nil {
			log.Printf("failed to log operator %s: %v", icao, err)
		}
	}

	log.Printf("Sample database created at %s", dbPath)
}
