package roundrobin

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRoundRobin checks that the round-robin distribution works as expected.
func TestRoundRobin(t *testing.T) {
	tests := []struct {
		name         string
		instances    []string
		expectedErr  error
		expectedInst []string
	}{
		{
			name:         "Valid round-robin distribution",
			instances:    []string{"localhost:8081", "localhost:8082", "localhost:8083"},
			expectedErr:  nil,
			expectedInst: []string{"localhost:8081", "localhost:8082", "localhost:8083", "localhost:8081", "localhost:8082", "localhost:8083"},
		},
		{
			name:         "No instances available",
			instances:    []string{},
			expectedErr:  errors.New("no instances available"),
			expectedInst: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := New(tt.instances)

			// If no instances, check the error
			if len(tt.instances) == 0 {
				_, err := rr.Next()
				assert.Equal(t, tt.expectedErr, err)
				return
			}

			// Test round-robin distribution
			for i := 0; i < len(tt.expectedInst); i++ {
				instance, err := rr.Next()
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedInst[i], instance, "Round-robin should return the instances in a cyclic order.")
			}
		})
	}
}
