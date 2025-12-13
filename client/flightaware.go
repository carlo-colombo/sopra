package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/carlo-colombo/sopra/database"
	"github.com/carlo-colombo/sopra/model"
)

// FlightAwareClient is a client for the FlightAware AeroAPI.
type FlightAwareClient struct {
	httpClient *http.Client
	apiKey     string
	baseURL    string
	db         *database.DB
}

// NewFlightAwareClient creates a new FlightAwareClient.
func NewFlightAwareClient(apiKey string, db *database.DB) *FlightAwareClient {
	return &FlightAwareClient{
		httpClient: &http.Client{Timeout: 10 * time.Second}, // Add a timeout for HTTP requests
		apiKey:     apiKey,
		baseURL:    "https://aeroapi.flightaware.com/aeroapi",
		db:         db,
	}
}

// GetFlightInfo retrieves detailed flight information from FlightAware AeroAPI by its ident (callsign).
func (c *FlightAwareClient) GetFlightInfo(ident string) (*model.FlightInfo, error) {
	// Try to get the flight info from the cache first.
	if cachedFlightInfo, _, err := c.db.GetFlight(ident); err == nil && cachedFlightInfo != nil {
		// Update last_seen even if retrieved from cache
		if err := c.db.LogFlight(ident, cachedFlightInfo); err != nil {
			log.Printf("Failed to update last_seen for cached flight %s: %v", ident, err)
		}
		return cachedFlightInfo, nil
	}

	url := fmt.Sprintf("%s/flights/%s", c.baseURL, ident)
	log.Printf("Requesting flight info from FlightAware API: %s\n", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("x-apikey", c.apiKey)
	req.Header.Add("Accept", "application/json; charset=UTF-8")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to FlightAware API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil // No flight found for the given ident, not an error
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("FlightAware API returned non-OK status: %s", resp.Status)
	}

	var faResponse model.FlightAwareResponse
	if err := json.NewDecoder(resp.Body).Decode(&faResponse); err != nil {
		return nil, fmt.Errorf("failed to decode FlightAware API response: %w", err)
	}

	if len(faResponse.Flights) > 0 {
		flightInfo := &faResponse.Flights[0]
		// Cache the result
		if err := c.db.LogFlight(ident, flightInfo); err != nil {
			log.Printf("Failed to cache flight info for ident %s: %v", ident, err)
		}
		return flightInfo, nil
	}

	return nil, nil // No flight info in the response
}

// GetOperator retrieves operator information from FlightAware AeroAPI by its ICAO code.
func (c *FlightAwareClient) GetOperator(icao string) (string, error) {
	url := fmt.Sprintf("%s/operators/%s", c.baseURL, icao)
	log.Printf("Requesting operator info from FlightAware API: %s\n", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("x-apikey", c.apiKey)
	req.Header.Add("Accept", "application/json; charset=UTF-8")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request to FlightAware API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", nil // No operator found for the given ICAO, not an error
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("FlightAware API returned non-OK status: %s", resp.Status)
	}

	var raw json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return "", fmt.Errorf("failed to decode FlightAware API response: %w", err)
	}

	return string(raw), nil
}
