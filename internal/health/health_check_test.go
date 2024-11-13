package health

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHealthCheckHandler tests the /health endpoint for the server health check.
func TestHealthCheckHandler(t *testing.T) {
	tests := []struct {
		name           string
		expectedStatus string
		expectedCode   int
	}{
		{
			name:           "Successful Health Check",
			expectedStatus: "OK",
			expectedCode:   http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a request to test the HealthCheckHandler
			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			rr := httptest.NewRecorder()

			// Call the handler function
			HealthCheckHandler(rr, req)

			// Check status code
			assert.Equal(t, tt.expectedCode, rr.Code)

			// Check response body
			var response HealthResponse
			err := json.NewDecoder(rr.Body).Decode(&response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, response.Status)
		})
	}
}
