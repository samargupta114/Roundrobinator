package health

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/samargupta114/Roundrobinator.git/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheckHandler(t *testing.T) {
	// Create a new HTTP request
	req := httptest.NewRequest(http.MethodGet, "/health", nil)

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Call the handler function
	HealthCheckHandler(rr, req)

	// Verify the status code is 200 OK
	assert.Equal(t, http.StatusOK, rr.Code)

	// Verify the content type is application/json
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	// Verify the response body
	var response HealthResponse
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "OK", response.Status)
}

// MockRoundTripper is a custom RoundTripper for mocking HTTP requests
type MockRoundTripper struct {
	ResponseMap map[string]*http.Response
}

// RoundTrip executes a single HTTP transaction and returns a mock response
func (mrt *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Check if the requested URL exists in the mock response map
	if resp, exists := mrt.ResponseMap[req.URL.String()]; exists {
		return resp, nil
	}
	// Return a generic error response if the URL is not found in the mock map
	return &http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       ioutil.NopCloser(bytes.NewBufferString("Server Error")),
	}, nil
}

func TestCheckHealth_WithRoundTripper(t *testing.T) {
	// Create a mock configuration
	cfg := &config.Config{
		Server: config.Server{
			Port: "8080",
		},
		Backend: config.Backend{
			Routes: []string{"9090", "7070"},
		},
	}

	// Define mock responses for specific URLs
	mockResponses := map[string]*http.Response{
		"http://localhost:8080/health": {
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString("OK")),
		},
		"http://localhost:9090/health": {
			StatusCode: http.StatusInternalServerError,
			Body:       ioutil.NopCloser(bytes.NewBufferString("Internal Server Error")),
		},
		"http://localhost:7070/health": {
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString("OK")),
		},
	}

	// Create the mock RoundTripper
	mockRoundTripper := &MockRoundTripper{ResponseMap: mockResponses}

	// Replace the default HTTP client with a mock client
	mockClient := &http.Client{Transport: mockRoundTripper}

	// Backup the original http.DefaultClient and restore it after the test
	originalClient := http.DefaultClient
	http.DefaultClient = mockClient
	defer func() {
		http.DefaultClient = originalClient
	}()

	// Capture log output
	var logOutput bytes.Buffer
	log.SetOutput(&logOutput)
	defer func() { log.SetOutput(nil) }() // Reset log output after the test

	// Call the checkHealth function
	checkHealth(cfg)

	// Verify the expected log output
	assert.Contains(t, logOutput.String(), "Health check succeeded for http://localhost:8080/health")
	assert.Contains(t, logOutput.String(), "Health check succeeded for http://localhost:7070/health")
}
