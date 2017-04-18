package transit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
)

type stop struct {
	StopID        json.Number `json:"StopId"`
	Name          string      `json:"Name"`
	Latitude      json.Number `json:"Latitude"`
	Longitude     json.Number `json:"Longitude"`
	ScheduledTime string      `json:"ScheduledTime"`
}

type Prediction struct {
	StopID                  json.Number `json:"StopId"`
	TripID                  json.Number `json:"TripId"`
	VehicleID               json.Number `json:"VehicleId"`
	RouteName               string      `json:"RouteName"`
	PredictedDelayInSeconds json.Number `json:"PredictedDelayInSeconds"`
	PredictedDeparture      string      `json:"PredictedDeparture"`
	PredictionDateTime      string      `json:"PredictionDateTime"`
}

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

// GetAllStops retrieves all available stops
func GetAllStops(URL string) (*[]stop, error) {
	if URL == "" {
		URL = "https://api.actransit.org/transit/stops"
	}

	stopsURL, err := appendAuthToURL(URL, nil)
	if err != nil {
		return nil, err
	}

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

// GetStopPredictions retrieves predictions for a stop by ID
func GetStopPredictions(stopID string, URL string) (*[]Prediction, error) {
	regex := regexp.MustCompile("^[0-9]+")
	match := regex.FindAllString(stopID, 1)
	if match == nil {
		return nil, fmt.Errorf("Invalid stop ID: %s", stopID)
	}

	if URL == "" {
		URL = fmt.Sprintf("https://api.actransit.org/transit/stops/%s/predictions", stopID)
	}

	stopsURL, err := appendAuthToURL(URL, nil)
	if err != nil {
		return nil, err
	}

	responseBody, err := httpRequest(*stopsURL)
	if err != nil {
		return nil, err
	}

	var predictions []Prediction
	err = json.Unmarshal(*responseBody, &predictions)
	if err != nil {
		return nil, err
	}

	return &predictions, nil
}

// UsefulStops a list of StopIDs
var UsefulStops = []uint16{
	58123,
	52246,
}
