package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io" // Use io.ReadAll instead of ioutil.ReadAll
	"log"
	"net/http"
	"time"

	"github.com/carlo-colombo/sopra/config"
	"github.com/carlo-colombo/sopra/database"
)

const climatiqAPIURL = "https://beta4.api.climatiq.io/estimate"
const co2CachePrefix = "climatiq_co2_"

// ClimatiqAPIClient defines the interface for the Climatiq API client.
type ClimatiqAPIClient interface {
	GetFlightEmission(aircraftType string, distanceKm float64) (float64, error)
}

// ClimatiqClient is a client for the Climatiq API.
type ClimatiqClient struct {
	apiKey string
	client *http.Client
	db     *database.DB // For caching
}

// NewClimatiqClient creates a new ClimatiqClient.
func NewClimatiqClient(cfg *config.Config, db *database.DB) *ClimatiqClient {
	return &ClimatiqClient{
		apiKey: cfg.Climatiq.APIKey,
		client: &http.Client{Timeout: 10 * time.Second},
		db:     db,
	}
}

// ClimatiqRequest represents the request body for the Climatiq API.
type ClimatiqRequest struct {
	EmissionFactor EmissionFactor `json:"emission_factor"`
	Parameters     Parameters     `json:"parameters"`
}

// EmissionFactor represents the emission factor details.
type EmissionFactor struct {
	ID string `json:"id"`
}

// Parameters represents the parameters for the emission calculation.
type Parameters struct {
	DistanceKm   float64 `json:"distance_km"`
	AircraftType string  `json:"aircraft_type"`
}

// ClimatiqResponse represents the response body from the Climatiq API.
type ClimatiqResponse struct {
	CO2e float64 `json:"co2e"`
	// Add other fields if needed for future use
}

// GetFlightEmission estimates the CO2 emissions for a flight using the Climatiq API.
// It uses caching to avoid repeated API calls for the same flight parameters.
func (c *ClimatiqClient) GetFlightEmission(aircraftType string, distanceKm float64) (float64, error) {
	if c.apiKey == "" {
		return 0, fmt.Errorf("Climatiq API key is not configured")
	}
	if distanceKm <= 0 {
		return 0, nil
	}

	cacheKey := fmt.Sprintf("%s%s_%f", co2CachePrefix, aircraftType, distanceKm)
	cachedCO2, err := c.db.Get(cacheKey)
	if err == nil && cachedCO2 != "" {
		var co2 float64
		if err := json.Unmarshal([]byte(cachedCO2), &co2); err == nil {
			log.Printf("Climatiq CO2 emission for aircraftType %s, distance %f km served from cache: %f kg", aircraftType, distanceKm, co2)
			return co2, nil
		}
		log.Printf("Failed to unmarshal cached Climatiq CO2 emission: %v", err)
		// Fall through to API call if cache unmarshal fails
	} else if err != nil {
		log.Printf("Error retrieving Climatiq CO2 emission from cache: %v", err)
	}

	// Prepare the request payload
	requestBody := ClimatiqRequest{
		EmissionFactor: EmissionFactor{ID: "passenger_flight-route_type_unknown-aircraft_type_unknown-lufthansa"}, // Using a generic one for now
		Parameters: Parameters{
			DistanceKm:   distanceKm,
			AircraftType: aircraftType,
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal Climatiq request: %w", err)
	}

	req, err := http.NewRequest("POST", climatiqAPIURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return 0, fmt.Errorf("failed to create Climatiq request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to make Climatiq API request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read Climatiq API response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("Climatiq API returned non-OK status: %d, response: %s", resp.StatusCode, body)
	}

	var climatiqResp ClimatiqResponse
	if err := json.Unmarshal(body, &climatiqResp); err != nil {
		return 0, fmt.Errorf("failed to unmarshal Climatiq API response: %w", err)
	}

	// Cache the response
	co2JSON, err := json.Marshal(climatiqResp.CO2e)
	if err != nil {
		log.Printf("Failed to marshal Climatiq CO2 for caching: %v", err)
	} else {
		if err := c.db.Set(cacheKey, string(co2JSON), 24*time.Hour); err != nil { // Cache for 24 hours
			log.Printf("Failed to cache Climatiq CO2 emission: %v", err)
		}
	}

	log.Printf("Climatiq CO2 emission for aircraftType %s, distance %f km: %f kg", aircraftType, distanceKm, climatiqResp.CO2e)
	return climatiqResp.CO2e, nil
}
