package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ianmdawson/go_server/config"
	"github.com/ianmdawson/go_server/transit"
)

func viewHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", "Title", "Body")
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int, message string) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprintf(w, "Not Found 👻 -- %s", message)
	}
}

func allTransitStopsHandler(w http.ResponseWriter, r *http.Request) {
	stops, err := transit.GetAllStops("")
	if err != nil {
		http.Error(w, "Something went wrong while trying to retrieve AC Transit stops: "+err.Error(), http.StatusBadGateway)
		return
	}

	if len(*stops) <= 0 {
		errorHandler(w, r, http.StatusNotFound, "No stops found")
	}

	fmt.Fprintf(w, "%v", stops)
}

func transitStopHandler(w http.ResponseWriter, r *http.Request) {
	stopID := r.URL.Path[len("/transit/stop/"):]
	// ensure stopID is a number
	stops, err := transit.GetStopPredictions(stopID, "")
	if err != nil {
		http.Error(w, "Something went wrong while trying to retrieve AC Transit stops: "+err.Error(), http.StatusBadGateway)
		return
	}

	if len(*stops) <= 0 {
		errorHandler(w, r, http.StatusNotFound, "No stops found")
	}

	fmt.Fprintf(w, "%v", stops)
}

// Log logging middleware
func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	if !config.IsProductionEnvironment() {
		config.LoadEnv("")
	}

	http.HandleFunc("/api/view", viewHandler)
	http.Handle("/", http.FileServer(http.Dir("./public")))

	http.HandleFunc("/transit/all", allTransitStopsHandler)
	http.HandleFunc("/transit/stop/", transitStopHandler)

	fmt.Printf("Getting ready to serve on port: %s", port)
	http.ListenAndServe(":"+port, Log(http.DefaultServeMux))
}
