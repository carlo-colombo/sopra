package service

import (
	"encoding/json"
	"log"
	"time"

	"github.com/carlo-colombo/sopra/config"
	"github.com/carlo-colombo/sopra/database"
	"github.com/carlo-colombo/sopra/model"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// OpenSkyAPIClient defines the interface for the OpenSky API client.
type OpenSkyAPIClient interface {
	GetStatesInRadius(lat, lon, radiusKm float64) ([]model.Flight, error)
}

// FlightAwareAPIClient defines the interface for the FlightAware AeroAPI client.
type FlightAwareAPIClient interface {
	GetFlightInfo(ident string) (*model.FlightInfo, error)
	GetOperator(icao string) (string, error)
}

// Service is the main service for the application.

type Service struct {
	openskyClient OpenSkyAPIClient

	flightawareClient FlightAwareAPIClient

	db *database.DB

	cfg *config.Config // Add config to the service struct

}

// NewService creates a new Service.

func NewService(openskyClient OpenSkyAPIClient, flightawareClient FlightAwareAPIClient, db *database.DB, cfg *config.Config) *Service {

	return &Service{

		openskyClient: openskyClient,

		flightawareClient: flightawareClient,

		db: db,

		cfg: cfg, // Store the config

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
			if flightInfo.OperatorIcao != "" {
				operatorInfo, err := s.getOperatorInfo(flightInfo.OperatorIcao)
				if err != nil {
					log.Printf("Could not get operator info for ICAO %s: %v", flightInfo.OperatorIcao, err)
				} else {
					flightInfo.OperatorInfo = *operatorInfo
				}
			}
			enrichedFlights = append(enrichedFlights, *flightInfo)
		}
	}
	return enrichedFlights, nil
}

func (s *Service) getOperatorInfo(icao string) (*model.OperatorInfo, error) {
	cachedOperator, err := s.db.GetOperator(icao)
	if err != nil {
		return nil, err
	}
	if cachedOperator != "" {
		var operatorInfo model.OperatorInfo
		if err := json.Unmarshal([]byte(cachedOperator), &operatorInfo); err != nil {
			return nil, err
		}
		caser := cases.Title(language.English)
		operatorInfo.Shortname = caser.String(operatorInfo.Shortname)
		return &operatorInfo, nil
	}
	operatorJSON, err := s.flightawareClient.GetOperator(icao)
	if err != nil {
		return nil, err
	}
	if operatorJSON == "" {
		return nil, nil
	}
	if err := s.db.LogOperator(icao, operatorJSON); err != nil {
		log.Printf("Failed to cache operator info for ICAO %s: %v", icao, err)
	}
	var operatorInfo model.OperatorInfo
	if err := json.Unmarshal([]byte(operatorJSON), &operatorInfo); err != nil {
		return nil, err
	}
	caser := cases.Title(language.English)
	operatorInfo.Shortname = caser.String(operatorInfo.Shortname)
	return &operatorInfo, nil
}

// LogFlights logs a slice of flights to the database.

func (s *Service) LogFlights(flights []model.FlightInfo) {

		for _, flight := range flights { // Changed back to use flight

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
