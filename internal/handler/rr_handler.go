package handler

import (
	"io"
	"net/http"

	"github.com/samargupta114/Roundrobinator.git/internal/roundrobin"
	"github.com/samargupta114/Roundrobinator.git/pkg/utils/httpclient"
)

// sendErrorResponse sends an error response with the specified message and status code.
func sendErrorResponse(w http.ResponseWriter, msg string, statusCode int) {
	http.Error(w, msg, statusCode) // Send the HTTP error response.
}

// RouteHandler handles forwarding HTTP requests using Round Robin to application instances.
// rr: RoundRobinInterface for selecting the next server instance.
// client: ClientInterface to forward the HTTP request to the chosen instance.
func RouteHandler(rr roundrobin.RoundRobinInterface, client httpclient.ClientInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the next instance/port from the Round Robin mechanism.
		port, err := rr.Next()
		if err != nil {
			// If an error occurred, send a 500 error response.
			sendErrorResponse(w, "Error getting next round-robin instance", http.StatusInternalServerError)
			return
		}

		// Construct the target URL for the request to the chosen instance.
		url := "http://localhost:" + port + "/mirror"

		// Forward the request to the target instance.
		resp, err := client.ForwardRequest(r, url)
		if err != nil {
			// If forwarding fails, send a 502 Bad Gateway error response.
			sendErrorResponse(w, "Error forwarding request", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close() // Ensure the response body is closed after streaming.

		// Copy headers from the target response to the client response.
		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value) // Add each header to the response.
			}
		}

		//_, err = w.Write(body) can also use this to direct write

		// Stream the response body directly to the client.
		if _, err := io.Copy(w, resp.Body); err != nil {
			// If streaming fails, send a 500 error response.
			http.Error(w, "Error streaming response body", http.StatusInternalServerError)
			return
		}
	}
}
