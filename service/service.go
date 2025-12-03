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

// FlightAwareAPIClient defines the interface for the FlightAware AeroAPI client.
type FlightAwareAPIClient interface {
	GetFlightInfo(icao24 string) (origin, destination string, err error)
}

// Service is the main service for the application.
type Service struct {
	openskyClient     OpenSkyAPIClient
	flightawareClient FlightAwareAPIClient
}

// NewService creates a new Service.
func NewService(openskyClient OpenSkyAPIClient, flightawareClient FlightAwareAPIClient) *Service {
	return &Service{
		openskyClient:     openskyClient,
		flightawareClient: flightawareClient,
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

	flights := states.ToFlights()

	for i := range flights {
		origin, destination, err := s.flightawareClient.GetFlightInfo(flights[i].Icao24)
		if err != nil {
			log.Printf("Could not get FlightAware info for ICAO24 %s: %v", flights[i].Icao24, err)
			// Continue even if FlightAware lookup fails for one flight
			continue
		}
		flights[i].Origin = origin
		flights[i].Destination = destination
	}

	return flights, nil
}
