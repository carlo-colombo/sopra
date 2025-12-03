package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/carlo-colombo/sopra/config"
	"github.com/carlo-colombo/sopra/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockService is a mock implementation of the Service.
type MockService struct {
	mock.Mock
}

func (m *MockService) GetFlightsInRadius(lat, lon, radius float64) ([]model.Flight, error) {
	args := m.Called(lat, lon, radius)
	return args.Get(0).([]model.Flight), args.Error(1)
}

func TestGetFlightsHandler(t *testing.T) {
	// Create a mock service
	mockService := new(MockService)
	expectedFlights := []model.Flight{
		{Icao24: "icao1", Callsign: "flight1"},
		{Icao24: "icao2", Callsign: "flight2"},
	}
	mockService.On("GetFlightsInRadius", mock.Anything, mock.Anything, mock.Anything).Return(expectedFlights, nil)

	// Create a new server with the mock service
	cfg := &config.Config{}
	server := NewServer(mockService, cfg)

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
	var actualFlights []model.Flight
	err = json.Unmarshal(rr.Body.Bytes(), &actualFlights)
	assert.NoError(t, err)
	assert.Equal(t, expectedFlights, actualFlights)

	// Assert that the mock was called
	mockService.AssertExpectations(t)
}
