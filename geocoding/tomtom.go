package geocoding

import (
	"fmt"
)

const TOMTOM_BASE_URL = "https://api.tomtom.com"

// Free with api key; rate = 2500 req/day.
// doc: https://developer.tomtom.com
type tomtom struct {
	apiKey string
}

func newTomTom(apiKey string) *tomtom {
	return &tomtom{apiKey}
}

type tomTomReverseGeocodingResponse struct {
	Addresses []struct {
		Address struct {
			Country     string `json:"country"`
			CountryCode string `json:"countryCode"`
			State       string `json:"countrySubdivisionName"`
			City        string `json:"municipality"`
		} `json:"address"`
	} `json:"addresses"`
}

// Doc: https://developer.tomtom.com/reverse-geocoding-api/api-explorer
func (t *tomtom) ReverseGeocoding(lat, long float64) (*Location, error) {
	url := fmt.Sprintf("%s/search/2/reverseGeocode/%f,%f.json?key=%s&language=en-US", TOMTOM_BASE_URL, lat, long, t.apiKey)

	var res tomTomReverseGeocodingResponse
	err := HttpGet(url, &res)
	if err != nil {
		return nil, err
	}

	addr := res.Addresses[0].Address
	return &Location{
		Country:     addr.Country,
		CountryCode: addr.CountryCode,
		State:       addr.State,
		City:        addr.City,
	}, nil
}

func (t *tomtom) Ping() error {
	_, err := t.ReverseGeocoding(40.730610, -73.935242)
	return err
}
