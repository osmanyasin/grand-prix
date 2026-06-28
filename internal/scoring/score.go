package scoring

import (
	"github.com/osmanyasin/grand-prix/internal/models"
	"github.com/osmanyasin/grand-prix/internal/simulator"
)

func CalculateFinalScore(level int, state *simulator.CarState, config *models.LevelConfig) (float64, float64, float64, float64) {
	baseScore := 1_000_000_000.0 / state.TotalTimeSeconds

	var fuelBonus float64 = 0.0
	var tyreBonus float64 = 0.0

	if level >= 2 {
		fuelUsed := state.TotalFuelUsedLitres
		softCap := config.Race.FuelSoftCapLimit

		ratio := 1.0 - (fuelUsed / softCap)
		fuelBonus = -1_000_000.0*(ratio*ratio) + 1_000_000.0
	}

	if level >= 4 {
		tyreBonus = (100_000.0 * state.TotalTyreDegradation) - (50_000.0 * float64(state.NumberOfBlowouts))
	}

	finalScore := baseScore + fuelBonus + tyreBonus
	return finalScore, baseScore, fuelBonus, tyreBonus
}
