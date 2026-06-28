package physics

func TimeToAccelerate(initialSpeed, finalSpeed, acceleration float64) float64 {
	return (finalSpeed - initialSpeed) / acceleration
}

func CalculateDistance(initialSpeed, finalSpeed, acceleration float64) float64 {
	return ((finalSpeed * finalSpeed) - (initialSpeed * initialSpeed)) / (2 * acceleration)
}

func DistanceGivenTime(initialSpeed, time, acceleration float64) float64 {
	return (initialSpeed * time) + (0.5 * acceleration * time * time)
}
