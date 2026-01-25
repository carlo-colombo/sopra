package model

import (
	"encoding/json"
	"io"
	"os"
	"testing"
)

func TestFlightAwareResponseDeserialization(t *testing.T) {
	jsonFilePath := "response.json"
	jsonFile, err := os.Open(jsonFilePath)
	if err != nil {
		t.Fatalf("Failed to open JSON file: %v", err)
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		t.Fatalf("Failed to read JSON file: %v", err)
	}

	var faResponse FlightAwareResponse
	err = json.Unmarshal(byteValue, &faResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if len(faResponse.Flights) == 0 {
		t.Fatal("No flights found in deserialized response")
	}

	// Basic check for one of the fields
	firstFlight := faResponse.Flights[0]
	if firstFlight.Ident != "TVS84J" {
		t.Errorf("Expected Ident TVS84J, got %s", firstFlight.Ident)
	}

	if firstFlight.Origin.Code != "LIML" {
		t.Errorf("Expected Origin Code LIML, got %s", firstFlight.Origin.Code)
	}

	if firstFlight.Destination.Code != "LDZA" {
		t.Errorf("Expected Destination Code LDZA, got %s", firstFlight.Destination.Code)
	}
}
