package physics

import "math"

const (
	Gravity = 9.8

	KStraight = 0.0000166
	KBraking  = 0.0398
	KCorner   = 0.000265

	CornerCrashPenaltyDegradation = 0.1
)

func CalculateTyreFriction(baseFriction, totalDegradation, weatherMultiplier float64) float64 {
	return (baseFriction - totalDegradation) * weatherMultiplier
}

func CalculateMaxCornerSpeed(tyreFriction, radius, crawlConstant float64) float64 {
	return math.Sqrt(tyreFriction*Gravity*radius) + crawlConstant
}

func DegradationStraight(degradationRate, length float64) float64 {
	return degradationRate * length * KStraight
}

func DegradationBraking(initialSpeed, finalSpeed, degradationRate float64) float64 {
	initialScaled := initialSpeed / 100.0
	finalScaled := finalSpeed / 100.0

	speedFactor := (initialScaled * initialScaled) - (finalScaled * finalScaled)

	return speedFactor * KBraking * degradationRate
}

func DegradationCorner(speed, radius, degradationRate float64) float64 {
	speedFactor := speed * speed
	return KCorner * speedFactor * degradationRate / radius
}
