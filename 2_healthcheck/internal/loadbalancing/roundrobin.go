package loadbalancing

import (
	"errors"
	"workspace/tinybandaid/internal/models"
)

type RoundRobin struct{}

// RoundRobin returns a pointer to the next available backend using round robin method
// This should be called in a resource-managing thread
func (l *RoundRobin) GetAvailableIndex(currentIndex int, backends []*models.Backend) (int, error) {
	// Round robin loadbalancing
	for i := 0; i < len(backends); i++ {
		// rotate index
		currentIndex = (currentIndex + 1) % len(backends)
		if backends[currentIndex].IsAlive {
			return currentIndex, nil
		}
	}
	return currentIndex, errors.New("no backend servers are available")
}
