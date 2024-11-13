package app

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestApplicationAPIHandler tests the ApplicationAPIHandler function.
func TestApplicationAPIHandler(t *testing.T) {
	// Define a table of test cases
	tests := []struct {
		name               string
		payload            interface{}
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name: "Valid JSON",
			payload: map[string]interface{}{
				"game":    "Mobile Legends",
				"gamerID": "GYUTDTE",
				"points":  20,
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"game":    "Mobile Legends",
				"gamerID": "GYUTDTE",
				"points":  20,
			},
		},
		{
			name:               "Invalid JSON",
			payload:            `{ "game": "Mobile Legends", "points": }`, // malformed JSON
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   "Invalid JSON\n",
		},
	}

	// Loop through each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal the payload to JSON if valid
			var jsonPayload []byte
			var err error
			if str, ok := tt.payload.(string); ok { // If payload is a raw string (for invalid JSON)
				jsonPayload = []byte(str)
			} else {
				jsonPayload, err = json.Marshal(tt.payload)
				if err != nil {
					t.Fatalf("Error marshaling payload: %v", err)
				}
			}

			// Create a new HTTP request to test the handler
			req := httptest.NewRequest(http.MethodPost, "/api", bytes.NewReader(jsonPayload))
			req.Header.Set("Content-Type", "application/json")

			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()

			// Create a handler function to pass to the test
			handler := http.HandlerFunc(ApplicationAPIHandler)

			// Serve the HTTP request using the handler
			handler.ServeHTTP(rr, req)

			// Check the status code
			if status := rr.Code; status != tt.expectedStatusCode {
				t.Errorf("Expected status %d, got %d", tt.expectedStatusCode, status)
			}

			// Check if the response body matches the expected response
			if rr.Code == http.StatusBadRequest {
				// Read the body for the Bad Request error and compare with expected response
				body := rr.Body.String()
				if body != tt.expectedResponse {
					t.Errorf("Expected response body %v, got %v", tt.expectedResponse, body)
				}
			} else {
				// For valid responses, decode the response and compare it with the expected response
				var responseBody interface{}
				if err := json.NewDecoder(rr.Body).Decode(&responseBody); err != nil {
					t.Fatalf("Error decoding response: %v", err)
				}

				// Compare the response payload with the expected response
				if !compareJSON(tt.expectedResponse, responseBody) {
					t.Errorf("Expected response payload to be %+v, got %+v", tt.expectedResponse, responseBody)
				}
			}
		})
	}
}

// compareJSON compares two maps by marshaling them to JSON and comparing their byte representations
func compareJSON(a, b interface{}) bool {
	aJSON, err := json.Marshal(a)
	if err != nil {
		return false
	}
	bJSON, err := json.Marshal(b)
	if err != nil {
		return false
	}

	return bytes.Equal(aJSON, bJSON)
}
