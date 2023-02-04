package geocoder

import (
	"encoding/json"
	"fmt"

	"github.com/Leyka/picor/cache"
)

var UNKNOWN_LOCATION = &Location{
	City:    "Unknown City",
	Country: "Unknown Country",
}

type GeoCode struct {
	Latitude  float64
	Longitude float64
}

type Location struct {
	City    string
	Country string
}

func NewGeoCode(lat, lng float64) *GeoCode {
	return &GeoCode{
		Latitude:  lat,
		Longitude: lng,
	}
}

func (g *GeoCode) format() string {
	// Keep 4 first decimals to avoid too many cache entries
	lat := fmt.Sprintf("%.4f", g.Latitude)
	lng := fmt.Sprintf("%.4f", g.Longitude)

	return fmt.Sprintf("%s,%s", lat, lng)
}

func (g *GeoCode) GetLocation() (*Location, error) {
	key := g.format()

	// Check if location is cached
	if res, err := cache.Instance.Get(key); res != nil && err == nil {
		bytes, ok := res.([]byte)
		if !ok {
			str := res.(string)
			bytes = []byte(str)
		}

		var location Location
		err := json.Unmarshal(bytes, &location)
		if err != nil {
			return nil, err
		}
		return &location, nil
	}

	// Get location from API
	location, err := api.ReverseGeocode(g)
	if err != nil {
		return nil, err
	}

	// Cache location
	locationJSON, err := json.Marshal(location)
	if err != nil {
		return nil, err
	}
	if err := cache.Instance.Set(key, locationJSON); err != nil {
		return nil, err
	}

	return location, nil
}
