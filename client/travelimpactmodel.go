package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/carlo-colombo/sopra/config"
	"github.com/carlo-colombo/sopra/database"
	"github.com/carlo-colombo/sopra/model"
)

const googleTravelImpactModelAPIURL = "https://travelimpactmodel.googleapis.com/v1/flights:computeFlightEmissions"
const co2CachePrefix = "google_tim_co2_"

// TravelImpactModelAPIClient defines the interface for the Google Travel Impact Model API client.
type TravelImpactModelAPIClient interface {
	GetFlightEmission(flightInfo *model.FlightInfo) (float64, error)
}

// TravelImpactModelClient is a client for the Google Travel Impact Model API.
type TravelImpactModelClient struct {
	apiKey string
	client *http.Client
	db     *database.DB // For caching
}

// NewTravelImpactModelClient creates a new TravelImpactModelClient.
func NewTravelImpactModelClient(cfg *config.Config, db *database.DB) *TravelImpactModelClient {
	return &TravelImpactModelClient{
		apiKey: cfg.TravelImpactModel.APIKey,
		client: &http.Client{Timeout: 30 * time.Second}, // Increased timeout for external API
		db:     db,
	}
}

// Date represents a whole or partial calendar date.
type Date struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
}

// Flight represents a single request item for direct flight emission estimates.
type Flight struct {
	Origin               string `json:"origin"`
	Destination          string `json:"destination"`
	DepartureDate        Date   `json:"departureDate"`
	FlightNumber         int    `json:"flightNumber"`
	OperatingCarrierCode string `json:"operatingCarrierCode"`
	// CabinClass string `json:"cabinClass"` // Removed as per API documentation for computeFlightEmissions
}

// ComputeFlightEmissionsRequest input definition.
type ComputeFlightEmissionsRequest struct {
	Flights []Flight `json:"flights"`
}

// EmissionsGramsPerPax grouped emissions per seating class results.
type EmissionsGramsPerPax struct {
	Economy        int `json:"economy"`
	PremiumEconomy int `json:"premiumEconomy"`
	Business       int `json:"business"`
	First          int `json:"first"`
}

// FlightWithEmissions direct flight with emission estimates.
type FlightWithEmissions struct {
	Flight               Flight               `json:"flight"`
	EmissionsGramsPerPax EmissionsGramsPerPax `json:"emissionsGramsPerPax"`
	Source               string               `json:"source"`
}

// ComputeFlightEmissionsResponse output definition.
type ComputeFlightEmissionsResponse struct {
	FlightEmissions []FlightWithEmissions `json:"flightEmissions"`
	ModelVersion    ModelVersion          `json:"modelVersion"`
}

// ModelVersion Travel Impact Model version.
type ModelVersion struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
	Patch int `json:"patch"`
}

