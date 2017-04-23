package transit

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/ianmdawson/go_server/config"
	"github.com/stretchr/testify/assert"
)

func TestAppendAuthToURL_Success(t *testing.T) {
	testURL := "https://www.example.com/test"
	testToken := "1234"

	u, err := appendAuthToURL(testURL, &testToken)
	if err != nil {
		t.Fatalf("Expected no error, but got: %s", err)
	}
	expectedURL, _ := url.Parse("https://www.example.com/test?token=1234")
	assert.Equal(t, expectedURL, u)
}

func TestAppendAuthToURL_ReturnsAnErrorForInvalidURLs(t *testing.T) {
	testURL := "badurl"
	testToken := "1234"

	u, err := appendAuthToURL(testURL, &testToken)
	assert.Error(t, err)
	assert.Nil(t, u)
}

// It creates an http client and sends a request, receives a response
func TestHTTPClient_Success(t *testing.T) {
	fakeResponse := map[string]string{"fakeData": "fakety"}
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fakeResponse)
		return
	}))
	defer s.Close()

	fakeURL, _ := url.Parse(s.URL)
	response, err := httpRequest(*fakeURL)
	if err != nil {
		t.Fatal(err)
	}

	// For some reason this response has a newline at the end
	trimmedResponse := strings.TrimSuffix(string(*response), "\n")
	if response != nil && trimmedResponse != `{"fakeData":"fakety"}` {
		t.Fatalf("response:%s!!", trimmedResponse)
	}
}

func TestHTTPClientErrorsIfServerErrors(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		return
	}))
	defer s.Close()

	fakeURL, _ := url.Parse(s.URL)
	response, err := httpRequest(*fakeURL)
	if err == nil && err != errors.New("Request failed, status code 404: %s") {
		t.Fatal(err)
	}
	if response != nil {
		t.Fatalf("Expected no response to be returned, but got: %v", response)
	}
}

func TestGetAllStopsErrorsIfServerErrors(t *testing.T) {
	setUp()
	fakeResponse := `A valid API token is required to use the AC Transit API.`
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "%s", fakeResponse)

		return
	}))
	defer s.Close()
	fakeURL, _ := url.Parse(s.URL)

	stops, err := GetAllStops(fakeURL.String())
	expectedError := "Request failed, status code 401: A valid API token is required to use the AC Transit API."
	if err == nil && err.Error() != expectedError {
		t.Fatalf("Expected error, but got: %s", err)
	}
	if stops != nil {
		t.Fatalf("Expected stops to be nil, but got: %v", stops)
	}
}

func TestGetAllStopsSuccess(t *testing.T) {
	setUp()
	fakeResponse := []map[string]string{
		map[string]string{
			"StopId":        "58123",
			"Name":          "3rd St:Santa Clara Av",
			"Latitude":      "37.7732681",
			"Longitude":     "-122.2882275",
			"ScheduledTime": "null",
		},
		map[string]string{
			"StopId":        "52246",
			"Name":          "8th St:Portola Av",
			"Latitude":      "37.7688136",
			"Longitude":     "-122.2729918",
			"ScheduledTime": "null",
		},
	}
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fakeResponse)
		return
	}))
	defer s.Close()
	fakeURL, _ := url.Parse(s.URL)

	stops, err := GetAllStops(fakeURL.String())

	if err != nil {
		t.Fatalf("Expected no error, but got: %s", err)
	}

	expectedStop := Stop{
		StopID:        "58123",
		Name:          "3rd St:Santa Clara Av",
		Latitude:      "37.7732681",
		Longitude:     "-122.2882275",
		ScheduledTime: "null",
	}

	actualStops := *stops
	assert.Equal(t, expectedStop, actualStops[0], "")
}

func TestGetStopPredictions_Success(t *testing.T) {
	setUp()
	fakeResponse := []map[string]string{
		map[string]string{
			"StopId":                  "55765",
			"TripId":                  "5340688",
			"VehicleId":               "5019",
			"RouteName":               "80",
			"PredictedDelayInSeconds": "-240",
			"PredictedDeparture":      "2017-04-17T22:30:00",
			"PredictionDateTime":      "2017-04-17T22:28:58",
		},
		map[string]string{
			"StopId":                  "55765",
			"TripId":                  "5340689",
			"VehicleId":               "5117",
			"RouteName":               "80",
			"PredictedDelayInSeconds": "-1860",
			"PredictedDeparture":      "2017-04-17T22:43:00",
			"PredictionDateTime":      "2017-04-17T22:28:48",
		},
	}
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fakeResponse)
		return
	}))
	defer s.Close()
	fakeURL, _ := url.Parse(s.URL)

	predictions, err := GetStopPredictions("55765", fakeURL.String())
	assert.NoError(t, err)

	expectedPredictedDeparture := "2017-04-17T22:43:00"
	expectedPredictionDateTime := "2017-04-17T22:28:48"
	expectedPrediction := getTestPrediction(&expectedPredictedDeparture, &expectedPredictionDateTime)
	actualPredictions := *predictions
	assert.Equal(t, expectedPrediction, actualPredictions[1], "")
}

func TestGetStopPredictions_RejectsNonNumberStopIDs(t *testing.T) {
	_, err := GetStopPredictions("zomgNotNumbers", "someurl")
	assert.EqualError(t, err, "Invalid stop ID: zomgNotNumbers")
}

func TestPredictionTimeUntilPredictedDeparture(t *testing.T) {
	testPrediction := getTestPrediction(nil, nil)
	timeUntil, err := testPrediction.TimeUntilPredictedDeparture()
	assert.NoError(t, err)
	expectedDurationMax := time.Minute * 15

	isTimeUntilAsExpected := false
	if *timeUntil <= expectedDurationMax {
		isTimeUntilAsExpected = true
	}

	assert.True(t, isTimeUntilAsExpected)
}

func TestIsDelayed(t *testing.T) {
	testPrediction := getTestPrediction(nil, nil)
	assert.True(t, testPrediction.IsDelayed())
}

// Helpers
func setUp() {
	config.LoadEnv("")
}

func getTestPrediction(fixedPredictedDeparture *string, fixedPredictionDateTime *string) Prediction {
	currentTime := time.Now().UTC()

	var predictedDeparture string
	if fixedPredictedDeparture != nil {
		predictedDeparture = *fixedPredictedDeparture
	} else {
		fifteenMinutesLater := currentTime.Add(time.Minute * 15)
		predictedDeparture = formatTimeForACTransit(fifteenMinutesLater)
	}

	var predictionDateTime string
	if fixedPredictionDateTime != nil {
		predictionDateTime = *fixedPredictionDateTime
	} else {
		predictionDateTime = formatTimeForACTransit(currentTime)
	}

	return Prediction{
		StopID:                  "55765",
		TripID:                  "5340689",
		VehicleID:               "5117",
		RouteName:               "80",
		PredictedDelayInSeconds: "-1860",
		PredictedDeparture:      predictedDeparture,
		PredictionDateTime:      predictionDateTime,
	}
}
