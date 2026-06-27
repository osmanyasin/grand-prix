package simulator

type CarState struct {
	// --- Physics & Positional Metrics ---
	CurrentSpeedMPS   float64
	TotalTimeSeconds  float64
	CurrentFuelLitres float64

	// --- Tyre Status ---
	ActiveTyreID     int
	ActiveTyreDegrad float64

	// --- Penalty Condition Flags ---
	IsLimpMode  bool
	IsCrawlMode bool

	// --- Accumulators for Final Scoring ---
	TotalFuelUsedLitres  float64
	TotalTyreDegradation float64
	NumberOfBlowouts     int
}

func NewCarState(initialFuel float64, startingTyreID int) *CarState {
	return &CarState{
		CurrentSpeedMPS:   0.0,
		TotalTimeSeconds:  0.0,
		CurrentFuelLitres: initialFuel,
		ActiveTyreID:      startingTyreID,
	}
}
