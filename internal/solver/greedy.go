package solver

import (
	"math"

	"github.com/osmanyasin/grand-prix/internal/models"
	"github.com/osmanyasin/grand-prix/internal/physics"
)

// GenerateGreedyStrategy builds a segment-by-segment plan using state-aware calculations.
func GenerateGreedyStrategy(config *models.LevelConfig) *models.Strategy {
	// 1. Pick a starting tyre
	initialTyreID := 1
	compoundName := ""
	if len(config.AvailableSets) > 0 {
		initialTyreID = config.AvailableSets[0].IDs[0]
		compoundName = config.AvailableSets[0].Compound
	}

	strategy := &models.Strategy{
		InitialTyreID: initialTyreID,
		Laps:          make([]models.Lap, config.Race.Laps),
	}

	projectedDegradation := 0.0
	expectedSpeed := 0.0 // Track the car's state across segments
	tyresUsedIndex := 0

	for l := 0; l < config.Race.Laps; l++ {
		lap := models.Lap{
			LapNumber: l + 1,
			Segments:  []models.SegmentAction{},
			Pit:       models.PitAction{Enter: false},
		}

		for i, seg := range config.Track.Segments {
			action := models.SegmentAction{ID: seg.ID, Type: seg.Type}
			props := config.Tyres.Properties[compoundName]

			if seg.Type == "straight" {
				action.TargetMPS = config.Car.MaxSpeedMPS

				// Look ahead to check if the next segment is a corner
				nextSegIndex := (i + 1) % len(config.Track.Segments)
				nextSeg := config.Track.Segments[nextSegIndex]

				if nextSeg.Type == "corner" {
					// 1. Calculate Safe Corner Entry Speed
					currentFriction := physics.CalculateTyreFriction(props.BaseFriction, projectedDegradation, props.DryFrictionMulti)
					maxSafeSpeed := physics.CalculateMaxCornerSpeed(currentFriction, nextSeg.RadiusM) * 0.85

					// 2. Calculate Peak Speed on this straight
					// v_f = sqrt(v_i^2 + 2ad)
					peakSpeed := math.Sqrt(math.Max(0, math.Pow(expectedSpeed, 2)+2*config.Car.AccelMPS2*seg.LengthM))
					peakSpeed = math.Min(peakSpeed, config.Car.MaxSpeedMPS)

					// 3. Determine if braking is required
					if peakSpeed > maxSafeSpeed {
						// d = (v_final^2 - v_initial^2) / (2 * a)
						brakeDist := (math.Pow(maxSafeSpeed, 2) - math.Pow(peakSpeed, 2)) / (2 * -config.Car.BrakeMPS2)
						action.BrakeStartMBeforeNext = math.Min(math.Max(brakeDist, 0), seg.LengthM)
						expectedSpeed = maxSafeSpeed
					} else {
						action.BrakeStartMBeforeNext = 0
						expectedSpeed = peakSpeed
					}
				} else {
					expectedSpeed = config.Car.MaxSpeedMPS
				}

				projectedDegradation += physics.DegradationStraight(props.DryDegradation, seg.LengthM)
			} else if seg.Type == "corner" {
				// We assume we exit the corner at the same speed we entered
				action.TargetMPS = expectedSpeed
				projectedDegradation += physics.DegradationCorner(expectedSpeed, seg.RadiusM, props.DryDegradation)
			}

			lap.Segments = append(lap.Segments, action)
		}

		// Pit Stop Logic
		if projectedDegradation > 0.70 {
			lap.Pit.Enter = true
			tyresUsedIndex++
			if tyresUsedIndex < len(config.AvailableSets[0].IDs) {
				lap.Pit.TyreChangeSetID = config.AvailableSets[0].IDs[tyresUsedIndex]
			}
			lap.Pit.FuelRefuelAmountL = 50.0
			projectedDegradation = 0.0
		}

		strategy.Laps[l] = lap
	}

	return strategy
}
