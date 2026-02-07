package model

import "time"

// FlightAwareResponse represents the top-level structure of the FlightAware AeroAPI response.
type FlightAwareResponse struct {
	Flights  []FlightInfo `json:"flights"`
	Links    *string      `json:"links"`
	NumPages int          `json:"num_pages"`
}

// FlightInfo represents detailed information about a flight from FlightAware.
type FlightInfo struct {
	Ident                         string        `json:"ident"`
	IdentIcao                     string        `json:"ident_icao"`
	IdentIata                     string        `json:"ident_iata"`
	ActualRunwayOff               string        `json:"actual_runway_off"`
	ActualRunwayOn                string        `json:"actual_runway_on"`
	FaFlightID                    string        `json:"fa_flight_id"`
	Operator                      string        `json:"operator"`
	OperatorIcao                  string        `json:"operator_icao"`
	OperatorIata                  string        `json:"operator_iata"`
	FlightNumber                  string        `json:"flight_number"`
	Registration                  string        `json:"registration"`
	AtcIdent                      *string       `json:"atc_ident"`
	InboundFaFlightID             string        `json:"inbound_fa_flight_id"`
	Codeshares                    []string      `json:"codeshares"`
	CodesharesIata                []string      `json:"codeshares_iata"`
	Blocked                       bool          `json:"blocked"`
	Diverted                      bool          `json:"diverted"`
	Cancelled                     bool          `json:"cancelled"`
	PositionOnly                  bool          `json:"position_only"`
	Origin                        AirportDetail `json:"origin"`
	Destination                   AirportDetail `json:"destination"`
	DepartureDelay                int           `json:"departure_delay"`
	ArrivalDelay                  int           `json:"arrival_delay"`
	FiledEte                      int           `json:"filed_ete"`
	ForesightPredictionsAvailable bool          `json:"foresight_predictions_available"`
	ScheduledOut                  *time.Time    `json:"scheduled_out"`
	EstimatedOut                  *time.Time    `json:"estimated_out"`
	ActualOut                     *time.Time    `json:"actual_out"`
	ScheduledOff                  time.Time     `json:"scheduled_off"`
	EstimatedOff                  time.Time     `json:"estimated_off"`
	ActualOff                     time.Time     `json:"actual_off"`
	ScheduledOn                   time.Time     `json:"scheduled_on"`
	EstimatedOn                   time.Time     `json:"estimated_on"`
	ActualOn                      time.Time     `json:"actual_on"`
	ScheduledIn                   *time.Time    `json:"scheduled_in"`
	EstimatedIn                   *time.Time    `json:"estimated_in"`
	ActualIn                      *time.Time    `json:"actual_in"`
	ProgressPercent               int           `json:"progress_percent"`
	Status                        string        `json:"status"`
	AircraftType                  string        `json:"aircraft_type"`
	RouteDistance                 int           `json:"route_distance"`
	FiledAirspeed                 *int          `json:"filed_airspeed"`
	FiledAltitude                 *int          `json:"filed_altitude"`
	Route                         *string       `json:"route"`
	BaggageClaim                  *string       `json:"baggage_claim"`
	SeatsCabinBusiness            *int          `json:"seats_cabin_business"`
	SeatsCabinCoach               *int          `json:"seats_cabin_coach"`
	SeatsCabinFirst               *int          `json:"seats_cabin_first"`
	GateOrigin                    *string       `json:"gate_origin"`
	GateDestination               *string       `json:"gate_destination"`
	TerminalOrigin                *string       `json:"terminal_origin"`
	TerminalDestination           *string       `json:"terminal_destination"`
	Type                          string        `json:"type"`
	Latitude                      float64       `json:"latitude"`
	Longitude                     float64       `json:"longitude"`
	Distance                      float64       `json:"distance_m"`
	IdentificationCount           int           `json:"-"`
}

// OperatorInfo represents detailed information about an operator.
type OperatorInfo struct {
	Name      string `json:"name"`
	Shortname string `json:"shortname"`
	Country   string `json:"country"`
}

// AirportDetail represents detailed information about an airport.
type AirportDetail struct {
	Code           string  `json:"code"`
	CodeIcao       string  `json:"code_icao"`
	CodeIata       string  `json:"code_iata"`
	CodeLid        *string `json:"code_lid"`
	Timezone       string  `json:"timezone"`
	Name           string  `json:"name"`
	City           string  `json:"city"`
	AirportInfoURL string  `json:"airport_info_url"`
}
