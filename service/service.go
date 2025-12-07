package service

import (
	"log"
	"time"

	"github.com/carlo-colombo/sopra/config"
	"github.com/carlo-colombo/sopra/database"
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

	db                *database.DB

	cfg               *config.Config // Add config to the service struct

}



// NewService creates a new Service.

func NewService(openskyClient OpenSkyAPIClient, flightawareClient FlightAwareAPIClient, db *database.DB, cfg *config.Config) *Service {

	return &Service{

		openskyClient:     openskyClient,

		flightawareClient: flightawareClient,

		db:                db,

		cfg:               cfg, // Store the config

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
			flightInfo.Latitude = flight.Latitude
			flightInfo.Longitude = flight.Longitude
			enrichedFlights = append(enrichedFlights, *flightInfo)
		}

	}



	return enrichedFlights, nil

}



// LogFlights logs a slice of flights to the database.

func (s *Service) LogFlights(flights []model.FlightInfo) {

	for _, flight := range flights {

		err := s.db.LogFlight(flight.Ident, &flight)

		if err != nil {

			log.Printf("Error logging flight %s: %v", flight.Ident, err)

		}

	}

}



// RunWatchMode continuously fetches and logs flights at a specified interval.

func (s *Service) RunWatchMode(interval int) {

	ticker := time.NewTicker(time.Duration(interval) * time.Second)

	defer ticker.Stop()



	for range ticker.C {

		log.Println("Watching for flights...")

		flights, err := s.GetFlightsInRadius(s.cfg.Service.Latitude, s.cfg.Service.Longitude, s.cfg.Service.Radius)

		if err != nil {

			log.Printf("Error getting flights: %v", err)

			continue

		}

		s.LogFlights(flights)

	}

}
