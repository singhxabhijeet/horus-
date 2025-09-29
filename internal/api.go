package internal

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

// API holds the database connection
type API struct {
	db *sql.DB
}

// NewAPI creates a new API handler
func NewAPI(db *sql.DB) *API {
	return &API{db: db}
}

// AddSiteRequest is the structure of the request body for adding a site
type AddSiteRequest struct {
	URL string `json:"url"`
}

// Site represents a monitored website
type Site struct {
	ID        int    `json:"id"`
	URL       string `json:"url"`
	CreatedAt string `json:"created_at"`
}

// AddSiteHandler handles requests to add a new site to monitor
func (a *API) AddSiteHandler(w http.ResponseWriter, r *http.Request) {
	var req AddSiteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	var siteID int
	err := a.db.QueryRow("INSERT INTO sites (url) VALUES ($1) RETURNING id", req.URL).Scan(&siteID)
	if err != nil {
		log.Printf("Error inserting site: %v", err)
		http.Error(w, "Failed to add site", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": siteID})
}

// ListSitesHandler handles requests to list all monitored sites
func (a *API) ListSitesHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := a.db.Query("SELECT id, url, created_at FROM sites ORDER BY created_at DESC")
	if err != nil {
		log.Printf("Error querying sites: %v", err)
		http.Error(w, "Failed to retrieve sites", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	sites := []Site{}
	for rows.Next() {
		var s Site
		if err := rows.Scan(&s.ID, &s.URL, &s.CreatedAt); err != nil {
			log.Printf("Error scanning site row: %v", err)
			continue
		}
		sites = append(sites, s)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sites)
}

// RegisterRoutes sets up the API routes
func (a *API) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/sites", a.AddSiteHandler)
	mux.HandleFunc("GET /api/sites", a.ListSitesHandler)
}
