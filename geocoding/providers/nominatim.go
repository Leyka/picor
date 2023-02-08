package providers

import (
	"fmt"

	"github.com/Leyka/picor/geocoding/location"
	"github.com/Leyka/picor/utils"
)

const NOMINATIM_BASE_URL = "https://nominatim.openstreetmap.org"

// Nominatim uses OpenStreetMap data to find locations.
// Free ; no api key required ; rate = 1 req/sec.
// Doc: https://nominatim.org/release-docs/develop/api/Overview
type Nominatim struct{}

func NewNominatim() *Nominatim {
	return &Nominatim{}
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

func (n *Nominatim) ReverseGeocoding(lat, long float64) (*location.Location, error) {
	url := fmt.Sprintf("%s/reverse?lat=%f&lon=%f&format=jsonv2", NOMINATIM_BASE_URL, lat, long)

	var res nominatimReverseGeocodingResponse
	err := utils.HttpGet(url, &res)
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

	return &location.Location{
		Country:     addr.Country,
		CountryCode: addr.CountryCode,
		State:       addr.State,
		City:        city,
	}, nil
}

func (n *Nominatim) Ping() error {
	_, err := n.ReverseGeocoding(40.730610, -73.935242)
	return err
}
