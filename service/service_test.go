package service

import (
	"errors"
	"testing"
	"time"

	"github.com/carlo-colombo/sopra/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockOpenSkyClient is a mock implementation of the OpenSkyAPIClient interface.
type MockOpenSkyClient struct {
	mock.Mock
}

func (m *MockOpenSkyClient) GetStatesInRadius(lat, lon, radiusKm float64) ([]model.Flight, error) {
	args := m.Called(lat, lon, radiusKm)
	return args.Get(0).([]model.Flight), args.Error(1)
}

// MockFlightAwareClient is a mock implementation of the FlightAwareAPIClient interface.
type MockFlightAwareClient struct {
	mock.Mock
}

func (m *MockFlightAwareClient) GetFlightInfo(ident string) (*model.FlightInfo, error) {
	args := m.Called(ident)
	// Check if the first argument is a nil pointer, indicating no flight info
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.FlightInfo), args.Error(1)
}

func TestGetFlightsInRadius(t *testing.T) {
	// Arrange
	mockOpenSkyClient := new(MockOpenSkyClient)
	mockFlightAwareClient := new(MockFlightAwareClient)

	// Mock OpenSky client to return a list of flights
	openskyFlights := []model.Flight{
		{
			Icao24:    "a1b2c3",
			Callsign:  "UAL123",
			Latitude:  40.0,
			Longitude: -74.0,
		},
		{
			Icao24:    "d4e5f6",
			Callsign:  "", // Flight without callsign
			Latitude:  41.0,
			Longitude: -75.0,
		},
	}
	mockOpenSkyClient.On("GetStatesInRadius", 40.7128, -74.0060, 100.0).Return(openskyFlights, nil)

	// Mock FlightAware client to return flight info for UAL123
	flightAwareInfo := &model.FlightInfo{
		Ident:        "UAL123",
		Operator:     "United Airlines",
		AircraftType: "B738",
		Origin:       model.AirportDetail{AirportCode: "KORD"},
		Destination:  model.AirportDetail{AirportCode: "KLAX"},
		Status:       "En Route",
		ScheduledOut: time.Now(),
	}
	mockFlightAwareClient.On("GetFlightInfo", "UAL123").Return(flightAwareInfo, nil)
	// For the flight without callsign, expect no call to FlightAware
	mockFlightAwareClient.On("GetFlightInfo", "").Return(nil, nil).Maybe()


	service := NewService(mockOpenSkyClient, mockFlightAwareClient)

	// Act
	flights, err := service.GetFlightsInRadius(40.7128, -74.0060, 100.0)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, flights)
	assert.Len(t, flights, 1, "Expected only one enriched flight to be returned")

	// Check the enriched flight details
	assert.Equal(t, flightAwareInfo.Ident, flights[0].Ident)
	assert.Equal(t, flightAwareInfo.Operator, flights[0].Operator)
	assert.Equal(t, flightAwareInfo.Origin.AirportCode, flights[0].Origin.AirportCode)
	assert.Equal(t, flightAwareInfo.Destination.AirportCode, flights[0].Destination.AirportCode)

	mockOpenSkyClient.AssertExpectations(t)
	mockFlightAwareClient.AssertExpectations(t)
}

func TestGetFlightsInRadius_OpenSkyError(t *testing.T) {
	mockOpenSkyClient := new(MockOpenSkyClient)
	mockFlightAwareClient := new(MockFlightAwareClient)

	mockOpenSkyClient.On("GetStatesInRadius", mock.Anything, mock.Anything, mock.Anything).Return([]model.Flight{}, errors.New("opensky error"))

	service := NewService(mockOpenSkyClient, mockFlightAwareClient)

	flights, err := service.GetFlightsInRadius(40.7128, -74.0060, 100.0)

	assert.Error(t, err)
	assert.Nil(t, flights)
	mockOpenSkyClient.AssertExpectations(t)
	mockFlightAwareClient.AssertNotCalled(t, "GetFlightInfo", mock.Anything)
}

func TestGetFlightsInRadius_FlightAwareError(t *testing.T) {
	mockOpenSkyClient := new(MockOpenSkyClient)
	mockFlightAwareClient := new(MockFlightAwareClient)

	openskyFlights := []model.Flight{
		{
			Icao24:    "a1b2c3",
			Callsign:  "UAL123",
			Latitude:  40.0,
			Longitude: -74.0,
		},
	}
	mockOpenSkyClient.On("GetStatesInRadius", mock.Anything, mock.Anything, mock.Anything).Return(openskyFlights, nil)
	mockFlightAwareClient.On("GetFlightInfo", "UAL123").Return(nil, errors.New("flightaware error"))

	service := NewService(mockOpenSkyClient, mockFlightAwareClient)

	flights, err := service.GetFlightsInRadius(40.7128, -74.0060, 100.0)

	assert.NoError(t, err) // Service continues on FlightAware error
	assert.Empty(t, flights, "Expected no enriched flights if FlightAware lookup fails")

	mockOpenSkyClient.AssertExpectations(t)
	mockFlightAwareClient.AssertExpectations(t)
}

func TestGetFlightsInRadius_NoCallsign(t *testing.T) {
	mockOpenSkyClient := new(MockOpenSkyClient)
	mockFlightAwareClient := new(MockFlightAwareClient)

	openskyFlights := []model.Flight{
		{
			Icao24:    "a1b2c3",
			Callsign:  "", // No callsign
			Latitude:  40.0,
			Longitude: -74.0,
		},
	}
	mockOpenSkyClient.On("GetStatesInRadius", mock.Anything, mock.Anything, mock.Anything).Return(openskyFlights, nil)

	service := NewService(mockOpenSkyClient, mockFlightAwareClient)

	flights, err := service.GetFlightsInRadius(40.7128, -74.0060, 100.0)

	assert.NoError(t, err)
	assert.Empty(t, flights, "Expected no enriched flights if OpenSky flight has no callsign")

	mockOpenSkyClient.AssertExpectations(t)
	mockFlightAwareClient.AssertNotCalled(t, "GetFlightInfo", mock.Anything)
}
