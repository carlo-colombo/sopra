package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// FlightAwareClient is a client for the FlightAware AeroAPI.
type FlightAwareClient struct {
	httpClient *http.Client
	apiKey     string
	baseURL    string
}

// NewFlightAwareClient creates a new FlightAwareClient.
func NewFlightAwareClient(apiKey string) *FlightAwareClient {
	return &FlightAwareClient{
		httpClient: &http.Client{Timeout: 10 * time.Second}, // Add a timeout for HTTP requests
		apiKey:     apiKey,
		baseURL:    "https://aeroapi.flightaware.com/aeroapi",
	}
}

// AeroAPIResponse represents the simplified structure of a FlightAware AeroAPI response for a flight.
// This is a placeholder and might need adjustment based on actual API response.
type AeroAPIResponse struct {
	Flights []struct {
		Origin struct {
			AirportCode string `json:"code_icao"`
		} `json:"origin"`
		Destination struct {
			AirportCode string `json:"code_icao"`
		} `json:"destination"`
	} `json:"flights"`
}

// GetFlightInfo retrieves flight information (origin and destination) from FlightAware AeroAPI.
func (c *FlightAwareClient) GetFlightInfo(icao24 string) (origin, destination string, err error) {
	// Construct the URL to get flight info by icao24. This is an assumption.
	// The actual API might require a callsign or a different identifier.
	// Assuming an endpoint like /flights/{icao24} or similar that returns flight details.
	url := fmt.Sprintf("%s/aircraft/%s/flights", c.baseURL, icao24)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("x-apikey", c.apiKey)
	req.Header.Add("Accept", "application/json; charset=UTF-8")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to make request to FlightAware API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("FlightAware API returned non-OK status: %s", resp.Status)
	}

	var aeroAPIResponse AeroAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&aeroAPIResponse); err != nil {
		return "", "", fmt.Errorf("failed to decode FlightAware API response: %w", err)
	}

	if len(aeroAPIResponse.Flights) > 0 {
		return aeroAPIResponse.Flights[0].Origin.AirportCode, aeroAPIResponse.Flights[0].Destination.AirportCode, nil
	}

	return "", "", fmt.Errorf("no flight information found for icao24: %s", icao24)
}
