package solver

import (
	"math"

	"github.com/osmanyasin/grand-prix/internal/models"
	"github.com/osmanyasin/grand-prix/internal/physics"
)

type tyreSet struct {
	id       int
	compound string
}

func GenerateGreedyStrategy(config *models.LevelConfig) *models.Strategy {
	lapDistM := totalLapDistance(config.Track.Segments)
	fuelPerLap := lapDistM * config.Car.FuelConsumption

	var tyrePool []tyreSet
	for _, s := range config.AvailableSets {
		for _, id := range s.IDs {
			tyrePool = append(tyrePool, tyreSet{id: id, compound: s.Compound})
		}
	}
	usedSetIDs := make(map[int]bool)

	pitLaps := planFuelStops(config, fuelPerLap)

	openingWeather := "dry"
	if len(config.Weather.Conditions) > 0 {
		openingWeather = config.Weather.Conditions[0].Condition
	}
	startSet := pickBestCompound(config, openingWeather, tyrePool, usedSetIDs)
	usedSetIDs[startSet.id] = true

	strategy := &models.Strategy{
		InitialTyreID: startSet.id,
		Laps:          make([]models.Lap, config.Race.Laps),
	}

	activeCompound := startSet.compound
	projectedDegrade := 0.0
	expectedSpeed := 0.0
	currentFuel := config.Car.InitialFuel
	elapsedTime := 0.0

	const degradeThreshold = 0.90

	for l := 0; l < config.Race.Laps; l++ {
		lap := models.Lap{
			LapNumber: l + 1,
			Segments:  []models.SegmentAction{},
			Pit:       models.PitAction{Enter: false},
		}

		currentWeather := getWeatherAtTime(config, elapsedTime)
		props := config.Tyres.Properties[activeCompound]
		frictionMulti, degradeRate := resolveModifiers(&props, currentWeather)

		lapTime := 0.0

		for i, seg := range config.Track.Segments {
			action := models.SegmentAction{ID: seg.ID, Type: seg.Type}

			if seg.Type == "straight" {
				action.TargetMPS = config.Car.MaxSpeedMPS

				nextSegIndex := (i + 1) % len(config.Track.Segments)
				nextSeg := config.Track.Segments[nextSegIndex]

				if nextSeg.Type == "corner" {
					currentFriction := physics.CalculateTyreFriction(props.BaseFriction, projectedDegrade, frictionMulti)
					maxSafeSpeed := physics.CalculateMaxCornerSpeed(currentFriction, nextSeg.RadiusM, config.Car.CrawlConstantMPS)

					peakSpeed := math.Sqrt(math.Max(0, expectedSpeed*expectedSpeed+2*config.Car.AccelMPS2*seg.LengthM))
					peakSpeed = math.Min(peakSpeed, config.Car.MaxSpeedMPS)

					if peakSpeed > maxSafeSpeed {
						brakeDist := (peakSpeed*peakSpeed - maxSafeSpeed*maxSafeSpeed) / (2 * config.Car.BrakeMPS2)
						action.BrakeStartMBeforeNext = math.Min(math.Max(brakeDist, 0), seg.LengthM)
						expectedSpeed = maxSafeSpeed
					} else {
						action.BrakeStartMBeforeNext = 0
						expectedSpeed = peakSpeed
					}
				} else {
					expectedSpeed = config.Car.MaxSpeedMPS
				}

				projectedDegrade += physics.DegradationStraight(degradeRate, seg.LengthM)
				lapTime += seg.LengthM / math.Max(expectedSpeed, 1)
			} else if seg.Type == "corner" {
				action.TargetMPS = expectedSpeed
				projectedDegrade += physics.DegradationCorner(expectedSpeed, seg.RadiusM, degradeRate)
				lapTime += seg.LengthM / math.Max(expectedSpeed, 1)
			}

			lap.Segments = append(lap.Segments, action)
		}

		elapsedTime += lapTime
		currentFuel -= fuelPerLap

		needsTyreChange := projectedDegrade >= degradeThreshold
		needsFuelStop := isFuelPitLap(pitLaps, l+1)

		if needsTyreChange || needsFuelStop {
			lap.Pit.Enter = true
			nextWeather := getWeatherAtTime(config, elapsedTime)
			nextSet := pickBestCompound(config, nextWeather, tyrePool, usedSetIDs)
			if nextSet.id != 0 {
				usedSetIDs[nextSet.id] = true
				lap.Pit.TyreChangeSetID = nextSet.id
				activeCompound = nextSet.compound
				projectedDegrade = 0.0
			}

			if needsFuelStop || currentFuel < fuelPerLap*3 {
				lapsRemaining := config.Race.Laps - (l + 1)
				nextStop := nextFuelStop(pitLaps, l+1)
				lapsToNextStop := lapsRemaining
				if nextStop > 0 {
					lapsToNextStop = nextStop - (l + 1)
				}
				fuelNeeded := float64(lapsToNextStop)*fuelPerLap + fuelPerLap*2
				refuelAmount := fuelNeeded - currentFuel
				refuelAmount = math.Max(0, refuelAmount)
				refuelAmount = math.Min(refuelAmount, config.Car.FuelTankCapacity-currentFuel)
				if refuelAmount > 0 {
					lap.Pit.FuelRefuelAmountL = refuelAmount
					currentFuel += refuelAmount
				}
			}
		}

		strategy.Laps[l] = lap
	}

	return strategy
}

