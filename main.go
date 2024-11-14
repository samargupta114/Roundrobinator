package main

import (
	"log"

	"github.com/samargupta114/Roundrobinator.git/internal/config"
	"github.com/samargupta114/Roundrobinator.git/internal/server"
)

// Main function is the entry point of the application.
// It starts the server using the Launch method.
func main() {

	// Load the configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		//push alerts
		log.Fatalf("Failed to load config: %v", err)
	}

	// Launch the servers
	server.Launch(cfg)
}
