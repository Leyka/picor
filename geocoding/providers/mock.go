package providers

import (
	"github.com/Leyka/picor/geocoding/location"
)

type Mock struct{}

func NewMock() *Mock {
	return &Mock{}
}

func (m *Mock) ReverseGeocoding(lat float64, long float64) (*location.Location, error) {
	return location.UNKNOWN_LOCATION, nil
}

func (m *Mock) Ping() error {
	return nil
}
