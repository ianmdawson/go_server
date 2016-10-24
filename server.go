package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ianmdawson/go_server/config"
)

func viewHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", "Title", "Body")
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

	fmt.Printf("Getting ready to serve on port: %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
