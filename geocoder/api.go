package geocoder

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var api geocodeAPI

type geocodeAPI interface {
	ReverseGeocode(*GeoCode) (*Location, error)
	Ping() error
}

// Check if TomTom API key is set and working, otherwise use OpenStreetMap API.
// If neither are working, use the mock API
func SetupGeocoderAPI() {
	tomtomAPIKey := os.Getenv("TOMTOM_API_KEY")
	if tomtomAPIKey != "" {
		tomtom := newTomTomGeocoder(tomtomAPIKey)
		if err := tomtom.Ping(); err == nil {
			fmt.Println("Using TomTom API as geocoding service")
			api = tomtom
			return
		}
	}

	osm := newOpenStreetMapGeocoder()
	if err := osm.Ping(); err == nil {
		fmt.Println("Using OpenStreetMap API as geocoding service")
		api = osm
		return
	}

	fmt.Println("No API as geocoding service")
	api = newMockGeocoder()
}

// ~ TomTom Geocoding ~
// Free. Rate limit 2,500 requests per day. API key required.
type tomTomGeocoder struct {
	apiKey string
}

type tomtomResponse struct {
	Addresses []struct {
		Address struct {
			City    string `json:"municipality"`
			Country string `json:"country"`
		} `json:"address"`
	} `json:"addresses"`
}

func newTomTomGeocoder(apiKey string) *tomTomGeocoder {
	return &tomTomGeocoder{apiKey: apiKey}
}

func (t *tomTomGeocoder) ReverseGeocode(geocode *GeoCode) (*Location, error) {
	url := fmt.Sprintf("https://api.tomtom.com/search/2/reverseGeocode/%f,%f.json?key=%s&language=en-US", geocode.Latitude, geocode.Longitude, t.apiKey)
	data, err := get(url)
	if err != nil {
		return nil, err
	}

	var response tomtomResponse
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(dataBytes, &response); err != nil {
		return nil, err
	}

	city := response.Addresses[0].Address.City
	country := response.Addresses[0].Address.Country
	return &Location{City: city, Country: country}, nil
}

func (t *tomTomGeocoder) Ping() error {
	_, err := t.ReverseGeocode(&GeoCode{Latitude: 40.730610, Longitude: -73.935242})
	return err
}

// ~ OpenStreetMap Geocoding ~
// Free. Rate limit 1 request per second. No API key required.
type OpenStreetMapGeocoder struct{}

func newOpenStreetMapGeocoder() *OpenStreetMapGeocoder {
	return &OpenStreetMapGeocoder{}
}

func (o *OpenStreetMapGeocoder) ReverseGeocode(geocode *GeoCode) (*Location, error) {
	url := fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?lat=%f&lon=%f&format=json", geocode.Latitude, geocode.Longitude)
	data, err := get(url)
	if err != nil {
		return nil, err
	}

	address := (*data)["address"].(map[string]interface{})

	var city string
	if maybeCity, ok := address["city"]; ok {
		city = maybeCity.(string)
	} else if state, ok := address["state"]; ok {
		city = state.(string)
	} else if suburb, ok := address["suburb"]; ok {
		city = suburb.(string)
	}
	country := address["country"].(string)

	return &Location{City: city, Country: country}, nil
}

func (o *OpenStreetMapGeocoder) Ping() error {
	_, err := o.ReverseGeocode(&GeoCode{Latitude: 40.730610, Longitude: -73.935242})
	return err
}

// ~ Mock Geocoding ~
type MockGeocoder struct{}

func newMockGeocoder() *MockGeocoder {
	return &MockGeocoder{}
}

func (m *MockGeocoder) ReverseGeocode(geocode *GeoCode) (*Location, error) {
	return UNKNOWN_LOCATION, nil
}

func (m *MockGeocoder) Ping() error {
	return nil
}

// Calls http get and returns JSON data
func get(url string) (*map[string]interface{}, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return &data, nil
}
