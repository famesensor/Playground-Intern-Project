package geolocation

import (
	"math"

	"github.com/mmcloughlin/geohash"
)

// GetGeoHashRange -> this is fn hash lat lng  for query
func GetGeoHashRange(latitude, longitude, distance float64) (string, string) {
	const lat = 0.0144927536231884 // degrees latitude per mile
	const lon = 0.0181818181818182 // degrees longitude per mile

	lowerLat := latitude - lat*distance
	lowerLon := longitude - lon*distance

	upperLat := latitude + lat*distance
	upperLon := longitude + lon*distance

	lower := geohash.Encode(lowerLat, lowerLon)
	upper := geohash.Encode(upperLat, upperLon)

	return lower, upper
}

func DistanceBetween(latRest, lngRest, latUser, lngUser float64, unit string) (distance float64) {
	const PI float64 = 3.141592653589793

	radlat1 := float64(PI * latRest / 180)
	radlat2 := float64(PI * latUser / 180)

	theta := float64(lngRest - lngUser)
	radtheta := float64(PI * theta / 180)

	distance = math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)

	if distance > 1 {
		distance = 1
	}

	distance = math.Acos(distance)
	distance = distance * 180 / PI
	distance = distance * 60 * 1.1515

	if unit == "K" {
		distance = distance * 1.609344
	} else if unit == "N" {
		distance = distance * 0.8684
	}

	return distance
}
