package haversine

import "math"

const (
	earthRadiusKm = 6371 // kilometers
)

// BoundingBox represents a bounding box with min and max latitude and longitude.
type BoundingBox struct {
	MinLat float64
	MinLon float64
	MaxLat float64
	MaxLon float64
}

// GetBoundingBox calculates a bounding box for a given location and radius in kilometers.
func GetBoundingBox(lat, lon, radiusKm float64) BoundingBox {
	rad := radiusKm / earthRadiusKm

	minLat := lat - radToDeg(rad)
	maxLat := lat + radToDeg(rad)

	deltaLon := math.Asin(math.Sin(rad) / math.Cos(degToRad(lat)))

	minLon := lon - radToDeg(deltaLon)
	maxLon := lon + radToDeg(deltaLon)

	return BoundingBox{
		MinLat: minLat,
		MinLon: minLon,
		MaxLat: maxLat,
		MaxLon: maxLon,
	}
}

func radToDeg(rad float64) float64 {
	return rad * 180 / math.Pi
}

func degToRad(deg float64) float64 {
	return deg * math.Pi / 180
}

// Distance calculates the shortest path between two coordinates on the surface of a sphere.
func Distance(lat1, lon1, lat2, lon2 float64) float64 {
	dLat := degToRad(lat2 - lat1)
	dLon := degToRad(lon2 - lon1)

	lat1 = degToRad(lat1)
	lat2 = degToRad(lat2)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1)*math.Cos(lat2)
	c := 2 * math.Asin(math.Sqrt(a))

	return earthRadiusKm * c
}
