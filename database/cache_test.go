package database

import (
	"fmt"
	"os"
	"testing"

	"github.com/carlo-colombo/sopra/model"
	"github.com/stretchr/testify/assert"
)

func TestAirportStats(t *testing.T) {
	dbName := fmt.Sprintf("%s.db", t.Name())
	os.Remove(dbName)
	db, err := NewDB(dbName)
	if err != nil {
		t.Fatalf("failed to create test db: %v", err)
	}
	t.Cleanup(func() {
		db.Close()
		os.Remove(dbName)
	})

	err = db.ClearFlightLog()
	assert.NoError(t, err)

	flights := []*model.FlightInfo{
		{
			Ident: "F1",
			Origin: model.AirportDetail{CodeIata: "ZRH", City: "Zurich"},
			Destination: model.AirportDetail{CodeIata: "JFK", City: "New York"},
		},
		{
			Ident: "F2",
			Origin: model.AirportDetail{CodeIata: "ZRH", City: "Zurich"},
			Destination: model.AirportDetail{CodeIata: "LAX", City: "Los Angeles"},
		},
		{
			Ident: "F3",
			Origin: model.AirportDetail{CodeIata: "LHR", City: "London"},
			Destination: model.AirportDetail{CodeIata: "JFK", City: "New York"},
		},
	}

	for _, f := range flights {
		err = db.LogFlight(f.Ident, f)
		assert.NoError(t, err)
	}

    // Log F1 again to increment identification_count
    err = db.LogFlight("F1", flights[0])
    assert.NoError(t, err)

	topDest, err := db.GetTopDestinations()
	assert.NoError(t, err)
	assert.Len(t, topDest, 2)
    // JFK should have count 3 (2 from F1, 1 from F3)
    // LAX should have count 1 (from F2)
    assert.Equal(t, "JFK", topDest[0].Iata)
    assert.Equal(t, 3, topDest[0].Count)
    assert.Equal(t, "LAX", topDest[1].Iata)
    assert.Equal(t, 1, topDest[1].Count)

	topSrc, err := db.GetTopSources()
	assert.NoError(t, err)
	assert.Len(t, topSrc, 2)
    // ZRH should have count 3 (2 from F1, 1 from F2)
    // LHR should have count 1 (from F3)
    assert.Equal(t, "ZRH", topSrc[0].Iata)
    assert.Equal(t, 3, topSrc[0].Count)
    assert.Equal(t, "LHR", topSrc[1].Iata)
    assert.Equal(t, 1, topSrc[1].Count)
}
