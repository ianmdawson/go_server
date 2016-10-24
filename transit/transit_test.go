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

	"github.com/ianmdawson/go_server/config"
)

// It creates an http client and sends a request, receives a response
func TestHTTPClient(t *testing.T) {
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
	// For some reason this response has a newline at the end :(
	trimmedResponse := strings.TrimSuffix(string(response), "\n")
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
		t.Fatal(string(response))
	}
}

func setUp() {
	config.LoadEnv("")
}

func TestGetStopsErrorsIfServerErrors(t *testing.T) {
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

	response, err := getStops(fakeURL.String())
	expectedError := "Request failed, status code 401: A valid API token is required to use the AC Transit API."
	if err == nil && err.Error() != expectedError {
		t.Fatalf("Expected error, but got: %s", err)
	}
	if response != nil {
		t.Fatalf("Expected response to be nil, but got: %s", string(response))
	}
}

func TestAppendTransitToken(t *testing.T) {
	testURL := "https://www.example.com/test"
	testToken := "1234"

	u, err := appendAuthToURL(testURL, &testToken)
	if err != nil {
		t.Fatalf("Expected no error, but got: %s", err)
	}
	expectedURL, _ := url.Parse("https://www.example.com/test?token=1234")
	if u == expectedURL {
		t.Fatalf("Expected url to be %s, but got: %s", expectedURL, u)
	}
}
