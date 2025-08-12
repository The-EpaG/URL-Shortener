package api

import (
	"crypto/sha1"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/The-EpaG/URL-Shortener/internal/storage"
)

type API struct {
	storage storage.Storage
}

func NewAPI(s storage.Storage) *API {
	return &API{storage: s}
}

func (a *API) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		OriginalURL string `json:"original_url"`
	}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := generateID(data.OriginalURL)

	shortURL, err := a.storage.GetShortURL(id)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err = a.storage.CreateShortURL(id, data.OriginalURL)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			shortURL = &storage.ShortURL{ID: id, Original: data.OriginalURL, AccessCount: 0}
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	response := shortURL

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	log.Default().Println("Short URL created successfully for " + data.OriginalURL + ".")
}

func (a *API) RedirectShortURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Path[1:] // Remove leading slash

	originalURL, err := a.storage.GetOriginalURL(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Short URL not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = a.storage.IncrementAccessCount(id)
	if err != nil {
		log.Printf("Failed to update access count for %s: %v", id, err)
	}

	http.Redirect(w, r, originalURL, http.StatusFound)
	log.Default().Println(r.RequestURI + " redirected to " + originalURL + ".")
}

func generateID(originalUrl string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(originalUrl)))[:8]
}

func (a *API) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check database connection
	if err := a.storage.Ping(); err != nil {
		log.Printf("Health check failed: database ping error: %v", err)
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
