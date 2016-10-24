package transit

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

// basic client
func httpRequest(u url.URL) ([]byte, error) {
	res, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("Request failed, status code %d: %s", res.StatusCode, resBody)
	}

	return resBody, nil
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

func getStops(URL string) ([]byte, error) {
	if URL == "" {
		URL = "https://api.actransit.org/transit/stops"
	}

	stopsURL, err := appendAuthToURL(URL, nil)
	if err != nil {
		return nil, err
	}

	response, err := httpRequest(*stopsURL)
	if err != nil {
		return nil, err
	}

	return response, nil
}
