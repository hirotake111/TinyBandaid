package loadbalancing_test

import (
	"errors"
	"testing"
	"workspace/tinybandaid/internal/loadbalancing"
	"workspace/tinybandaid/internal/models"
)

func TestRoundRobin(t *testing.T) {
	backends := []*models.Backend{
		{IsAlive: true},
		{IsAlive: true},
		{IsAlive: true},
		{IsAlive: false},
	}
	cases := []struct {
		description string
		index       int
		backends    []*models.Backend
		expected    int
		expectedErr error
	}{
		{
			description: "next server is available",
			index:       0,
			backends:    backends,
			expected:    1,
			expectedErr: nil,
		},
		{
			description: "next server is down",
			index:       2,
			backends:    backends,
			expected:    0,
			expectedErr: nil,
		},
		{
			description: "no backend servers are available",
			index:       0,
			backends: []*models.Backend{
				{IsAlive: false},
				{IsAlive: false},
			},
			expected:    0,
			expectedErr: errors.New("no backend servers are available"),
		},
	}

	for _, tt := range cases {
		t.Run(tt.description, func(t *testing.T) {
			rr := loadbalancing.RoundRobin{}
			result, err := rr.GetAvailableIndex(tt.index, tt.backends)
			if result != tt.expected {
				t.Errorf("Expected %d but got %d", tt.expected, result)
			}
			if err != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("Expected error is %s but got %s", tt.expectedErr, err)
			}
		})
	}
}
