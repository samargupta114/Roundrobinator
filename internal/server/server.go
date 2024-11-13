package server

import (
	"context"
	"log"
	_ "net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/samargupta114/Roundrobinator.git/internal/config"
	"github.com/samargupta114/Roundrobinator.git/internal/health"
)

const (
	Healthcheck = "healthcheck"
)

// Launch initializes and starts both the Application API and Round Robin API.
func Launch(cfg *config.Config) {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// List of servers implementing ServerLauncher interface
	servers := []ServerLauncher{
		&ApplicationServer{},
		&RoundRobinServer{},
	}

	// Launch all servers
	for _, server := range servers {
		wg.Add(1)
		go func(srv ServerLauncher) {
			defer wg.Done()
			srv.Launch(ctx, cfg, &wg)
		}(server)
	}

	// Start health check monitoring
	go health.StartHealthCheck(cfg, &wg)

	// Graceful shutdown
	go func() {
		<-signalChan
		log.Println("Shutdown signal received. Gracefully shutting down...")
		cancel()
	}()

	// Wait for all goroutines to finish
	wg.Wait()
	log.Println("All servers stopped. Exiting.")
}
