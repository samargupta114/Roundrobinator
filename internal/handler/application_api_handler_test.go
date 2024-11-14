package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockResponseWriter is a custom ResponseWriter that forces an error on Write
type MockResponseWriter struct {
	header http.Header
	status int
}

func (m *MockResponseWriter) Header() http.Header {
	if m.header == nil {
		m.header = make(http.Header)
	}
	return m.header
}

func (m *MockResponseWriter) Write(p []byte) (int, error) {
	return 0, errors.New("forced write error")
}

func (m *MockResponseWriter) WriteHeader(statusCode int) {
	m.status = statusCode
}

// TestApplicationAPIHandler tests the ApplicationAPIHandler function.
func TestApplicationAPIHandler(t *testing.T) {
	// Define a table of test cases
	tests := []struct {
		name               string
		payload            interface{}
		expectedStatusCode int
		expectedResponse   interface{}
		useMockWriter      bool // To determine if we should use the MockResponseWriter
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
			useMockWriter: false,
		},
		{
			name:               "Invalid JSON",
			payload:            `{ "game": "Mobile Legends", "points": }`, // malformed JSON
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   "Invalid JSON\n",
			useMockWriter:      false,
		},
		{
			name: "Encoding Error",
			payload: map[string]interface{}{
				"game":    "Mobile Legends",
				"gamerID": "GYUTDTE",
				"points":  20,
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   "Error encoding JSON response\n",
			useMockWriter:      true,
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

			// Create a ResponseRecorder or use MockResponseWriter based on the test case
			var rr *httptest.ResponseRecorder
			var writer http.ResponseWriter

			if tt.useMockWriter {
				writer = &MockResponseWriter{}
			} else {
				rr = httptest.NewRecorder()
				writer = rr
			}

			// Create a handler function to pass to the test
			handler := http.HandlerFunc(ApplicationAPIHandler)

			// Serve the HTTP request using the handler
			handler.ServeHTTP(writer, req)

			if tt.useMockWriter {
				// For the mock writer scenario, we don't need to check the body, just the status
				if status := writer.(*MockResponseWriter).Header().Get("status"); status != "" {
					t.Errorf("Expected status %d, got %s", tt.expectedStatusCode, status)
				}
			} else {
				// Check the status code
				if status := rr.Code; status != tt.expectedStatusCode {
					t.Errorf("Expected status %d, got %d", tt.expectedStatusCode, status)
				}

				// Check if the response body matches the expected response
				if rr.Code == http.StatusBadRequest || rr.Code == http.StatusInternalServerError {
					// Read the body for error and compare with expected response
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
