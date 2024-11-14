package server

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/samargupta114/Roundrobinator.git/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRoundRobin struct {
	mock.Mock
}

func (m *MockRoundRobin) Next() string {
	args := m.Called()
	return args.String(0)
}

func TestLaunch(t *testing.T) {
	// Setting up a mock configuration
	cfg := &config.Config{
		Server: config.Server{
			Port:    "8080",
			Timeout: 5,
		},
		Backend: config.Backend{
			Routes: []string{"http://localhost:8081", "http://localhost:8082"},
			Endpoint: map[string]config.Endpoint{
				"healthcheck": {URL: "/health"},
			},
		},
		GracefulTimeoutSeconds: 5,
	}

	// Create mock instances for roundrobin and httpclient
	mockRoundRobin := new(MockRoundRobin)
	// Mock the Next method of roundrobin
	mockRoundRobin.On("Next").Return("http://localhost:8081").Once()

	// Capture log output for verification
	var logOutput bytes.Buffer
	log.SetOutput(&logOutput)
	defer func() { log.SetOutput(nil) }() // Reset log output after the test

	// Create the RoundRobinServer instance
	server := &RoundRobinServer{}

	// Create context with cancellation for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)

	// Start Launch method in a goroutine
	go func() {
		defer wg.Done()
		err := server.Launch(ctx, cfg, &wg)
		assert.NoError(t, err, "Server should start without error")
	}()

	// Allow some time for the server to initialize
	time.Sleep(1 * time.Second)

	// Now, simulate an HTTP request to the /route endpoint to trigger Next()
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:8080/route", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Test that the expected logs were captured
	assert.Contains(t, logOutput.String(), "Round Robin API running on http://localhost:8080/route", "Expected server start log")

	// Simulate stopping the server by canceling the context
	cancel()

	// Wait for the goroutine to finish
	wg.Wait()

	// Assert that the graceful shutdown log appears
	assert.Contains(t, logOutput.String(), "Shutting down Round Robin API on port 8080...", "Expected graceful shutdown log")

}
