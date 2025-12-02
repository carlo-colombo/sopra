package service

import (
	"sopra/client"
	"sopra/haversine"
	"sopra/model"
)

// Service is the main service for the application.
type Service struct {
	openskyClient *client.OpenSkyClient
}

// NewService creates a new Service.
func NewService(openskyClient *client.OpenSkyClient) *Service {
	return &Service{
		openskyClient: openskyClient,
	}
}

// GetFlightsInRadius returns a list of flights within a given radius from a location.
func (s *Service) GetFlightsInRadius(lat, lon, radius float64) ([]model.Flight, error) {
	bbox := haversine.GetBoundingBox(lat, lon, radius)

	states, err := s.openskyClient.GetStatesWithBoundingBox(bbox.MinLat, bbox.MinLon, bbox.MaxLat, bbox.MaxLon)
	if err != nil {
		return nil, err
	}

	if states == nil {
		return []model.Flight{}, nil
	}

	return states.ToFlights(), nil
}
