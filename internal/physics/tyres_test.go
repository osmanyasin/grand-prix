package physics_test

import (
	"math"
	"testing"

	"github.com/osmanyasin/grand-prix/internal/physics"
)

func TestCalculateTyreFriction(t *testing.T) {
	baseFriction := 1.8
	totalDegradation := 0.5
	weatherMultiplier := 1.0
	expected := 1.3

	got := physics.CalculateTyreFriction(baseFriction, totalDegradation, weatherMultiplier)

	if math.Abs(got-expected) > 1e-9 {
		t.Errorf("CalculateTyreFriction() = %v; want %v", got, expected)
	}
}

func TestCalculateMaxCornerSpeed(t *testing.T) {
	tyreFriction := 0.9
	radius := 50.0
	expected := 21.0
	crawlConstant := 10.0

	got := physics.CalculateMaxCornerSpeed(tyreFriction, radius, crawlConstant)

	if math.Abs(got-expected) > 1e-9 {
		t.Errorf("CalculateMaxCornerSpeed() = %v; want %v", got, expected)
	}
}
