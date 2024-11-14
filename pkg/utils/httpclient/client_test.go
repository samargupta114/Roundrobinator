package httpclient

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRoundTripper implements RoundTripper for mocking HTTP requests
type MockRoundTripper struct {
	mock.Mock
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestForwardRequest(t *testing.T) {
	tests := []struct {
		name          string
		method        string
		body          string
		mockResponse  *http.Response
		mockError     error
		expectedError bool
		expectedBody  string
	}{
		{
			name:         "Successful request forwarding",
			method:       http.MethodPost,
			body:         `{"key": "value"}`,
			mockResponse: &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewReader([]byte(`{"key": "value"}`)))},
			mockError:    nil,
			expectedBody: `{"key": "value"}`,
		},
		{
			name:          "Client error during request",
			method:        http.MethodGet,
			body:          "",
			mockResponse:  nil,
			mockError:     errors.New("network error"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new HTTP request
			req := httptest.NewRequest(tt.method, "http://localhost", bytes.NewReader([]byte(tt.body)))

			// Mock RoundTripper
			mockRoundTripper := &MockRoundTripper{}
			mockRoundTripper.On("RoundTrip", mock.Anything).Return(tt.mockResponse, tt.mockError)

			// Create a client with mocked RoundTripper
			httpClient := &http.Client{Transport: mockRoundTripper, Timeout: time.Second * 10}
			client := &Client{client: httpClient}

			// Call the ForwardRequest method
			resp, err := client.ForwardRequest(req, "http://mocked-url.com")

			// Assert based on expected results
			if tt.expectedError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.mockResponse.StatusCode, resp.StatusCode)

			// Check the response body
			respBody, _ := io.ReadAll(resp.Body)
			assert.Equal(t, tt.expectedBody, string(respBody))

			// Verify the expectations
			mockRoundTripper.AssertExpectations(t)
		})
	}
}

func TestNewClient(t *testing.T) {
	tests := []struct {
		name            string
		timeoutSeconds  int
		expectedTimeout time.Duration
	}{
		{
			name:            "Client with 10 seconds timeout",
			timeoutSeconds:  10,
			expectedTimeout: 10 * time.Second,
		},
		{
			name:            "Client with 0 seconds timeout",
			timeoutSeconds:  0,
			expectedTimeout: 0 * time.Second,
		},
		{
			name:            "Client with 30 seconds timeout",
			timeoutSeconds:  30,
			expectedTimeout: 30 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.timeoutSeconds)

			// Check that the returned client is not nil
			assert.NotNil(t, client)
			assert.NotNil(t, client.client)

			// Check that the client timeout is set correctly
			assert.Equal(t, tt.expectedTimeout, client.client.Timeout)
		})
	}
}
