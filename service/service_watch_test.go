package service_test

import (
	"log"
	"os"
	"sync" // Added sync import
	"testing"
	"time"

	"github.com/carlo-colombo/sopra/config"
	"github.com/carlo-colombo/sopra/database"
	"github.com/carlo-colombo/sopra/model"
	"github.com/carlo-colombo/sopra/service" // Import the service package
)

// MockOpenSkyClient is a mock implementation of OpenSkyAPIClient for testing.
type MockOpenSkyClient struct {
	mu              sync.Mutex
	GetStatesCalls  int
	FlightsToReturn []model.Flight // Flights to return on GetStatesInRadius call
	ErrToReturn     error          // Error to return on GetStatesInRadius call
}

// GetStatesInRadius increments the call counter and returns predefined flights or an error.
func (m *MockOpenSkyClient) GetStatesInRadius(lat, lon, radiusKm float64) ([]model.Flight, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.GetStatesCalls++
	return m.FlightsToReturn, m.ErrToReturn
}

// MockFlightAwareClient is a mock implementation of FlightAwareAPIClient for testing.
type MockFlightAwareClient struct {
	mu             sync.Mutex
	GetFlightCalls int
	FlightToReturn *model.FlightInfo // FlightInfo to return on GetFlightInfo call
	ErrToReturn    error             // Error to return on GetFlightInfo call
}

// GetFlightInfo increments the call counter and returns predefined flight info or an error.
func (m *MockFlightAwareClient) GetFlightInfo(ident string) (*model.FlightInfo, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.GetFlightCalls++
	return m.FlightToReturn, m.ErrToReturn
}

// Ensure that logging doesn't print during tests
func TestMain(m *testing.M) {
	log.SetOutput(os.Stderr) // Or os.Stdout, or ioutil.Discard
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestService_RunWatchMode(t *testing.T) {
	// 1. Setup a mock OpenSky client
	mockOpenSky := &MockOpenSkyClient{}
	mockOpenSky.FlightsToReturn = []model.Flight{
		{Icao24: "a8093f", Callsign: "UAL123"},
	}

	// 2. Setup a mock FlightAware client
	mockFlightAware := &MockFlightAwareClient{}
	mockFlightAware.FlightToReturn = &model.FlightInfo{
		Ident:      "UAL123",
		FaFlightID: "UAL123-12345",
	}

	// 3. Create a temporary database for testing
	tempDBPath := "test_sopra.db"
	db, err := database.NewDB(tempDBPath)
	if err != nil {
		t.Fatalf("Failed to create temporary database: %v", err)
	}
	defer func() {
		db.Close()
		os.Remove(tempDBPath)
	}()

	// 4. Configure a test config
	testCfg := &config.Config{
		Service: struct {
			Latitude  float64 `mapstructure:"latitude"`
			Longitude float64 `mapstructure:"longitude"`
			Radius    float64 `mapstructure:"radius"`
		}{
			Latitude:  34.052235,
			Longitude: -118.243683,
			Radius:    100,
		},
		OpenSkyClient: struct {
			ID     string `mapstructure:"id"`
			Secret string `mapstructure:"secret"`
		}{
			ID:     "test",
			Secret: "test",
		},
		FlightAware: struct {
			APIKey string `mapstructure:"api_key"`
		}{
			APIKey: "test",
		},
		Watch:    true,
		Interval: 1, // Short interval for testing
	}

	// 5. Initialize the service with mocks and test config
	appService := service.NewService(mockOpenSky, mockFlightAware, db, testCfg)

	// 6. Run RunWatchMode in a goroutine
	done := make(chan struct{})
	go func() {
		appService.RunWatchMode(testCfg.Interval)
		close(done)
	}()

	// 7. Wait for a few intervals
	time.Sleep(2500 * time.Millisecond) // Wait for 2-3 ticks (interval is 1 second)

	// To stop the RunWatchMode gracefully in test environment we need to close done channel or implement a stop chan
	// For now, we will simply not wait for `done` and let the test finish.
	// In a real application, a context.Context with cancellation would be used to stop the watcher.

	// 8. Assert that GetStatesInRadius was called multiple times
	expectedCalls := 2 // At least 2 calls for 2.5 seconds with 1 second interval
	if mockOpenSky.GetStatesCalls < expectedCalls {
		t.Errorf("Expected at least %d calls to GetStatesInRadius, but got %d", expectedCalls, mockOpenSky.GetStatesCalls)
	}
}
