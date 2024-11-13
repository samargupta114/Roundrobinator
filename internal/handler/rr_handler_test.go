package handler

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/samargupta114/Roundrobinator.git/internal/roundrobin"
	"github.com/stretchr/testify/assert"
)

// MockHttpClient is a mock implementation of the httpclient.ClientInterface to mock the ForwardRequest method.
type MockHttpClient struct {
	ForwardRequestFunc func(*http.Request, string) (*http.Response, error)
}

// ForwardRequest is the mocked method for forwarding requests.
func (m *MockHttpClient) ForwardRequest(r *http.Request, url string) (*http.Response, error) {
	return m.ForwardRequestFunc(r, url)
}

func TestRouteHandler(t *testing.T) {
	tests := []struct {
		name           string
		roundRobin     *roundrobin.RoundRobin
		mockClientFunc func(*http.Request, string) (*http.Response, error)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:       "Successful request forwarding",
			roundRobin: roundrobin.New([]string{"8081", "8082"}),
			mockClientFunc: func(r *http.Request, url string) (*http.Response, error) {
				// Simulating a successful response from the backend
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(`{"message": "success"}`))),
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message": "success"}`,
		},
		{
			name:       "Error forwarding request",
			roundRobin: roundrobin.New([]string{"8081", "8082"}),
			mockClientFunc: func(r *http.Request, url string) (*http.Response, error) {
				// Simulating an error while forwarding the request
				return nil, errors.New("forwarding error")
			},
			expectedStatus: http.StatusBadGateway,
			expectedBody:   "Error forwarding request",
		},
		{
			name:       "Error writing response",
			roundRobin: roundrobin.New([]string{"8081", "8082"}),
			mockClientFunc: func(r *http.Request, url string) (*http.Response, error) {
				// Simulating a successful response but error while writing to the response
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(`{"message": "success"}`))),
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message": "success"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock HTTP client that implements ClientInterface
			mockClient := &MockHttpClient{
				ForwardRequestFunc: tt.mockClientFunc,
			}

			// Create a test request
			req := httptest.NewRequest(http.MethodGet, "/test", nil)

			// Create a response recorder to capture the response
			rr := httptest.NewRecorder()

			// Call the RouteHandler
			handler := RouteHandler(tt.roundRobin, mockClient)
			handler.ServeHTTP(rr, req)

			// Assert the status code and response body (trim spaces and newlines)
			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Equal(t, tt.expectedBody, strings.TrimSpace(rr.Body.String()))
		})
	}
}
