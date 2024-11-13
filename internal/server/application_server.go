package server

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/samargupta114/Roundrobinator.git/internal/app"
	"github.com/samargupta114/Roundrobinator.git/internal/config"
	"github.com/samargupta114/Roundrobinator.git/internal/health"
)

// ApplicationServer implements the ServerLauncher interface for Application APIs.
type ApplicationServer struct{}

// Launch starts multiple Application API instances.
func (as *ApplicationServer) Launch(ctx context.Context, cfg *config.Config, wg *sync.WaitGroup) error {
	// Launch the application servers on different ports as configured
	for _, port := range cfg.Backend.Routes {
		wg.Add(1)
		go func(port string) {
			defer wg.Done()
			// Start a single application server
			if err := startAppServer(ctx, port, cfg); err != nil {
				//push alerts
				log.Printf("Error starting Application API on port %s: %v", port, err)
			}
		}(port)
	}
	return nil
}

// startAppServer launches a single instance of the Application API.
func startAppServer(ctx context.Context, port string, cfg *config.Config) error {
	mux := http.NewServeMux()
	// Register /health route for health check
	mux.HandleFunc(cfg.Backend.Endpoint[Healthcheck].URL, health.HealthCheckHandler)
	// Register the /mirror route
	mux.HandleFunc("/mirror", app.ApplicationAPIHandler)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// Start the server in a separate goroutine so that it can be gracefully shut down
	go func() {
		log.Printf("Starting Application API server on http://localhost:%s/mirror", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			//push alerts
			log.Printf("Application API server failed on port %s: %v", port, err)
		}
	}()

	// Gracefully shutdown the server when the context is canceled
	go func() {
		<-ctx.Done() // Wait for cancellation signal
		log.Printf("Shutting down Application API on port %s...", port)
		shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.GracefulTimeoutSeconds)*time.Second)
		defer cancel()

		// Attempt graceful shutdown
		if err := server.Shutdown(shutdownCtx); err != nil {
			//push alerts
			log.Printf("Graceful shutdown failed on port %s: %v", port, err)
		} else {
			log.Printf("Successfully shut down Application API on port %s", port)
		}
	}()

	// Block here until the server shuts down (returns)
	<-ctx.Done()
	return nil
}
