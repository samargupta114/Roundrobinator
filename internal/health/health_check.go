package health

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/samargupta114/Roundrobinator.git/internal/config"
)

// HealthResponse is the structure for the health check response
type HealthResponse struct {
	Status string `json:"status"`
}

// HealthCheckHandler handles the /health endpoint for checking the health of the server
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("Health check requested")

	// Set response headers for JSON
	w.Header().Set("Content-Type", "application/json")

	// Respond with a status indicating that the server is healthy
	response := HealthResponse{Status: "OK"}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode health response", http.StatusInternalServerError)
		return
	}

	//for debugging
	log.Println("Health check response : " + response.Status)
}

// StartHealthCheck performs periodic health checks for backend services
func StartHealthCheck(cfg *config.Config, wg *sync.WaitGroup) {
	defer wg.Done()

	// Set up a ticker to run health checks periodically
	ticker := time.NewTicker(time.Duration(cfg.HealthCheckTickerTimeInSeconds) * time.Second)
	defer ticker.Stop()

	// Run health checks indefinitely in a loop
	for {
		select {
		case <-ticker.C:
			// Perform health checks for the configured services
			checkHealth(cfg)
		}
	}
}

// checkHealth checks the health of the application and round-robin API servers.
func checkHealth(cfg *config.Config) {
	// Define health check URLs based on configuration
	healthURLs := []string{
		"http://localhost:" + cfg.Server.Port + "/health", // Round Robin API
	}

	// Add all application API URLs from config
	for _, port := range cfg.Backend.Routes {
		healthURLs = append(healthURLs, "http://localhost:"+port+"/health")
	}

	// Iterate through all health check URLs and log the results
	for _, url := range healthURLs {
		resp, err := http.Get(url)
		if err != nil && resp.StatusCode != http.StatusOK {
			//push alerts
			log.Printf("Health check failed for %s: %v\n", url, err)
		} else {
			log.Printf("Health check succeeded for %s: %d\n", url, resp.StatusCode)
			resp.Body.Close() // Close the response body to avoid resource leakage
		}
	}
}
