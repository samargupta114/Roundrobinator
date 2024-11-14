package server

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/samargupta114/Roundrobinator.git/internal/health"

	"github.com/samargupta114/Roundrobinator.git/internal/config"
	"github.com/samargupta114/Roundrobinator.git/internal/handler"
)

// ServerInterface defines the methods we care about for mocking
type ServerInterface interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

// ServerWrapper is a wrapper around http.Server
type ServerWrapper struct {
	*http.Server
}

// ApplicationServer implements the ServerLauncher interface for Application APIs.
type ApplicationServer struct{}

// Launch starts multiple Application API instances.
func (as *ApplicationServer) Launch(ctx context.Context, cfg *config.Config, wg *sync.WaitGroup) error {
	// Loop through each route in the configuration and launch a server on each port.
	for _, port := range cfg.Backend.Routes {
		wg.Add(1)                                // Add to the WaitGroup for each server to be launched concurrently.
		go as.startAppServer(ctx, port, cfg, wg) // Start the server in a separate goroutine.
	}
	return nil
}

// startAppServer launches a single instance of the Application API.
func (as *ApplicationServer) startAppServer(ctx context.Context, port string, cfg *config.Config, wg *sync.WaitGroup) {
	defer wg.Done() // Decrement the WaitGroup counter when this function completes.

	// Create a new ServeMux to handle HTTP routes
	mux := http.NewServeMux()
	// Register the `/health` route to monitor the health of the server
	mux.HandleFunc(cfg.Backend.Endpoint[Healthcheck].URL, health.HealthCheckHandler)
	// Register the `/mirror` route to handle the main API functionality
	mux.HandleFunc("/mirror", handler.ApplicationAPIHandler)

	// Create a new HTTP server instance with the specified port and handler
	server := &ServerWrapper{Server: &http.Server{Addr: ":" + port, Handler: mux}}

	// Start the HTTP server and handle graceful shutdown
	go as.runServer(ctx, server, port, cfg)
}

func (as *ApplicationServer) runServer(ctx context.Context, server *ServerWrapper, port string, cfg *config.Config) {
	// Log to indicate the server is starting
	log.Printf("Starting Application API server on http://localhost:%s/mirror", port)

	// Create a channel to report any errors encountered by the server
	errCh := make(chan error)

	// Start the server in a new goroutine
	go func() {
		// Listen and serve the application API, report errors through the error channel
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err // Send the error to the error channel
		}
	}()

	// Start another goroutine to listen for the shutdown signal (context cancellation)
	go func() {
		<-ctx.Done() // Wait for the cancellation signal
		// Once canceled, start the graceful shutdown with a timeout
		shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.GracefulTimeoutSeconds)*time.Second)
		defer cancel()

		// Attempt to shut down the server gracefully, and report errors if any
		if err := server.Shutdown(shutdownCtx); err != nil {
			errCh <- err // Send the error to the error channel
		}
		close(errCh) // Close the error channel after attempting shutdown
	}()

	// Block until an error occurs (server failure) or the shutdown completes
	if err := <-errCh; err != nil {
		// Log failure if the server encountered an error
		log.Printf("Application API server failed on port %s: %v", port, err)
	} else {
		// Log successful shutdown
		log.Printf("Successfully shut down Application API on port %s", port)
	}
}
