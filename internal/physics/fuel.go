package physics

const (
	KBase = 0.0005
	KDrag = 0.0000000015
)

func FuelUsed(initialSpeed, finalSpeed, distance float64) float64 {
	avgSpeed := (initialSpeed + finalSpeed) / 2.0
	return (KBase + KDrag*(avgSpeed*avgSpeed)) * distance
}
