package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	confPath = "ROUND_ROBIN_CONF_PATH" // Environment variable that specifies the config file path
)

// Config holds the overall configuration for the application, including server settings, backend configurations, and health check , graceful shutdown intervals.
type Config struct {
	Server  Server  `json:"server"`  // Server configuration settings
	Backend Backend `json:"backend"` // Backend configuration for API routing and endpoints

	// HealthCheckTickerTimeInSeconds defines the interval for health check ticks in seconds.
	// It specifies how often the system should check the health of the backend services.
	HealthCheckTickerTimeInSeconds int64 `json:"healthCheck_ticker_time_seconds"`

	// GracefulTimeoutSeconds specifies the time allowed for graceful shutdown of the server.
	GracefulTimeoutSeconds int64 `json:"graceful_timeout_seconds"`
}

// Server represents the configuration for the server settings.
type Server struct {
	Port    string `json:"port"`    // The port on which the server should listen
	Timeout int    `json:"timeout"` // The timeout in seconds for server requests
}

// Backend holds the configuration for backend services, including server routes and endpoints.
type Backend struct {
	// Routes is a list of routes (e.g., ports) where the backend services are available.
	Routes []string `json:"routes"`

	// Endpoint is a map of endpoint configurations, where the key is the endpoint name (e.g., "health_check")
	// and the value holds the specific URL and timeout for that endpoint.
	Endpoint map[string]Endpoint `json:"endpoints"`
}

// Endpoint defines the configuration for a single backend endpoint.
type Endpoint struct {
	// URL is the endpoint URL (e.g., "http://localhost:8080/health_check")
	URL string `json:"url"`

	// Timeout specifies the timeout for requests to this endpoint.
	Timeout int `json:"timeout"`
}

// LoadConfig reads the configuration file specified by the environment variable
// ROUND_ROBIN_CONF_PATH, unmarshal its content, and returns the populated Config struct.
func LoadConfig() (*Config, error) {
	// Get the config file path from environment variable
	configPath := os.Getenv(confPath)
	if configPath == "" {
		// push alerts
		return nil, fmt.Errorf("config path not set in environment variable : %s", confPath)
	}

	// Open the config file
	file, err := os.Open(configPath)
	if err != nil {
		// push alerts
		return nil, fmt.Errorf("failed to load config: %v", err)
	}
	defer file.Close()

	// Unmarshal JSON content into Config struct
	var cfg Config
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		// push alerts
		return nil, fmt.Errorf("failed to decode config: %v", err)
	}

	// Return the populated Config struct
	return &cfg, nil
}
