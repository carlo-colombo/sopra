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
