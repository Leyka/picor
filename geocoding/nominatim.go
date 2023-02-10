package geocoding

import (
	"fmt"
)

const NOMINATIM_BASE_URL = "https://nominatim.openstreetmap.org"

// nominatim uses OpenStreetMap data to find locations.
// Free ; no api key required ; rate = 1 req/sec.
// Doc: https://nominatim.org/release-docs/develop/api/Overview
type nominatim struct{}

func newNominatim() *nominatim {
	return &nominatim{}
}

type nominatimReverseGeocodingResponse struct {
	Address struct {
		Country     string `json:"country"`
		CountryCode string `json:"country_code"`
		State       string `json:"state"`
		City        string `json:"city"`
		County      string `json:"county"`
		Suburb      string `json:"suburb"`
	} `json:"address"`
}

func (n *nominatim) ReverseGeocoding(lat, long float64) (*Location, error) {
	url := fmt.Sprintf("%s/reverse?lat=%f&lon=%f&format=jsonv2", NOMINATIM_BASE_URL, lat, long)

	var res nominatimReverseGeocodingResponse
	err := HttpGet(url, &res)
	if err != nil {
		return nil, err
	}

	addr := res.Address

	var city string
	if addr.City != "" {
		city = addr.City
	} else if addr.Suburb != "" {
		city = addr.Suburb
	} else if addr.County != "" {
		city = addr.County
	}

	return &Location{
		Country:     addr.Country,
		CountryCode: addr.CountryCode,
		State:       addr.State,
		City:        city,
	}, nil
}

func (n *nominatim) Ping() error {
	_, err := n.ReverseGeocoding(40.730610, -73.935242)
	return err
}
