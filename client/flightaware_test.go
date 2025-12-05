package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/carlo-colombo/sopra/database"
	"github.com/carlo-colombo/sopra/model"
)

// newTestCache creates a new in-memory cache for testing.
func newTestCache(t *testing.T) *database.Cache {
	t.Helper()
	cache, err := database.NewCache(":memory:")
	if err != nil {
		t.Fatalf("failed to create test cache: %v", err)
	}
	return cache
}

func TestNewFlightAwareClient(t *testing.T) {
	apiKey := "test_api_key"
	cache := newTestCache(t)
	client := NewFlightAwareClient(apiKey, cache)

	if client == nil {
		t.Fatal("Expected NewFlightAwareClient to return a client, but got nil")
	}
	if client.apiKey != apiKey {
		t.Errorf("Expected API key %s, but got %s", apiKey, client.apiKey)
	}
	if client.baseURL != "https://aeroapi.flightaware.com/aeroapi" {
		t.Errorf("Expected base URL %s, but got %s", "https://aeroapi.flightaware.com/aeroapi", client.baseURL)
	}
	if client.cache == nil {
		t.Error("Expected cache to be initialized, but it was nil")
	}
}

func TestGetFlightInfo_Success(t *testing.T) {
	expectedIdent := "UAL123"
	mockResponse := model.FlightAwareResponse{
		Flights: []model.FlightInfo{
			{
				Ident:        expectedIdent,
				Operator:     "United Airlines",
				AircraftType: "B738",
				Origin: model.AirportDetail{
					AirportCode: "ORD",
					AirportName: "Chicago O'Hare International Airport",
				},
				Destination: model.AirportDetail{
					AirportCode: "LAX",
					AirportName: "Los Angeles International Airport",
				},
				ScheduledOut: time.Now(),
				EstimatedOut: time.Now(),
				ActualOut:    nil,
				ScheduledOn:  time.Now(),
				EstimatedOn:  time.Now(),
				ActualOn:     nil,
				ScheduledIn:  time.Now(),
				EstimatedIn:  time.Now(),
				ActualIn:     nil,
				Status:       "En Route",
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != fmt.Sprintf("/aeroapi/flights/%s", expectedIdent) {
			t.Errorf("Expected to request '/aeroapi/flights/%s', got '%s'", expectedIdent, r.URL.Path)
		}
		if r.Header.Get("x-apikey") != "test_api_key" {
			t.Errorf("Expected API key header 'test_api_key', got '%s'", r.Header.Get("x-apikey"))
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	cache := newTestCache(t)
	client := &FlightAwareClient{
		httpClient: server.Client(),
		apiKey:     "test_api_key",
		baseURL:    server.URL + "/aeroapi", // Adjust base URL for mock server
		cache:      cache,
	}

	flightInfo, err := client.GetFlightInfo(expectedIdent)
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}
	if flightInfo == nil {
		t.Fatal("Expected flight info, but got nil")
	}
	if flightInfo.Ident != expectedIdent {
		t.Errorf("Expected ident %s, but got %s", expectedIdent, flightInfo.Ident)
	}
	if flightInfo.Origin.AirportCode != mockResponse.Flights[0].Origin.AirportCode {
		t.Errorf("Expected origin airport code %s, but got %s", mockResponse.Flights[0].Origin.AirportCode, flightInfo.Origin.AirportCode)
	}
	if flightInfo.Destination.AirportCode != mockResponse.Flights[0].Destination.AirportCode {
		t.Errorf("Expected destination airport code %s, but got %s", mockResponse.Flights[0].Destination.AirportCode, flightInfo.Destination.AirportCode)
	}
	if flightInfo.Status != "En Route" {
		t.Errorf("Expected status %s, but got %s", "En Route", flightInfo.Status)
	}
}

func TestGetFlightInfo_Cache(t *testing.T) {
	expectedIdent := "UAL123"
	mockResponse := model.FlightAwareResponse{
		Flights: []model.FlightInfo{
			{
				Ident:    expectedIdent,
				Operator: "United Airlines",
			},
		},
	}

	serverHitCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serverHitCount++
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	cache := newTestCache(t)
	client := &FlightAwareClient{
		httpClient: server.Client(),
		apiKey:     "test_api_key",
		baseURL:    server.URL + "/aeroapi",
		cache:      cache,
	}

	// First call - should hit the server
	flightInfo, err := client.GetFlightInfo(expectedIdent)
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}
	if flightInfo == nil {
		t.Fatal("Expected flight info, but got nil")
	}
	if serverHitCount != 1 {
		t.Errorf("Expected server to be hit once, but it was hit %d times", serverHitCount)
	}

	// get last seen time
	_, lastSeen1, err := cache.Get(expectedIdent)
	if err != nil {
		t.Fatalf("Error getting from cache: %v", err)
	}

	// Second call - should be served from cache
	flightInfo, err = client.GetFlightInfo(expectedIdent)
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}
	if flightInfo == nil {
		t.Fatal("Expected flight info, but got nil")
	}
	if serverHitCount != 1 {
		t.Errorf("Expected server to be hit once, but it was hit %d times", serverHitCount)
	}

	_, lastSeen2, err := cache.Get(expectedIdent)
	if err != nil {
		t.Fatalf("Error getting from cache: %v", err)
	}

	if !lastSeen1.Equal(lastSeen2) {
		t.Errorf("Expected last seen time to be the same, but it was updated")
	}
}

func TestGetFlightInfo_NotFound(t *testing.T) {
	expectedIdent := "NONEXISTENT123"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	cache := newTestCache(t)
	client := &FlightAwareClient{
		httpClient: server.Client(),
		apiKey:     "test_api_key",
		baseURL:    server.URL + "/aeroapi",
		cache:      cache,
	}

	flightInfo, err := client.GetFlightInfo(expectedIdent)
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}
	if flightInfo != nil {
		t.Fatalf("Expected no flight info, but got: %+v", flightInfo)
	}
}

func TestGetFlightInfo_ServerError(t *testing.T) {
	expectedIdent := "UAL123"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	cache := newTestCache(t)
	client := &FlightAwareClient{
		httpClient: server.Client(),
		apiKey:     "test_api_key",
		baseURL:    server.URL + "/aeroapi",
		cache:      cache,
	}

	flightInfo, err := client.GetFlightInfo(expectedIdent)
	if err == nil {
		t.Fatal("Expected an error, but got none")
	}
	if flightInfo != nil {
		t.Fatalf("Expected no flight info, but got: %+v", flightInfo)
	}
}

func TestGetFlightInfo_NoFlightsInResponse(t *testing.T) {
	expectedIdent := "UAL123"
	mockResponse := model.FlightAwareResponse{
		Flights: []model.FlightInfo{}, // Empty flights array
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	cache := newTestCache(t)
	client := &FlightAwareClient{
		httpClient: server.Client(),
		apiKey:     "test_api_key",
		baseURL:    server.URL + "/aeroapi",
		cache:      cache,
	}

	flightInfo, err := client.GetFlightInfo(expectedIdent)
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}
	if flightInfo != nil {
		t.Fatalf("Expected no flight info, but got: %+v", flightInfo)
	}
}
