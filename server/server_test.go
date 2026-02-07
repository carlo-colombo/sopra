package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/carlo-colombo/sopra/config"
	"github.com/carlo-colombo/sopra/database"
	"github.com/carlo-colombo/sopra/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockService is a mock implementation of the Service.
type MockService struct {
	mock.Mock
}

func (m *MockService) GetFlightsInRadius(lat, lon, radius float64) ([]model.FlightInfo, error) {
	args := m.Called(lat, lon, radius)
	return args.Get(0).([]model.FlightInfo), args.Error(1)
}

func newTestDB(t *testing.T) *database.DB {
	t.Helper()
	dbName := fmt.Sprintf("%s.db", t.Name())
	os.Remove(dbName) // Clean up before test
	db, err := database.NewDB(dbName)
	if err != nil {
		t.Fatalf("failed to create test db: %v", err)
	}
	t.Cleanup(func() {
		db.Close()
		os.Remove(dbName)
	})
	return db
}

func TestGetFlightsHandler(t *testing.T) {
	// Create a mock service
	mockService := new(MockService)
	expectedFlights := []model.FlightInfo{
		{Ident: "UAL123", Operator: "United Airlines"},
		{Ident: "DAL456", Operator: "Delta Airlines"},
	}
	mockService.On("GetFlightsInRadius", mock.Anything, mock.Anything, mock.Anything).Return(expectedFlights, nil)

	// Create a new server with the mock service
	cfg := &config.Config{}
	server := NewServer(mockService, cfg, nil)

	// Create a new HTTP request
	req, err := http.NewRequest("GET", "/flights", nil)
	assert.NoError(t, err)

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.getFlightsHandler)

	// Serve the HTTP request
	handler.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check the response body
	var actualFlights []model.FlightInfo
	err = json.Unmarshal(rr.Body.Bytes(), &actualFlights)
	assert.NoError(t, err)
	assert.Equal(t, expectedFlights, actualFlights)

	// Assert that the mock was called
	mockService.AssertExpectations(t)
}

func TestGetAllFlightsHandler(t *testing.T) {
	// Create a new in-memory database for testing
	db := newTestDB(t)
	if err := db.ClearFlightLog(); err != nil {
		t.Fatalf("failed to clear flight log: %v", err)
	}

	// Log some dummy flight data
	flight1 := &model.FlightInfo{
		Ident:    "FL001",
		Operator: "TestAir",
		Origin: model.AirportDetail{
			City: "Testville",
			Code: "TST",
		},
		Destination: model.AirportDetail{
			City: "Testburg",
			Code: "TSB",
		},
		AircraftType: "B737",
		Distance:     1234.5,
	}
	flight2 := &model.FlightInfo{
		Ident:    "FL002",
		Operator: "TestAir",
		Origin: model.AirportDetail{
			City: "Testburg",
			Code: "TSB",
		},
		Destination: model.AirportDetail{
			City: "Testville",
			Code: "TST",
		},
		AircraftType: "A320",
		Distance:     5678.9,
	}
	if err := db.LogFlight("FL001", flight1); err != nil {
		t.Fatalf("failed to log flight FL001: %v", err)
	}
	if err := db.LogFlight("FL002", flight2); err != nil {
		t.Fatalf("failed to log flight FL002: %v", err)
	}

	// Create a new server with the test database
	cfg := &config.Config{}
	server := NewServer(nil, cfg, db)

	// Create a new HTTP request
	req, err := http.NewRequest("GET", "/all-flights", nil)
	assert.NoError(t, err)

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.getAllFlightsHandler)

	// Serve the HTTP request
	handler.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check the response body
	var actualFlights []map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &actualFlights)
	assert.NoError(t, err)
	assert.Len(t, actualFlights, 2)
	assert.Equal(t, "FL002", actualFlights[0]["flight"])
	assert.Equal(t, 5678.9, actualFlights[0]["distance_m"])
	assert.Equal(t, "FL001", actualFlights[1]["flight"])
	assert.Equal(t, 1234.5, actualFlights[1]["distance_m"])
}
func TestGetStatsHandler(t *testing.T) {
	// Create a new in-memory database for testing
	db := newTestDB(t)

	// Log some dummy flight data
	flight1 := &model.FlightInfo{
		Ident:    "FL001",
		Operator: "TestAir",
		Origin: model.AirportDetail{
			City: "Testville",
			Code: "TST",
		},
		Destination: model.AirportDetail{
			City: "Testburg",
			Code: "TSB",
		},
		AircraftType: "B737",
		Distance:     1234.5,
	}
	flight2 := &model.FlightInfo{
		Ident:    "FL002",
		Operator: "TestAir",
		Origin: model.AirportDetail{
			City: "Testburg",
			Code: "TSB",
		},
		Destination: model.AirportDetail{
			City: "Testville",
			Code: "TST",
		},
		AircraftType: "A320",
		Distance:     5678.9,
	}
	if err := db.LogFlight("FL001", flight1); err != nil {
		t.Fatalf("failed to log flight FL001: %v", err)
	}
	if err := db.LogFlight("FL002", flight2); err != nil {
		t.Fatalf("failed to log flight FL002: %v", err)
	}

	// Create a new server with the test database
	cfg := &config.Config{}
	server := NewServer(nil, cfg, db)

	// Create a new HTTP request
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.getStatsHandler)

	// Serve the HTTP request
	handler.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check the response body
	body := rr.Body.String()
	assert.Contains(t, body, "<h1>Flight Statistics</h1>")
	assert.Contains(t, body, "<h2>Last Flight Seen</h2>")
	assert.Contains(t, body, "<td>FL002</td>")
	assert.Contains(t, body, "<h2>Last 10 Flights Seen</h2>")
	assert.Contains(t, body, "<td>FL001</td>")
	assert.Contains(t, body, "<h2>5 Most Common Flights</h2>")
}
