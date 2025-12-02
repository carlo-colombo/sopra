package model

import "fmt"

// States represents the response from the OpenSky Network API
type States struct {
	Time   int             `json:"time"`
	States [][]interface{} `json:"states"`
}

// Flight represents a flight with named fields.
type Flight struct {
	Icao24         string  `json:"icao24"`
	Callsign       string  `json:"callsign"`
	OriginCountry  string  `json:"origin_country"`
	TimePosition   int     `json:"time_position"`
	LastContact    int     `json:"last_contact"`
	Longitude      float64 `json:"longitude"`
	Latitude       float64 `json:"latitude"`
	BaroAltitude   float64 `json:"baro_altitude"`
	OnGround       bool    `json:"on_ground"`
	Velocity       float64 `json:"velocity"`
	TrueTrack      float64 `json:"true_track"`
	VerticalRate   float64 `json:"vertical_rate"`
	Sensors        []int   `json:"sensors"`
	GeoAltitude    float64 `json:"geo_altitude"`
	Squawk         string  `json:"squawk"`
	Spi            bool    `json:"spi"`
	PositionSource int     `json:"position_source"`
}

// ToFlights converts the States object to a slice of Flight objects.
func (s *States) ToFlights() []Flight {
	var flights []Flight
	for _, state := range s.States {
		flight := Flight{}
		if len(state) > 0 {
			flight.Icao24, _ = state[0].(string)
		}
		if len(state) > 1 {
			flight.Callsign, _ = state[1].(string)
		}
		if len(state) > 2 {
			flight.OriginCountry, _ = state[2].(string)
		}
		if len(state) > 3 {
			if val, ok := state[3].(float64); ok {
				flight.TimePosition = int(val)
			}
		}
		if len(state) > 4 {
			if val, ok := state[4].(float64); ok {
				flight.LastContact = int(val)
			}
		}
		if len(state) > 5 {
			flight.Longitude, _ = state[5].(float64)
		}
		if len(state) > 6 {
			flight.Latitude, _ = state[6].(float64)
		}
		if len(state) > 7 {
			flight.BaroAltitude, _ = state[7].(float64)
		}
		if len(state) > 8 {
			flight.OnGround, _ = state[8].(bool)
		}
		if len(state) > 9 {
			flight.Velocity, _ = state[9].(float64)
		}
		if len(state) > 10 {
			flight.TrueTrack, _ = state[10].(float64)
		}
		if len(state) > 11 {
			flight.VerticalRate, _ = state[11].(float64)
		}
		if len(state) > 12 {
			// Sensors are not handled in this example
		}
		if len(state) > 13 {
			flight.GeoAltitude, _ = state[13].(float64)
		}
		if len(state) > 14 {
			flight.Squawk, _ = state[14].(string)
		}
		if len(state) > 15 {
			flight.Spi, _ = state[15].(bool)
		}
		if len(state) > 16 {
			if val, ok := state[16].(float64); ok {
				flight.PositionSource = int(val)
			}
		}
		flights = append(flights, flight)
	}
	return flights
}

func (s *States) String() string {
	var result string
	for _, flight := range s.ToFlights() {
		result += fmt.Sprintf("Flight: %s\n", flight.Callsign)
	}
	return result
}
