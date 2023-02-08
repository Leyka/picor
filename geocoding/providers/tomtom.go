package providers

import (
	"fmt"

	"github.com/Leyka/picor/geocoding/location"
	"github.com/Leyka/picor/utils"
)

const TOMTOM_BASE_URL = "https://api.tomtom.com"

// Free with api key; rate = 2500 req/day.
// doc: https://developer.tomtom.com
type TomTom struct {
	apiKey string
}

func NewTomTom(apiKey string) *TomTom {
	return &TomTom{apiKey}
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
func (t *TomTom) ReverseGeocoding(lat, long float64) (*location.Location, error) {
	url := fmt.Sprintf("%s/search/2/reverseGeocode/%f,%f.json?key=%s&language=en-US", TOMTOM_BASE_URL, lat, long, t.apiKey)

	var res tomTomReverseGeocodingResponse
	err := utils.HttpGet(url, &res)
	if err != nil {
		return nil, err
	}

	addr := res.Addresses[0].Address
	return &location.Location{
		Country:     addr.Country,
		CountryCode: addr.CountryCode,
		State:       addr.State,
		City:        addr.City,
	}, nil
}

func (t *TomTom) Ping() error {
	_, err := t.ReverseGeocoding(40.730610, -73.935242)
	return err
}
