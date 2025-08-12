package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/The-EpaG/URL-Shortener/internal/api"
	"github.com/The-EpaG/URL-Shortener/internal/storage"
)

func main() {
	const port = 8080

	// Initialize storage
	s := storage.NewSQLiteStorage()
	err := s.Init()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize API handlers with storage dependency
	apiHandlers := api.NewAPI(s)

	// Set up HTTP routes
	http.HandleFunc("/shorten", apiHandlers.CreateShortURL)
	http.HandleFunc("/healthz", apiHandlers.HealthCheck)
	http.HandleFunc("/", apiHandlers.RedirectShortURL)

	log.Default().Println("Server started on :" + fmt.Sprint(port))
	log.Fatal(http.ListenAndServe(":"+fmt.Sprint(port), nil))
}