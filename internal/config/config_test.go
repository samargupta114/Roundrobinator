package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestLoadConfig tests the LoadConfig function for various cases using table-driven tests.
func TestLoadConfig(t *testing.T) {
	// Define test cases
	tests := []struct {
		name           string
		configContent  string
		expectedError  bool
		expectedConfig *Config
	}{
		{
			name: "ValidConfig",
			configContent: `{
				"server": {
					"port": "8080",
					"timeout": 30
				},
				"backend": {
					"routes": ["8081", "8082"],
					"endpoints": {
						"health_check": {
							"url": "http://localhost:8080/health",
							"timeout": 10
						}
					}
				},
				"healthCheck_ticker_time_seconds": 30,
				"graceful_timeout_seconds": 15
			}`,
			expectedError: false,
			expectedConfig: &Config{
				Server: Server{
					Port:    "8080",
					Timeout: 30,
				},
				Backend: Backend{
					Routes: []string{"8081", "8082"},
					Endpoint: map[string]Endpoint{
						"health_check": {
							URL:     "http://localhost:8080/health",
							Timeout: 10,
						},
					},
				},
				HealthCheckTickerTimeInSeconds: 30,
				GracefulTimeoutSeconds:         15,
			},
		},
		{
			name: "InvalidJSON",
			configContent: `{
				"server": {
					"port": "8080",
					"timeout": 30
				},
				"backend": {
					"routes": ["8081", "8082"],
					"endpoints": {
						"health_check": {
							"url": "http://localhost:8080/health",
							"timeout": 10
						}
					}
				},
				"healthCheck_ticker_time_seconds": 30,
				"graceful_timeout_seconds": 15
			`,
			expectedError:  true,
			expectedConfig: nil,
		},
		{
			name: "MissingConfigPathEnv",
			// Empty config content as we won't be using it here
			configContent:  "",
			expectedError:  true,
			expectedConfig: nil,
		},
	}

	// Run all test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup a temporary config file if the test requires it
			var file *os.File
			if tt.configContent != "" {
				var err error
				file, err = os.CreateTemp("", "app-config.json")
				assert.NoError(t, err)
				defer os.Remove(file.Name())

				// Write the content to the file
				_, err = file.WriteString(tt.configContent)
				assert.NoError(t, err)
				err = file.Close()
				assert.NoError(t, err)

				// Set the environment variable for the config path
				os.Setenv(confPath, file.Name())
				defer os.Unsetenv(confPath)
			} else {
				// Test case for missing confPath environment variable
				os.Unsetenv(confPath)
			}

			// Call LoadConfig
			cfg, err := LoadConfig()

			// Validate error and config
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, cfg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, cfg)

				// Validate fields of expectedConfig
				assert.Equal(t, tt.expectedConfig.Server.Port, cfg.Server.Port)
				assert.Equal(t, tt.expectedConfig.Server.Timeout, cfg.Server.Timeout)
				assert.Equal(t, tt.expectedConfig.HealthCheckTickerTimeInSeconds, cfg.HealthCheckTickerTimeInSeconds)
				assert.Equal(t, tt.expectedConfig.GracefulTimeoutSeconds, cfg.GracefulTimeoutSeconds)
				assert.Len(t, cfg.Backend.Routes, len(tt.expectedConfig.Backend.Routes))
				for i, route := range tt.expectedConfig.Backend.Routes {
					assert.Equal(t, route, cfg.Backend.Routes[i])
				}
				assert.Equal(t, tt.expectedConfig.Backend.Endpoint["health_check"].URL, cfg.Backend.Endpoint["health_check"].URL)
			}
		})
	}
}
