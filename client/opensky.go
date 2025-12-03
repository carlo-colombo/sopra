package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/carlo-colombo/sopra/model"
	"golang.org/x/oauth2/clientcredentials"
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
