package roundrobin

import (
	"errors"
	"log"
	"sync"
)

// RoundRobinInterface defines the methods for RoundRobin
type RoundRobinInterface interface {
	Next() (string, error)
}

// RoundRobin struct holds the list of instances and the current index for round-robin distribution.
type RoundRobin struct {
	instances []string   // List of instances/ports to balance the load across
	index     int        // Current index in the round-robin rotation
	mu        sync.Mutex // Ensure thread-safety for accessing the index
}

// New creates a new instance of RoundRobin with the given list of API instances.
func New(ports []string) *RoundRobin {
	return &RoundRobin{
		instances: ports,
		index:     0,
	}
}

// Next selects the next API instance in a round-robin fashion and ensures thread-safety.
func (rr *RoundRobin) Next() (string, error) {
	rr.mu.Lock()
	defer rr.mu.Unlock()

	if len(rr.instances) == 0 {
		//push alerts
		return "", errors.New("no instances available")
	}

	// Get the next instance in the round-robin cycle
	instance := rr.instances[rr.index]
	log.Printf("Routed the application to the instance  : %s", instance)

	// Update the index for the next request
	rr.index = (rr.index + 1) % len(rr.instances)

	return instance, nil
}
