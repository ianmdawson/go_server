package transit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

func appendAuthToURL(URLPrefix string, testToken *string) (*url.URL, error) {
	var actransitToken string
	if testToken == nil {
		actransitToken = os.Getenv("ACTRANSIT_TOKEN")
	} else {
		actransitToken = *testToken
	}

	var tokenSuffix = "?token=" + actransitToken
	u, err := url.Parse(URLPrefix + tokenSuffix)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// super basic client
func httpRequest(u url.URL) (*[]byte, error) {
	res, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("Request failed, status code %d", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &body, nil
}

type stop struct {
	StopID        json.Number `json:"StopId"`
	Name          string      `json:"Name"`
	Latitude      json.Number `json:"Latitude"`
	Longitude     json.Number `json:"Longitude"`
	ScheduledTime string      `json:"ScheduledTime"`
}

// GetAllStops retrieves all available stops
func GetAllStops(URL string) (*[]stop, error) {
	if URL == "" {
		URL = "https://api.actransit.org/transit/stops"
	}

	stopsURL, err := appendAuthToURL(URL, nil)
	if err != nil {
		return nil, err
	}

	fmt.Println("getting stops")
	responseBody, err := httpRequest(*stopsURL)
	if err != nil {
		return nil, err
	}

	var stops []stop
	err = json.Unmarshal(*responseBody, &stops)
	if err != nil {
		return nil, err
	}

	return &stops, nil
}

// UsefulStops a list of StopIDs
var UsefulStops = []uint16{
	1234,
	5678,
}
