package geocoding

type Location struct {
	Country     string
	CountryCode string
	State       string
	City        string
}

type GeocodingAPI interface {
	Ping() error
	ReverseGeocoding(lat, lon float64) (*Location, error)
}

type GeocodingSettings struct {
	TomTomApiKey string
}

var API GeocodingAPI

func Setup(settings GeocodingSettings) {
	// TODO
}
