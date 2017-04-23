package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ianmdawson/go_server/config"
	"github.com/ianmdawson/go_server/handlers"
)

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

	http.Handle("/", http.FileServer(http.Dir("./public")))

	http.HandleFunc("/transit/all", handlers.AllTransitStopsHandler)
	http.HandleFunc("/transit/stop/", handlers.TransitStopHandler)

	fmt.Printf("Getting ready to serve on port: %s", port)
	http.ListenAndServe(":"+port, Log(http.DefaultServeMux))
}
