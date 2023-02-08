package geocoding

import (
	"encoding/json"
	"fmt"

	"github.com/Leyka/picor/cache"
	"github.com/Leyka/picor/geocoding/location"
	"github.com/Leyka/picor/geocoding/providers"
)

type GeocodingAPI interface {
	ReverseGeocoding(lat, long float64) (*location.Location, error)
	Ping() error
}

type GeocodingSettings struct {
	TomTomApiKey string
}

var api GeocodingAPI

func Setup(settings GeocodingSettings) {
	if settings.TomTomApiKey != "" {
		tomtom := providers.NewTomTom(settings.TomTomApiKey)
		if err := tomtom.Ping(); err == nil {
			fmt.Println("~ Using TomTom as geocoding API")
			api = tomtom
			return
		}
	}

	nominatim := providers.NewNominatim()
	if err := nominatim.Ping(); err == nil {
		fmt.Println("~ Using Nominatim as geocoding API")
		api = nominatim
		return
	}

	fmt.Println("~ No geocoding API available, fallback to mock")
	api = providers.NewMock()
}

func ReverseGeocoding(lat, long float64) (*location.Location, error) {
	var location *location.Location

	// Check if location is cached
	key := serializeLatLong(lat, long)
	if err := cache.Get(key, &location); err == nil {
		return location, nil
	}

	// Get location from API
	if api == nil {
		panic("Geocoding API is not initialized")
	}

	location, err := api.ReverseGeocoding(lat, long)
	if err != nil {
		return nil, err
	}

	// Cache location
	locationJSON, err := json.Marshal(location)
	if err != nil {
		return nil, err
	}
	cache.Set(key, locationJSON)

	return location, nil
}

func serializeLatLong(lat, long float64) string {
	// Keep 4 first decimals to avoid too many cache entries
	lat4 := fmt.Sprintf("%.4f", lat)
	long4 := fmt.Sprintf("%.4f", long)
	return fmt.Sprintf("%s,%s", lat4, long4)
}
