package handler

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// MockRoundRobin is a mock implementation of RoundRobin for testing.
type MockRoundRobin struct {
	ports   []string
	counter int
	err     error
}

func (m *MockRoundRobin) Next() (string, error) {
	if m.err != nil {
		return "", m.err
	}
	port := m.ports[m.counter%len(m.ports)]
	m.counter++
	return port, nil
}

// MockHttpClient is a mock implementation of ClientInterface for testing.
type MockHttpClient struct {
	resp *http.Response
	err  error
}

func (m *MockHttpClient) ForwardRequest(r *http.Request, url string) (*http.Response, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.resp, nil
}

// TestRouteHandler tests the RouteHandler function using a table-driven approach.
func TestRouteHandler(t *testing.T) {
	tests := []struct {
		name               string
		roundRobinError    error
		forwardRequestResp *http.Response
		forwardRequestErr  error
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:               "Round Robin Next Error",
			roundRobinError:    errors.New("round robin error"),
			forwardRequestResp: nil,
			forwardRequestErr:  nil,
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       "Error getting next round-robin instance\n",
		},
		{
			name:               "Forwarding Request Error",
			roundRobinError:    nil,
			forwardRequestResp: nil,
			forwardRequestErr:  errors.New("forwarding request error"),
			expectedStatusCode: http.StatusBadGateway,
			expectedBody:       "Error forwarding request\n",
		},
		{
			name:            "Successful Forwarding",
			roundRobinError: nil,
			forwardRequestResp: &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(strings.NewReader(`{"message": "Success"}`)),
			},
			forwardRequestErr:  nil,
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"message": "Success"}`,
		},
		{
			name:            "Error in Streaming Response Body",
			roundRobinError: nil,
			forwardRequestResp: &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(&errorReader{}), // Simulate error during body read
			},
			forwardRequestErr:  nil,
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       "Error streaming response body\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up mocks
			mockRoundRobin := &MockRoundRobin{
				ports:   []string{"8081"},
				counter: 0,
				err:     tt.roundRobinError,
			}
			mockHttpClient := &MockHttpClient{
				resp: tt.forwardRequestResp,
				err:  tt.forwardRequestErr,
			}

			// Create a test request
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			rr := httptest.NewRecorder()

			// Create the handler function
			handler := RouteHandler(mockRoundRobin, mockHttpClient)

			// Call the handler
			handler.ServeHTTP(rr, req)

			// Check status code
			if rr.Code != tt.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatusCode, rr.Code)
			}

			// Check response body
			if body := rr.Body.String(); body != tt.expectedBody {
				t.Errorf("Expected response body %q, got %q", tt.expectedBody, body)
			}
		})
	}
}

// errorReader is a mock reader that simulates an error during Read.
type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("error reading body")
}

func (e *errorReader) Close() error {
	return nil
}
