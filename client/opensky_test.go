package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/carlo-colombo/sopra/model"
)

func TestNewOpenSkyClient(t *testing.T) {
	client := NewOpenSkyClient("test_id", "test_secret")
	if client == nil {
		t.Error("Expected NewOpenSkyClient to return a client, but got nil")
	}
}

func TestGetStates(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		states := model.States{
			Time: 1,
			States: [][]interface{}{
				{"icao24", "callsign", "origin_country", 1.0, 1.0, 1.0, 1.0, 1.0, true, 1.0, 1.0, 1.0, nil, 1.0, nil, false, 0},
			},
		}
		json.NewEncoder(w).Encode(states)
	}))
	defer server.Close()

	client := &OpenSkyClient{
		httpClient: server.Client(),
		baseURL:    server.URL,
	}

	states, err := client.GetStates()
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	if states == nil {
		t.Fatal("Expected states, but got nil")
	}

	if len(states.States) != 1 {
		t.Errorf("Expected 1 state, but got %d", len(states.States))
	}
}

func TestGetStatesWithBoundingBox(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		states := model.States{
			Time: 1,
			States: [][]interface{}{
				{"icao24", "callsign", "origin_country", 1.0, 1.0, 1.0, 1.0, 1.0, true, 1.0, 1.0, 1.0, nil, 1.0, nil, false, 0},
			},
		}
		json.NewEncoder(w).Encode(states)
	}))
	defer server.Close()

	client := &OpenSkyClient{
		httpClient: server.Client(),
		baseURL:    server.URL,
	}

	states, err := client.GetStatesWithBoundingBox(1, 1, 1, 1)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	if states == nil {
		t.Fatal("Expected states, but got nil")
	}

	if len(states.States) != 1 {
		t.Errorf("Expected 1 state, but got %d", len(states.States))
	}
}

func TestGetStatesInRadius(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		states := model.States{
			Time: 1,
			States: [][]interface{}{
				// Flight 1: Inside radius (approx 15.7 km from 0,0)
				{"icao24_1", "CALL1", "country1", 1.0, 1.0, 0.1, 0.1, 1.0, false, 1.0, 1.0, 1.0, nil, 1.0, nil, false, 0, "", ""},
				// Flight 2: Outside radius (approx 157.2 km from 0,0)
				{"icao24_2", "CALL2", "country2", 1.0, 1.0, 1.0, 1.0, 1.0, false, 1.0, 1.0, 1.0, nil, 1.0, nil, false, 0, "", ""},
			},
		}
		json.NewEncoder(w).Encode(states)
	}))
	defer server.Close()

	client := &OpenSkyClient{
		httpClient: server.Client(),
		baseURL:    server.URL,
	}

	centerLat, centerLon := 0.0, 0.0
	radiusKm := 100.0 // 100 km radius

	flights, err := client.GetStatesInRadius(centerLat, centerLon, radiusKm)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	if len(flights) != 1 {
		t.Fatalf("Expected 1 flight within radius, but got %d", len(flights))
	}

	if flights[0].Callsign != "CALL1" {
		t.Errorf("Expected flight with callsign CALL1, but got %s", flights[0].Callsign)
	}
}
