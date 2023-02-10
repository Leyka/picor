package geocoding

type mock struct{}

func newMock() *mock {
	return &mock{}
}

func (m *mock) ReverseGeocoding(lat float64, long float64) (*Location, error) {
	return UNKNOWN_LOCATION, nil
}

func (m *mock) Ping() error {
	return nil
}
