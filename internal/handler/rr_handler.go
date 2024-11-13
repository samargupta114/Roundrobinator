package handler

import (
	"io"
	"net/http"

	"github.com/samargupta114/Roundrobinator.git/internal/roundrobin"
	"github.com/samargupta114/Roundrobinator.git/pkg/utils/httpclient"
)

// RouteHandler handles forwarding HTTP requests using Round Robin to application instances
func RouteHandler(rr *roundrobin.RoundRobin, client httpclient.ClientInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the next instance/port from the Round Robin
		port, _ := rr.Next()

		// Construct the full URL for forwarding the request
		url := "http://localhost:" + port + "/mirror"

		// Forward the request to the chosen instance/port
		resp, err := client.ForwardRequest(r, url)
		if err != nil {
			//push alerts
			http.Error(w, "Error forwarding request", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close() // Close the response body after reading

		// Set the response status code
		w.WriteHeader(resp.StatusCode)

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			//push alerts
			http.Error(w, "Error reading response body", http.StatusInternalServerError)
			return
		}

		// Write the response body to the client
		_, err = w.Write(body)
		if err != nil {
			//push alerts
			http.Error(w, "Error writing response", http.StatusInternalServerError)
			return
		}
	}
}