func totalLapDistance(segments []models.Segment) float64 {
	total := 0.0
	for _, s := range segments {
		total += s.LengthM
	}
	return total
}

func planFuelStops(config *models.LevelConfig, fuelPerLap float64) []int {
	var stops []int
	fuel := config.Car.InitialFuel
	lap := 0
	totalLaps := config.Race.Laps

	for lap < totalLaps {
		lapsLeft := int(fuel / fuelPerLap)
		if lap+lapsLeft >= totalLaps {
			break
		}
		stopLap := lap + intMax(1, lapsLeft-2)
		stops = append(stops, stopLap)
		fuelAtStop := fuel - float64(stopLap-lap)*fuelPerLap
		fuel = fuelAtStop + (config.Car.FuelTankCapacity - fuelAtStop)
		lap = stopLap
	}
	return stops
}

func isFuelPitLap(pitLaps []int, lap int) bool {
	for _, p := range pitLaps {
		if p == lap {
			return true
		}
	}
	return false
}

func nextFuelStop(pitLaps []int, afterLap int) int {
	for _, p := range pitLaps {
		if p > afterLap {
			return p
		}
	}
	return 0
}

func pickBestCompound(config *models.LevelConfig, weather string, pool []tyreSet, used map[int]bool) tyreSet {
	best := tyreSet{}
	bestScore := -math.MaxFloat64

	for _, ts := range pool {
		if used[ts.id] {
			continue
		}
		props, ok := config.Tyres.Properties[ts.compound]
		if !ok {
			continue
		}
		frictionMulti, degradeRate := resolveModifiers(&props, weather)
		score := (props.BaseFriction * frictionMulti) / math.Max(degradeRate, 0.001)
		if score > bestScore {
			bestScore = score
			best = ts
		}
	}

	return best
}

func getWeatherAtTime(config *models.LevelConfig, t float64) string {
	if len(config.Weather.Conditions) == 0 {
		return "dry"
	}
	elapsed := 0.0
	for _, w := range config.Weather.Conditions {
		elapsed += w.DurationS
		if t < elapsed {
			return w.Condition
		}
	}
	return config.Weather.Conditions[len(config.Weather.Conditions)-1].Condition
}

func resolveModifiers(props *models.TyreProperty, weather string) (float64, float64) {
	switch weather {
	case "cold":
		return props.ColdFrictionMulti, props.ColdDegradation
	case "light_rain":
		return props.LightRainFrictionMulti, props.LightRainDegradation
	case "heavy_rain":
		return props.HeavyRainFrictionMulti, props.HeavyRainDegradation
	default:
		return props.DryFrictionMulti, props.DryDegradation
	}
}

func intMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}
