package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToFlights(t *testing.T) {
	states := &States{
		Time: 1,
		States: [][]interface{}{
			{"icao24", "callsign", "origin_country", float64(1), float64(2), 3.0, 4.0, 5.0, true, 6.0, 7.0, 8.0, nil, 9.0, "squawk", false, float64(0)},
			{"icao24_2", nil, "origin_country_2", float64(10), float64(20), 30.0, 40.0, 50.0, false, 60.0, 70.0, 80.0, nil, 90.0, nil, true, float64(1)},
		},
	}

	flights := states.ToFlights()

	assert.Len(t, flights, 2)

	// Test first flight
	assert.Equal(t, "icao24", flights[0].Icao24)
	assert.Equal(t, "callsign", flights[0].Callsign)
	assert.Equal(t, "origin_country", flights[0].OriginCountry)
	assert.Equal(t, 1, flights[0].TimePosition)
	assert.Equal(t, 2, flights[0].LastContact)
	assert.Equal(t, 3.0, flights[0].Longitude)
	assert.Equal(t, 4.0, flights[0].Latitude)
	assert.Equal(t, 5.0, flights[0].BaroAltitude)
	assert.Equal(t, true, flights[0].OnGround)
	assert.Equal(t, 6.0, flights[0].Velocity)
	assert.Equal(t, 7.0, flights[0].TrueTrack)
	assert.Equal(t, 8.0, flights[0].VerticalRate)
	assert.Equal(t, 9.0, flights[0].GeoAltitude)
	assert.Equal(t, "squawk", flights[0].Squawk)
	assert.Equal(t, false, flights[0].Spi)
	assert.Equal(t, 0, flights[0].PositionSource)

	// Test second flight with nil values
	assert.Equal(t, "icao24_2", flights[1].Icao24)
	assert.Equal(t, "", flights[1].Callsign)
	assert.Equal(t, "origin_country_2", flights[1].OriginCountry)
	assert.Equal(t, 10, flights[1].TimePosition)
	assert.Equal(t, 20, flights[1].LastContact)
	assert.Equal(t, 30.0, flights[1].Longitude)
	assert.Equal(t, 40.0, flights[1].Latitude)
	assert.Equal(t, 50.0, flights[1].BaroAltitude)
	assert.Equal(t, false, flights[1].OnGround)
	assert.Equal(t, 60.0, flights[1].Velocity)
	assert.Equal(t, 70.0, flights[1].TrueTrack)
	assert.Equal(t, 80.0, flights[1].VerticalRate)
	assert.Equal(t, 90.0, flights[1].GeoAltitude)
	assert.Equal(t, "", flights[1].Squawk)
	assert.Equal(t, true, flights[1].Spi)
	assert.Equal(t, 1, flights[1].PositionSource)
}
