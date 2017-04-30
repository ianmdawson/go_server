package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/ianmdawson/goactransit"
)

// PredictionsPage is an page html template struct
type PredictionsPage struct {
	Title       string
	Predictions []actransit.Prediction
}

func compilePredictionPageFromTemplate(title string, predictions []actransit.Prediction) (*template.Template, *PredictionsPage, error) {
	page := &PredictionsPage{Title: title, Predictions: predictions}
	_template, err := template.ParseFiles("public/template/transit/predictionsTemplate.html")
	if err != nil {
		return nil, nil, err
	}

	return _template, page, nil
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int, message string) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprintf(w, "Not Found 👻 -- %s", message)
	}
}

// AllTransitStopsHandler gets all transit stops then returns them
func AllTransitStopsHandler(w http.ResponseWriter, r *http.Request) {
	stops, err := actransit.GetAllStops("")
	if err != nil {
		http.Error(w, "Something went wrong while trying to retrieve AC Transit stops: "+err.Error(), http.StatusBadGateway)
		return
	}

	if len(*stops) <= 0 {
		errorHandler(w, r, http.StatusNotFound, "No stops found")
	}

	fmt.Fprintf(w, "%v", stops)
}

// TransitStopHandler returns predictions for a specific stop
func TransitStopHandler(w http.ResponseWriter, r *http.Request) {
	stopID := r.URL.Path[len("/transit/stop/"):]
	predictions, err := actransit.GetPredictionsForStop(stopID, "")
	if err != nil {
		http.Error(w, "Something went wrong while trying to retrieve AC Transit stops: "+err.Error(), http.StatusBadGateway)
		return
	}

	if len(*predictions) <= 0 {
		errorHandler(w, r, http.StatusNotFound, "No stops found")
	}

	_template, page, err := compilePredictionPageFromTemplate("Predictions", *predictions)
	if err != nil {
		http.Error(w, "Something went wrong while trying to parse template: "+err.Error(), http.StatusBadGateway)
		return
	}

	_template.Execute(w, page)
}
