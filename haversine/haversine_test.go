package haversine

import (
	"math"
	"testing"
)

const float64EqualityThreshold = 1e-6

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
}

func TestGetBoundingBox(t *testing.T) {
	lat, lon := 40.7128, -74.0060 // New York City
	radiusKm := 10.0

	bbox := GetBoundingBox(lat, lon, radiusKm)

	expectedMinLat := 40.622868
	expectedMaxLat := 40.802732
	expectedMinLon := -74.124646
	expectedMaxLon := -73.887354

	if !almostEqual(bbox.MinLat, expectedMinLat) {
		t.Errorf("Expected MinLat to be around %f, but got %f", expectedMinLat, bbox.MinLat)
	}
	if !almostEqual(bbox.MaxLat, expectedMaxLat) {
		t.Errorf("Expected MaxLat to be around %f, but got %f", expectedMaxLat, bbox.MaxLat)
	}
	if !almostEqual(bbox.MinLon, expectedMinLon) {
		t.Errorf("Expected MinLon to be around %f, but got %f", expectedMinLon, bbox.MinLon)
	}
	if !almostEqual(bbox.MaxLon, expectedMaxLon) {
		t.Errorf("Expected MaxLon to be around %f, but got %f", expectedMaxLon, bbox.MaxLon)
	}
}

func TestRadToDeg(t *testing.T) {
	rad := math.Pi
	expectedDeg := 180.0
	if deg := radToDeg(rad); !almostEqual(deg, expectedDeg) {
		t.Errorf("Expected %f, but got %f", expectedDeg, deg)
	}
}

func TestDegToRad(t *testing.T) {
	deg := 180.0
	expectedRad := math.Pi
	if rad := degToRad(deg); !almostEqual(rad, expectedRad) {
		t.Errorf("Expected %f, but got %f", expectedRad, rad)
	}
}
