package physics_test

import (
	"math"
	"testing"

	"github.com/osmanyasin/grand-prix/internal/physics"
)

func TestTimeToAccelerate(t *testing.T) {
	tests := []struct {
		name         string
		initialSpeed float64
		finalSpeed   float64
		acceleration float64
		expected     float64
	}{
		{
			name:         "Accelerate 0 to 70 at 10m/s2",
			initialSpeed: 0.0,
			finalSpeed:   70.0,
			acceleration: 10.0,
			expected:     7.0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := physics.TimeToAccelerate(tc.initialSpeed, tc.finalSpeed, tc.acceleration)
			if math.Abs(got-tc.expected) > 1e-9 {
				t.Errorf("TimeToAccelerate() = %v; want %v", got, tc.expected)
			}
		})
	}
}
