package app

import (
	"encoding/json"
	"log"
	"net/http"
)

// ApplicationAPIHandler handles POST requests and mirrors back the received JSON payload.
func ApplicationAPIHandler(w http.ResponseWriter, r *http.Request) {
	// Decode incoming request body into a map
	var payload map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		// Return a 400 status code and an "Invalid JSON" message if decoding fails
		//push alerts
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Log the received payload for debugging
	log.Printf("Received payload: %+v", payload)

	// Set response header as JSON
	w.Header().Set("Content-Type", "application/json")

	// Return the same JSON payload in the response
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		//push alerts
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
	}
}
