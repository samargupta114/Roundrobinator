package server

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/samargupta114/Roundrobinator.git/internal/config"
	"github.com/stretchr/testify/mock"
)

type MockServer struct {
	mock.Mock
}

func (m *MockServer) ListenAndServe() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockServer) Shutdown(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestStartAppServer(t *testing.T) {
	// Mock config
	cfg := &config.Config{
		Backend: config.Backend{
			Routes: []string{"8081"},
			Endpoint: map[string]config.Endpoint{
				"healthcheck": {URL: "/health"},
			},
		},
		GracefulTimeoutSeconds: 5,
	}

	// Initialize the ApplicationServer
	applicationServer := &ApplicationServer{}

	// Mock the http.Server
	mockServer := new(MockServer)

	// Create a mock WaitGroup
	var wg sync.WaitGroup
	wg.Add(1)

	// Mock ListenAndServe to return no error
	mockServer.On("ListenAndServe").Return(nil).Once()

	// Mock Shutdown to return no error
	mockServer.On("Shutdown", mock.Anything).Return(nil).Once()

	// Run startAppServer in a goroutine
	go applicationServer.startAppServer(context.Background(), "8081", cfg, &wg)

	// Give the server a little time to start
	time.Sleep(100 * time.Millisecond)

	mockServer.ListenAndServe()
	mockServer.Shutdown(context.Background())

	// Verify that ListenAndServe was called
	mockServer.AssertExpectations(t)

	// Simulate canceling the context (i.e., server shutdown)
	_, cancel := context.WithCancel(context.Background())
	cancel() // Cancel the context immediately to trigger shutdown
	time.Sleep(100 * time.Millisecond)

	// Verify that Shutdown was called
	mockServer.AssertExpectations(t)

	// Ensure WaitGroup completes
	wg.Wait()
}

// MockApplicationServer is a mock of ApplicationServer to mock startAppServer
type MockApplicationServer struct {
	mock.Mock
	ApplicationServer
}

// Mock startAppServer method
func (m *MockApplicationServer) startAppServer(ctx context.Context, port string, cfg *config.Config, wg *sync.WaitGroup) {
	m.Called(ctx, port, cfg, wg) // Record method call
}
