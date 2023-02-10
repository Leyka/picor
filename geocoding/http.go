package geocoding

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Returns body as bytes of GET request to url
func HttpGet[T any](url string, res *T) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bodyBytes, res); err != nil {
		return err
	}

	return nil
}