// GetFlightEmission estimates the CO2 emissions for a flight using the Google Travel Impact Model API.
// It uses caching to avoid repeated API calls for the same flight parameters.
func (c *TravelImpactModelClient) GetFlightEmission(flightInfo *model.FlightInfo) (float64, error) {
	if c.apiKey == "" {
		return 0, fmt.Errorf("Google Travel Impact Model API key is not configured")
	}

	carrierCode := flightInfo.OperatorIata
	if carrierCode == "" {
		carrierCode = flightInfo.OperatorIcao
	}

	// Required fields for the Google TIM API
	if flightInfo.Origin.CodeIata == "" || flightInfo.Destination.CodeIata == "" ||
		flightInfo.FlightNumber == "" || carrierCode == "" ||
		flightInfo.ScheduledOut == nil {
		missingFields := []string{}
		if flightInfo.Origin.CodeIata == "" {
			missingFields = append(missingFields, "Origin.CodeIata")
		}
		if flightInfo.Destination.CodeIata == "" {
			missingFields = append(missingFields, "Destination.CodeIata")
		}
		if flightInfo.FlightNumber == "" {
			missingFields = append(missingFields, "FlightNumber")
		}
		if carrierCode == "" {
			missingFields = append(missingFields, "OperatingCarrierCode (Iata or Icao)")
		}
		if flightInfo.ScheduledOut == nil {
			missingFields = append(missingFields, "ScheduledOut")
		}
		log.Printf("Missing required flight information for Google TIM API for flight %s. Missing fields: %v. FlightInfo: %+v", flightInfo.Ident, missingFields, flightInfo)
		return 0, fmt.Errorf("missing required flight information for Google Travel Impact Model API")
	}

	departureDate := flightInfo.ScheduledOut.Format("2006-01-02")
	cacheKey := fmt.Sprintf("%s%s_%s_%s_%s_%s_%s",
		co2CachePrefix,
		flightInfo.Origin.CodeIata,
		flightInfo.Destination.CodeIata,
		departureDate,
		flightInfo.FlightNumber,
		carrierCode,
		"ECONOMY", // Assuming economy class for now
	)

	cachedCO2, err := c.db.Get(cacheKey)
	if err == nil && cachedCO2 != "" {
		var co2 float64
		if err := json.Unmarshal([]byte(cachedCO2), &co2); err == nil {
			log.Printf("Google TIM CO2 emission for flight %s served from cache: %f kg", flightInfo.Ident, co2)
			return co2, nil
		}
		log.Printf("Failed to unmarshal cached Google TIM CO2 emission: %v", err)
		// Fall through to API call if cache unmarshal fails
	} else if err != nil {
		log.Printf("Error retrieving Google TIM CO2 emission from cache: %v", err)
	}

	flightNumberInt, err := strconv.Atoi(flightInfo.FlightNumber)
	if err != nil {
		return 0, fmt.Errorf("failed to parse flight number '%s': %w", flightInfo.FlightNumber, err)
	}

	requestBody := ComputeFlightEmissionsRequest{
		Flights: []Flight{
			{
				Origin:      flightInfo.Origin.CodeIata,
				Destination: flightInfo.Destination.CodeIata,
				DepartureDate: Date{
					Year:  flightInfo.ScheduledOut.Year(),
					Month: int(flightInfo.ScheduledOut.Month()),
					Day:   flightInfo.ScheduledOut.Day(),
				},
				FlightNumber:         flightNumberInt,
				OperatingCarrierCode: carrierCode,
			},
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal Google TIM request: %w", err)
	}

	req, err := http.NewRequest("POST", googleTravelImpactModelAPIURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return 0, fmt.Errorf("failed to create Google TIM request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Goog-Api-Key", c.apiKey) // Google API key in header

	resp, err := c.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to make Google TIM API request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read Google TIM API response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("Google TIM API returned non-OK status: %d, response: %s", resp.StatusCode, body)
	}

	var timResp ComputeFlightEmissionsResponse
	if err := json.Unmarshal(body, &timResp); err != nil {
		return 0, fmt.Errorf("failed to unmarshal Google TIM API response: %w", err)
	}

	if len(timResp.FlightEmissions) == 0 {
		return 0, fmt.Errorf("no flight emissions found in Google TIM response for flight %s", flightInfo.Ident)
	}

	// Calculate the average CO2 for passenger across the cabin classes
	emissions := timResp.FlightEmissions[0].EmissionsGramsPerPax
	totalEmissions := float64(emissions.Economy + emissions.PremiumEconomy + emissions.Business + emissions.First)
	count := 0
	if emissions.Economy > 0 {
		count++
	}
	if emissions.PremiumEconomy > 0 {
		count++
	}
	if emissions.Business > 0 {
		count++
	}
	if emissions.First > 0 {
		count++
	}

	var co2Grams float64
	if count > 0 {
		co2Grams = totalEmissions / float64(count)
	} else {
		co2Grams = 0.0 // No emissions data available
	}
	co2Kg := co2Grams / 1000.0 // Convert grams to kilograms

	// Cache the response
	co2JSON, err := json.Marshal(co2Kg)
	if err != nil {
		log.Printf("Failed to marshal Google TIM CO2 for caching: %v", err)
	} else {
		if err := c.db.Set(cacheKey, string(co2JSON), 24*time.Hour); err != nil { // Cache for 24 hours
			log.Printf("Failed to cache Google TIM CO2 emission: %v", err)
		}
	}

	log.Printf("Google TIM CO2 emission for flight %s (Origin: %s, Dest: %s, Flight: %s, Carrier: %s, Date: %s): %f kg",
		flightInfo.Ident,
		flightInfo.Origin.CodeIata,
		flightInfo.Destination.CodeIata,
		flightInfo.FlightNumber,
		flightInfo.OperatorIata, // Keep this for original FlightInfo context
		departureDate,
		co2Kg)
	return co2Kg, nil
}
