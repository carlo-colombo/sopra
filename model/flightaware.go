package model

import "time"

// FlightAwareResponse represents the top-level structure of the FlightAware AeroAPI response.
type FlightAwareResponse struct {
	Flights []FlightInfo `json:"flights"`
}

// FlightInfo represents detailed information about a flight from FlightAware.
type FlightInfo struct {
	IdentificationCount int         `json:"identification_count"`
	Ident          string      `json:"ident"`
	FaFlightID     string      `json:"fa_flight_id"`
	Operator       string      `json:"operator"`
	OperatorICAO   string      `json:"operator_icao"`
	FlightNumber   string      `json:"flight_number"`
	AircraftType   string      `json:"aircraft_type"`
	AircraftTypeFA string      `json:"aircraft_type_fa"`
	Origin         AirportDetail `json:"origin"`
	Destination    AirportDetail `json:"destination"`
	ScheduledOut   time.Time   `json:"scheduled_out"`
	EstimatedOut   time.Time   `json:"estimated_out"`
	ActualOut      *time.Time  `json:"actual_out"` // Pointer because it can be null
	ScheduledOn    time.Time   `json:"scheduled_on"`
	EstimatedOn    time.Time   `json:"estimated_on"`
	ActualOn       *time.Time  `json:"actual_on"`
	ScheduledIn    time.Time   `json:"scheduled_in"`
	EstimatedIn    time.Time   `json:"estimated_in"`
	ActualIn       *time.Time  `json:"actual_in"`
	Status         string      `json:"status"`
	ProgressPercent int         `json:"progress_percent"`
	FiledEte       int         `json:"filed_ete"`
	Route          string      `json:"route"`
	Altitude       int         `json:"altitude"`
	Groundspeed    int         `json:"groundspeed"`
	TrueAirspeed   int         `json:"true_airspeed"`
	Heading        int         `json:"heading"`
	Latitude       float64     `json:"latitude"`
	Longitude      float64     `json:"longitude"`
	FlightRules    string      `json:"flight_rules"`
	DelayReason    *string     `json:"delay_reason"` // Pointer because it can be null
	BaggageClaim   string      `json:"baggage_claim"`
	GateOut        string      `json:"gate_out"`
	GateIn         string      `json:"gate_in"`
	TerminalOut    string      `json:"terminal_out"`
	TerminalIn     string      `json:"terminal_in"`
	Codeshares     []string `json:"codeshares"`
}

// AirportDetail represents detailed information about an airport.
type AirportDetail struct {
	CodeIATA    string `json:"code_iata"`
	AirportName string `json:"airport_name"`
	City        string `json:"city"`
	State       string `json:"state"`
	Country     string `json:"country"`
	Timezone    string `json:"timezone"`
}