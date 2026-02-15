package services

import "math"

const earthRadiusMeters = 6_371_000

// GeoService provides geographic calculations (distance, scoring).
// It is currently stateless; the type is retained for future algorithm variants and mocking in tests.
type GeoService struct {
}

// NewGeoService returns a GeoService instance for geographic operations.
// The constructor is retained for future algorithm variants and mocking in tests.
func NewGeoService() *GeoService {
	return &GeoService{}
}

// CalculateDistance returns the distance in meters between two coordinates using the Haversine formula.
func (s *GeoService) CalculateDistance(latitude, longitude float64, targetLatitude, targetLongitude float64) float64 {
	return haversineDistance(latitude, longitude, targetLatitude, targetLongitude)
}

// GeoGuessr scoring constants.
const (
	maxScore       = 5000.0
	sigmaKm        = 22.5 // GeoGuessr sigma: at this distance you get ~60% of max score
	metersPerKm    = 1000.0
)

// CalculateScoreFromDistance returns the GeoGuessr-style score (0-5000) based on distance in meters.
// Formula: score = 5000 * exp(-0.5 * (distance_km / sigma)^2)
func (s *GeoService) CalculateScoreFromDistance(distanceMeters float64) int {
	distanceKm := distanceMeters / metersPerKm
	score := maxScore * math.Exp(-0.5*math.Pow(distanceKm/sigmaKm, 2))
	return int(math.Round(score))
}

func haversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = earthRadiusMeters

	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	dlat := lat2Rad - lat1Rad
	dlon := lon2Rad - lon1Rad

	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dlon/2)*math.Sin(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}