package physics_test

import (
	"math"
	"testing"

	"github.com/osmanyasin/grand-prix/internal/physics"
)

func TestFuelUsed(t *testing.T) {
	tests := []struct {
		name         string
		initialSpeed float64
		finalSpeed   float64
		distance     float64
		expected     float64
	}{
		{
			name:         "Spec Example: 50m/s to 70m/s over 800m",
			initialSpeed: 50.0,
			finalSpeed:   70.0,
			distance:     800.0,
			expected:     0.40432,
		},
		{
			name:         "Zero Distance uses zero fuel",
			initialSpeed: 50.0,
			finalSpeed:   70.0,
			distance:     0.0,
			expected:     0.0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := physics.FuelUsed(tc.initialSpeed, tc.finalSpeed, tc.distance)

			if math.Abs(got-tc.expected) > 1e-9 {
				t.Errorf("FuelUsed(%v, %v, %v) = %v; want %v",
					tc.initialSpeed, tc.finalSpeed, tc.distance, got, tc.expected)
			}
		})
	}
}
