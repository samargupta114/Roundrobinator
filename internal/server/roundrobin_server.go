package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/samargupta114/Roundrobinator.git/internal/config"
	"github.com/samargupta114/Roundrobinator.git/internal/handler"
	"github.com/samargupta114/Roundrobinator.git/internal/health"
	"github.com/samargupta114/Roundrobinator.git/internal/roundrobin"
	"github.com/samargupta114/Roundrobinator.git/pkg/utils/httpclient"
)

// RoundRobinServer implements the ServerLauncher interface for Round Robin API.
type RoundRobinServer struct{}

// Launch starts the Round Robin API server.
func (rrs *RoundRobinServer) Launch(ctx context.Context, cfg *config.Config, wg *sync.WaitGroup) error {
	mux := http.NewServeMux()

	// Create a round-robin instance to distribute requests to backend servers
	rr := roundrobin.New(cfg.Backend.Routes)
	client := httpclient.NewClient(cfg.Server.Timeout)

	// Healthcheck endpoint
	mux.HandleFunc(cfg.Backend.Endpoint[Healthcheck].URL, health.HealthCheckHandler)

	// Route for handling round-robin logic
	mux.HandleFunc("/route", handler.RouteHandler(rr, client))

	server := &http.Server{
		Addr:    ":" + cfg.Server.Port, // Ensure this is the correct port
		Handler: mux,
	}

	// Graceful shutdown
	go func() {
		<-ctx.Done()
		log.Printf("Shutting down Round Robin API on port %s...", cfg.Server.Port)
		shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.GracefulTimeoutSeconds)*time.Second)
		defer cancel()
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			return
		}
	}()

	log.Printf("Round Robin API running on http://localhost:%s/route", cfg.Server.Port)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		//push alerts
		log.Fatalf("Round Robin server failed: %v", err)
	}

	return nil
}
