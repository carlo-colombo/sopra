package service

import (
	"encoding/json"
	"log"
	"time"

	"github.com/carlo-colombo/sopra/config"
	"github.com/carlo-colombo/sopra/database"
	"github.com/carlo-colombo/sopra/haversine"
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

// ClimatiqAPIClient defines the interface for the Climatiq API client.
type ClimatiqAPIClient interface {
	GetFlightEmission(aircraftType string, distanceKm float64) (float64, error)
}

// Service is the main service for the application.
type Service struct {
	openskyClient     OpenSkyAPIClient
	flightawareClient FlightAwareAPIClient
	climatiqClient    ClimatiqAPIClient // Add Climatiq client
	db                *database.DB
	cfg               *config.Config // Add config to the service struct
}

// NewService creates a new Service.
func NewService(openskyClient OpenSkyAPIClient, flightawareClient FlightAwareAPIClient, climatiqClient ClimatiqAPIClient, db *database.DB, cfg *config.Config) *Service {
	return &Service{
		openskyClient:     openskyClient,
		flightawareClient: flightawareClient,
		climatiqClient:    climatiqClient, // Store the Climatiq client
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
			flightInfo.Distance = haversine.Distance(lat, lon, flight.Latitude, flight.Longitude) * 1000

			// --- START Climatiq Integration ---
			// Convert nautical miles to kilometers for Climatiq API (1 NM = 1.852 km)
			distanceKm := float64(flightInfo.RouteDistance) * 1.852
			co2, err := s.climatiqClient.GetFlightEmission(flightInfo.AircraftType, distanceKm)
			if err != nil {
				log.Printf("Error getting CO2 emission from Climatiq for aircraft %s, distance %.2f km: %v", flightInfo.AircraftType, distanceKm, err)
				flightInfo.CO2KG = 0.0 // Set to 0 or handle as appropriate
			} else {
				flightInfo.CO2KG = co2
			}
			// --- END Climatiq Integration ---

			if flightInfo.OperatorIcao != "" {
				_, err := s.getOperatorInfo(flightInfo.OperatorIcao)
				if err != nil {
					log.Printf("Could not get operator info for ICAO %s: %v", flightInfo.OperatorIcao, err)
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

// EstimateCO2 estimates the CO2 emissions of a flight in kilograms.
// This function is now deprecated and will use the Climatiq API through the service struct.
func EstimateCO2(aircraftType string, distanceNm int) float64 {
	// This function is no longer used directly.
	// The CO2 estimation is now handled by the climatiqClient within the GetFlightsInRadius method.
	log.Println("Warning: Deprecated EstimateCO2 function called. Use service.climatiqClient.GetFlightEmission instead.")
	return 0.0 // Return 0 as this function is deprecated.
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
