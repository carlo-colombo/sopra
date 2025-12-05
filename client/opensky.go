package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/carlo-colombo/sopra/model"
	"golang.org/x/oauth2/clientcredentials"
	"github.com/carlo-colombo/sopra/haversine"
)

// OpenSkyClient is a client for the OpenSky Network API.
type OpenSkyClient struct {
	httpClient *http.Client
	baseURL    string
}

// NewOpenSkyClient creates a new OpenSkyClient.
func NewOpenSkyClient(clientID, clientSecret string) *OpenSkyClient {
	config := clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     "https://auth.opensky-network.org/auth/realms/opensky-network/protocol/openid-connect/token",
	}
	ctx := context.Background()
	httpClient := config.Client(ctx)

	return &OpenSkyClient{
		httpClient: httpClient,
		baseURL:    "https://opensky-network.org/api",
	}
}

// GetStates retrieves all flight states from the OpenSky Network API.
func (c *OpenSkyClient) GetStates() (*model.States, error) {
	log.Printf("Requesting all states from OpenSky API: %s/states/all\n", c.baseURL)
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/states/all", c.baseURL))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get states: %s", resp.Status)
	}

	var states model.States
	if err := json.NewDecoder(resp.Body).Decode(&states); err != nil {
		return nil, err
	}

	return &states, nil
}

// GetStatesWithBoundingBox retrieves flight states within a specified bounding box from the OpenSky Network API.
func (c *OpenSkyClient) GetStatesWithBoundingBox(lamin, lomin, lamax, lomax float64) (*model.States, error) {
	url := fmt.Sprintf("%s/states/all?lamin=%f&lomin=%f&lamax=%f&lomax=%f", c.baseURL, lamin, lomin, lamax, lomax)
	log.Printf("Requesting states from OpenSky API: %s\n", url)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get states: %s", resp.Status)
	}

	var states model.States
	if err := json.NewDecoder(resp.Body).Decode(&states); err != nil {
		return nil, err
	}

	return &states, nil
}

// GetStatesInRadius retrieves flight states within a specified radius from a given central point.
func (c *OpenSkyClient) GetStatesInRadius(lat, lon, radiusKm float64) ([]model.Flight, error) {
	bbox := haversine.GetBoundingBox(lat, lon, radiusKm)
	states, err := c.GetStatesWithBoundingBox(bbox.MinLat, bbox.MinLon, bbox.MaxLat, bbox.MaxLon)
	if err != nil {
		return nil, err
	}

	var filteredFlights []model.Flight
	for _, flight := range states.ToFlights() {
		if flight.Latitude != 0 && flight.Longitude != 0 &&
			haversine.Distance(lat, lon, flight.Latitude, flight.Longitude) <= radiusKm {
			filteredFlights = append(filteredFlights, flight)
		}
	}
	return filteredFlights, nil
}

