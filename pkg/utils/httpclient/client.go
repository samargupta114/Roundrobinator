package httpclient

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

// Client struct wraps an HTTP client with a configurable timeout.
type Client struct {
	client *http.Client
}

// ClientInterface defines the methods that our mock and actual client will implement.
type ClientInterface interface {
	ForwardRequest(*http.Request, string) (*http.Response, error)
}

// NewClient creates a new HTTP client with the specified timeout in seconds.
func NewClient(timeout int) *Client {
	return &Client{client: &http.Client{Timeout: time.Duration(timeout) * time.Second}}
}

// ForwardRequest forwards an incoming HTTP request to the target URL and returns the response.
func (c *Client) ForwardRequest(req *http.Request, url string) (*http.Response, error) {
	// Read the body of the incoming request
	body, _ := io.ReadAll(req.Body)

	// Create a new HTTP request with the same method, headers, and body
	newReq, _ := http.NewRequest(req.Method, url, bytes.NewReader(body))
	newReq.Header = req.Header

	// Send the request using the client's HTTP client
	return c.client.Do(newReq)
}
