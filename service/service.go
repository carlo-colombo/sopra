package service

import (
	"log"

	"github.com/carlo-colombo/sopra/model"
)

// OpenSkyAPIClient defines the interface for the OpenSky API client.
type OpenSkyAPIClient interface {
	GetStatesInRadius(lat, lon, radiusKm float64) ([]model.Flight, error)
}

// FlightAwareAPIClient defines the interface for the FlightAware AeroAPI client.
type FlightAwareAPIClient interface {
	GetFlightInfo(ident string) (*model.FlightInfo, error)
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

// GetFlightsInRadius returns a list of enriched FlightInfo objects within a given radius from a location.
func (s *Service) GetFlightsInRadius(lat, lon, radius float64) ([]model.FlightInfo, error) {
	log.Printf("Request for flights in radius %f from position (%f, %f)\n", radius, lat, lon)

	openskyFlights, err := s.openskyClient.GetStatesInRadius(lat, lon, radius)
	if err != nil {
		return nil, err
	}

	var enrichedFlights []model.FlightInfo
	for _, flight := range openskyFlights {
		if flight.Callsign == "" {
			continue // Skip flights without a callsign for FlightAware lookup
		}

		flightInfo, err := s.flightawareClient.GetFlightInfo(flight.Callsign)
		if err != nil {
			log.Printf("Could not get FlightAware info for callsign %s (ICAO24: %s): %v", flight.Callsign, flight.Icao24, err)
			continue // Continue even if FlightAware lookup fails for one flight
		}
		if flightInfo != nil {
			enrichedFlights = append(enrichedFlights, *flightInfo)
		}
	}

	return enrichedFlights, nil
}
