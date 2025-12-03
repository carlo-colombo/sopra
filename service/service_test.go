package service

import (
	"testing"

	"github.com/carlo-colombo/sopra/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockOpenSkyClient is a mock implementation of the OpenSkyAPIClient interface.
type MockOpenSkyClient struct {
	mock.Mock
}

func (m *MockOpenSkyClient) GetStatesWithBoundingBox(lamin, lomin, lamax, lomax float64) (*model.States, error) {
	args := m.Called(lamin, lomin, lamax, lomax)
	return args.Get(0).(*model.States), args.Error(1)
}

// MockFlightAwareClient is a mock implementation of the FlightAwareAPIClient interface.
type MockFlightAwareClient struct {
	mock.Mock
}

func (m *MockFlightAwareClient) GetFlightInfo(icao24 string) (origin, destination string, err error) {
	args := m.Called(icao24)
	return args.String(0), args.String(1), args.Error(2)
}

func TestGetFlightsInRadius(t *testing.T) {
	// Create a mock OpenSky client
	mockOpenSkyClient := new(MockOpenSkyClient)
	expectedStates := &model.States{
		Time: 1,
		States: [][]interface{}{
			{"icao1", "flight1", "USA", 1.0, 1.0, 1.0, 1.0, 1.0, false, 1.0, 1.0, 1.0, nil, 1.0, "1234", false, 0},
		},
	}
	mockOpenSkyClient.On("GetStatesWithBoundingBox", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(expectedStates, nil)

	// Create a mock FlightAware client
	mockFlightAwareClient := new(MockFlightAwareClient)
	mockFlightAwareClient.On("GetFlightInfo", "icao1").Return("KJFK", "LAX", nil)

	// Create a new service with the mock clients
	service := NewService(mockOpenSkyClient, mockFlightAwareClient)

	// Call the method to test
	flights, err := service.GetFlightsInRadius(40.7128, -74.0060, 100)

	// Assert the results
	assert.NoError(t, err)
	assert.NotNil(t, flights)
	assert.Len(t, flights, 1)
	assert.Equal(t, "icao1", flights[0].Icao24)
	assert.Equal(t, "flight1", flights[0].Callsign)
	assert.Equal(t, "KJFK", flights[0].Origin)
	assert.Equal(t, "LAX", flights[0].Destination)

	// Assert that the mocks were called
	mockOpenSkyClient.AssertExpectations(t)
	mockFlightAwareClient.AssertExpectations(t)
}
