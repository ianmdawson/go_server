package actransit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"time"
)

type Stop struct {
	StopID        json.Number `json:"StopId"`
	Name          string      `json:"Name"`
	Latitude      json.Number `json:"Latitude,omitempty"`
	Longitude     json.Number `json:"Longitude,omitempty"`
	ScheduledTime string      `json:"ScheduledTime,omitempty"`
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

// TimeUntilPredictedDeparture returns a time.Duration from now until the
// predicted departure time of a stop
func (prediction *Prediction) TimeUntilPredictedDeparture() (*time.Duration, error) {
	arrivalTime, err := getTimeFromACTransit(prediction.PredictedDeparture)
	if err != nil {
		return nil, err
	}

	currentTime := time.Now()
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		return nil, err
	}
	currentTimeInPST := currentTime.In(loc)
	difference := arrivalTime.Sub(currentTimeInPST)

	// It's not necessary to have precision beyond seconds here
	differenceWithTruncatedSeconds := truncateSeconds(difference)
	return &differenceWithTruncatedSeconds, nil
}

func truncateSeconds(duration time.Duration) time.Duration {
	return duration - (duration % time.Second)
}

// IsDelayed method returns true if the line for this prediction is delayed
func (prediction *Prediction) IsDelayed() bool {
	predictedDelay, _ := prediction.PredictedDelayInSeconds.Int64()
	if predictedDelay != 0 {
		return true
	}
	return false
}

func getTimeFromACTransit(acTransitTime string) (*time.Time, error) {
	PST, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		return nil, err
	}
	t, err := time.ParseInLocation("2006-01-02T15:04:05", acTransitTime, PST)
	if err != nil {
		return nil, err
	}
	return &t, err
}

func formatTimeForACTransit(t time.Time) string {
	return t.Format("2006-01-02T15:04:05")
}

func appendAuthToURL(URLPrefix string, testToken *string) (*url.URL, error) {
	var actransitToken string
	if testToken == nil {
		actransitToken = os.Getenv("ACTRANSIT_TOKEN")
	} else {
		actransitToken = *testToken
	}

	var tokenSuffix = "?token=" + actransitToken
	_url, err := url.ParseRequestURI(URLPrefix + tokenSuffix)
	if err != nil {
		return nil, err
	}

	return _url, nil
}

// super basic client
func httpRequest(u url.URL) (*[]byte, error) {
	res, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		// TODO: handle 404s
		// res.StatusCode == 404 {
		// }
		return nil, fmt.Errorf("Request failed, status code %d", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &body, nil
}

// GetAllStops retrieves all available stops
func GetAllStops(URL string) (*[]Stop, error) {
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

	var stops []Stop
	err = json.Unmarshal(*responseBody, &stops)
	if err != nil {
		return nil, err
	}

	return &stops, nil
}

// GetPredictionsForStop retrieves predictions for a stop by ID
func GetPredictionsForStop(stopID string, URL string) (*[]Prediction, error) {
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

	// TODO: handle 404 not founds
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
