package internal

import (
	"database/sql"
	"log"
	"net/http"
	"time"
)

// HealthCheckResult contains the outcome of a single health check
type HealthCheckResult struct {
	SiteID         int
	IsUp           bool
	StatusCode     int
	ResponseTimeMs int64
}

// StartChecker initializes the health checking process.
func StartChecker(db *sql.DB, publisher *Publisher) {
	log.Println("Starting health checker...")
	// Run a check immediately on start, then on a ticker
	go checkAllSites(db, publisher)

	ticker := time.NewTicker(1 * time.Minute) // Check every 1 minute
	go func() {
		for range ticker.C {
			checkAllSites(db, publisher)
		}
	}()
}

// checkAllSites fetches all sites from the DB and checks them concurrently
func checkAllSites(db *sql.DB, publisher *Publisher) {
	log.Println("Running health checks for all sites...")
	rows, err := db.Query("SELECT id, url FROM sites")
	if err != nil {
		log.Printf("Error fetching sites for checking: %v", err)
		return
	}
	defer rows.Close()

	resultsChan := make(chan HealthCheckResult)
	siteCount := 0

	for rows.Next() {
		var s Site
		if err := rows.Scan(&s.ID, &s.URL); err != nil {
			log.Printf("Error scanning site row for checking: %v", err)
			continue
		}
		siteCount++
		go performCheck(s, resultsChan)
	}

	// Collect results
	for i := 0; i < siteCount; i++ {
		result := <-resultsChan
		saveCheckResult(db, result)

		if err := publisher.Publish(result); err != nil {
			log.Printf("Failed to publish message for site ID %d: %v", result.SiteID, err)
		}
	}
	log.Println("Finished health checks.")
}

// performCheck executes an HTTP GET request to the site's URL
func performCheck(s Site, c chan<- HealthCheckResult) {
	start := time.Now()

	// Create a client with a timeout
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(s.URL)
	duration := time.Since(start).Milliseconds()

	result := HealthCheckResult{
		SiteID:         s.ID,
		ResponseTimeMs: duration,
	}

	if err != nil {
		log.Printf("Error checking %s: %v", s.URL, err)
		result.IsUp = false
		result.StatusCode = 0 // Or some other indicator of failure
	} else {
		defer resp.Body.Close()
		result.IsUp = resp.StatusCode >= 200 && resp.StatusCode < 300
		result.StatusCode = resp.StatusCode
	}

	c <- result
}

// saveCheckResult inserts a health check result into the database.
func saveCheckResult(db *sql.DB, result HealthCheckResult) {
	_, err := db.Exec(
		"INSERT INTO health_checks (site_id, is_up, status_code, response_time_ms) VALUES ($1, $2, $3, $4)",
		result.SiteID, result.IsUp, result.StatusCode, result.ResponseTimeMs,
	)
	if err != nil {
		log.Printf("Error saving health check for site ID %d: %v", result.SiteID, err)
	}
}
