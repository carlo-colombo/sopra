package service

import (
	"log"

	"github.com/carlo-colombo/sopra/haversine"
	"github.com/carlo-colombo/sopra/model"
)

// OpenSkyAPIClient defines the interface for the OpenSky API client.
type OpenSkyAPIClient interface {
	GetStatesWithBoundingBox(lamin, lomin, lamax, lomax float64) (*model.States, error)
}

// Service is the main service for the application.
type Service struct {
	openskyClient OpenSkyAPIClient
}

// NewService creates a new Service.
func NewService(openskyClient OpenSkyAPIClient) *Service {
	return &Service{
		openskyClient: openskyClient,
	}
}

// GetFlightsInRadius returns a list of flights within a given radius from a location.
func (s *Service) GetFlightsInRadius(lat, lon, radius float64) ([]model.Flight, error) {
	log.Printf("Request for flights in radius %f from position (%f, %f)\n", radius, lat, lon)
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
