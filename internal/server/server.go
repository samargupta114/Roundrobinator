package server

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/samargupta114/Roundrobinator.git/internal/config"
	"github.com/samargupta114/Roundrobinator.git/internal/health"
)

const (
	Healthcheck = "healthcheck" // Constant for health check endpoint
)

// Launch initializes and starts both the Application API and Round Robin API.
func Launch(cfg *config.Config) {
	//  WaitGroup to manage concurrent server launches and graceful shutdown
	var wg sync.WaitGroup

	// context that will be canceled when shutdown is triggered
	ctx, cancel := context.WithCancel(context.Background())

	// channel to listen for system signals like interrupt (Ctrl+C) or termination (SIGTERM)
	signalChan := make(chan os.Signal, 1)
	// Notify signalChan on receiving Interrupt or SIGTERM signals
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// List of servers that implement the ServerLauncher interface
	// These servers will be launched concurrently
	servers := []ServerLauncher{
		&ApplicationServer{}, // Application API server
		&RoundRobinServer{},  // Round Robin API server
	}

	// Launch each server in a separate goroutine
	for _, server := range servers {
		wg.Add(1) // Increment the WaitGroup counter for each server launch
		go func(srv ServerLauncher) {
			defer wg.Done()           // Decrement the WaitGroup counter when the server finishes
			srv.Launch(ctx, cfg, &wg) // Start the server, passing the context, config, and WaitGroup
		}(server)
	}

	// Start health check monitoring in a separate goroutine
	go health.StartHealthCheck(cfg, &wg)

	// Start a goroutine to handle graceful shutdown on receiving system signals
	go func() {
		// Wait for a signal (Interrupt or SIGTERM)
		<-signalChan
		// Log the shutdown signal
		log.Println("Shutdown signal received. Gracefully shutting down...")
		// Cancel the context, which will trigger shutdown of servers
		cancel()
	}()

	// Wait for all goroutines (server launches, health check, shutdown) to finish
	wg.Wait()

	// Log that all servers have stopped
	log.Println("All servers stopped. Exiting.")
}
