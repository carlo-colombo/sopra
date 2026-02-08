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

// TravelImpactModelAPIClient defines the interface for the Google Travel Impact Model API client.
type TravelImpactModelAPIClient interface {
	GetFlightEmission(flightInfo *model.FlightInfo) (float64, error)
}

// Service is the main service for the application.
type Service struct {
	openskyClient           OpenSkyAPIClient
	flightawareClient       FlightAwareAPIClient
	travelImpactModelClient TravelImpactModelAPIClient // Add Travel Impact Model client
	db                      *database.DB
	cfg                     *config.Config // Add config to the service struct
}

// NewService creates a new Service.
func NewService(openskyClient OpenSkyAPIClient, flightawareClient FlightAwareAPIClient, travelImpactModelClient TravelImpactModelAPIClient, db *database.DB, cfg *config.Config) *Service {
	return &Service{
		openskyClient:           openskyClient,
		flightawareClient:       flightawareClient,
		travelImpactModelClient: travelImpactModelClient, // Store the Travel Impact Model client
		db:                      db,
		cfg:                     cfg, // Store the config
	}
}

// GetFlightsInRadius returns a list of enriched FlightInfo objects within a given radius from a location.
func (s *Service) GetFlightsInRadius(lat, lon, radius float64) ([]model.FlightInfo, error) {
	log.Printf("Request for flights in radius %f from position (%f, %f)\n", radius, lat, lon)
	startOpenSky := time.Now()
	openskyFlights, err := s.openskyClient.GetStatesInRadius(lat, lon, radius)
	if err != nil {
		return nil, err
	}
	log.Printf("OpenSky API call took %s to get %d flights\n", time.Since(startOpenSky), len(openskyFlights))

	var enrichedFlights []model.FlightInfo
	for _, flight := range openskyFlights {
		if flight.Callsign == "" {
			continue // Skip flights without a callsign for FlightAware lookup
		}

		startFlightAwareInfo := time.Now()
		flightInfo, err := s.flightawareClient.GetFlightInfo(flight.Callsign)
		if err != nil {
			log.Printf("Could not get FlightAware info for callsign %s (ICAO24: %s): %v. Took %s\n", flight.Callsign, flight.Icao24, err, time.Since(startFlightAwareInfo))
			continue // Continue even if FlightAware lookup fails for one flight
		}
		log.Printf("FlightAware GetFlightInfo for %s took %s\n", flight.Callsign, time.Since(startFlightAwareInfo))

		if flightInfo != nil {
			flightInfo.Latitude = flight.Latitude
			flightInfo.Longitude = flight.Longitude
			flightInfo.Distance = haversine.Distance(lat, lon, flight.Latitude, flight.Longitude) * 1000

			// --- START Google Travel Impact Model Integration ---
			startTIM := time.Now()
			co2, err := s.travelImpactModelClient.GetFlightEmission(flightInfo)
			if err != nil {
				log.Printf("Error getting CO2 emission from Google Travel Impact Model for flight %s: %v. Took %s\n", flightInfo.Ident, err, time.Since(startTIM))
				flightInfo.CO2KG = 0.0 // Set to 0 or handle as appropriate
			} else {
				flightInfo.CO2KG = co2
			}
			log.Printf("Google Travel Impact Model API call took %s\n", time.Since(startTIM))
			// --- END Google Travel Impact Model Integration ---

			if flightInfo.OperatorIcao != "" {
				startOperatorInfo := time.Now()
				_, err := s.getOperatorInfo(flightInfo.OperatorIcao)
				if err != nil {
					log.Printf("Could not get operator info for ICAO %s: %v. Took %s\n", flightInfo.OperatorIcao, err, time.Since(startOperatorInfo))
				}
				log.Printf("GetOperatorInfo for %s took %s\n", flightInfo.OperatorIcao, time.Since(startOperatorInfo))
			}
			enrichedFlights = append(enrichedFlights, *flightInfo)
		}
	}
	return enrichedFlights, nil
}

func (s *Service) getOperatorInfo(icao string) (*model.OperatorInfo, error) {
	startDbGet := time.Now()
	cachedOperator, err := s.db.GetOperator(icao)
	if err != nil {
		log.Printf("Error getting operator %s from DB: %v. Took %s\n", icao, err, time.Since(startDbGet))
		return nil, err
	}
	log.Printf("DB GetOperator for %s took %s\n", icao, time.Since(startDbGet))

	if cachedOperator != "" {
		var operatorInfo model.OperatorInfo
		if err := json.Unmarshal([]byte(cachedOperator), &operatorInfo); err != nil {
			return nil, err
		}
		caser := cases.Title(language.English)
		operatorInfo.Shortname = caser.String(operatorInfo.Shortname)
		return &operatorInfo, nil
	}

	startFlightAwareOperator := time.Now()
	operatorJSON, err := s.flightawareClient.GetOperator(icao)
	if err != nil {
		log.Printf("Error getting operator %s from FlightAware: %v. Took %s\n", icao, err, time.Since(startFlightAwareOperator))
		return nil, err
	}
	log.Printf("FlightAware GetOperator for %s took %s\n", icao, time.Since(startFlightAwareOperator))

	if operatorJSON == "" {
		return nil, nil
	}

	startDbLog := time.Now()
	if err := s.db.LogOperator(icao, operatorJSON); err != nil {
		log.Printf("Failed to cache operator info for ICAO %s: %v. Took %s\n", icao, err, time.Since(startDbLog))
	} else {
		log.Printf("DB LogOperator for %s took %s\n", icao, time.Since(startDbLog))
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
