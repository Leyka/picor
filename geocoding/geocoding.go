package geocoding

import (
	"encoding/json"
	"fmt"
)

type Location struct {
	Country     string
	CountryCode string
	State       string
	City        string
}

var UNKNOWN_LOCATION = &Location{
	Country:     "~ unknown country",
	CountryCode: "~ unknown country code",
	State:       "~ unknown state",
	City:        "~ unknown city",
}

type GeocodingAPI interface {
	ReverseGeocoding(lat, long float64) (*Location, error)
	Ping() error
}

type GeocodingSettings struct {
	TomTomApiKey string
}

var api GeocodingAPI

func Setup(settings GeocodingSettings) {
	setupCache()

	if settings.TomTomApiKey != "" {
		tomtom := newTomTom(settings.TomTomApiKey)
		if err := tomtom.Ping(); err == nil {
			fmt.Println("~ Using TomTom as geocoding API")
			api = tomtom
			return
		}
	}

	nominatim := newNominatim()
	if err := nominatim.Ping(); err == nil {
		fmt.Println("~ Using Nominatim as geocoding API")
		api = nominatim
		return
	}

	fmt.Println("~ No geocoding API available, fallback to mock")
	api = newMock()
}

func Close() {
	closeCache()
}

func ReverseGeocoding(lat, long float64) (*Location, error) {
	var location *Location

	// Check if location is cached
	key := serializeLatLong(lat, long)
	if err := getCache(key, &location); err == nil {
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
	setCache(key, locationJSON)

	return location, nil
}

func serializeLatLong(lat, long float64) string {
	// Keep 4 first decimals to avoid too many cache entries
	lat4 := fmt.Sprintf("%.4f", lat)
	long4 := fmt.Sprintf("%.4f", long)
	return fmt.Sprintf("%s,%s", lat4, long4)
}
