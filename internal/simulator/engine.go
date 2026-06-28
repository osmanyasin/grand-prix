package simulator

import (
	"fmt"
	"math"

	"github.com/osmanyasin/grand-prix/internal/models"
	"github.com/osmanyasin/grand-prix/internal/physics"
)

func EvaluateRace(config *models.LevelConfig, strategy *models.Strategy) (*CarState, error) {
	state := NewCarState(config.Car.InitialFuel, strategy.InitialTyreID)

	trackMap := make(map[int]models.Segment)
	for _, seg := range config.Track.Segments {
		trackMap[seg.ID] = seg
	}

	for _, lapStrategy := range strategy.Laps {
		for _, action := range lapStrategy.Segments {
			trackSeg, ok := trackMap[action.ID]
			if !ok {
				return nil, fmt.Errorf("strategy references invalid segment ID %d", action.ID)
			}

			if state.IsLimpMode {
				processLimpMode(state, trackSeg, config)
			} else {
				if trackSeg.Type == "straight" {
					processStraight(state, action, trackSeg, config)
				} else if trackSeg.Type == "corner" {
					processCorner(state, trackSeg, config)
				}
			}

			checkFailureModes(state)
		}

		if lapStrategy.Pit.Enter {
			processPitStop(state, lapStrategy.Pit, config)
		}
	}

	return state, nil
}

func processStraight(state *CarState, action models.SegmentAction, seg models.Segment, config *models.LevelConfig) error {
	state.IsCrawlMode = false
	weather := getCurrentWeather(config, state.TotalTimeSeconds)
	tyreProps, err := getActiveTyreProperties(config, state.ActiveTyreID)
	if err != nil {
		return err
	}
	_, degradeRate := resolveTyreModifiers(tyreProps, weather.Condition)

	actualAccel := config.Car.AccelMPS2 * weather.AccelerationMulti
	actualDecel := config.Car.BrakeMPS2 * weather.DecelerationMulti

	initialSpeed := state.CurrentSpeedMPS
	targetSpeed := math.Min(action.TargetMPS, config.Car.MaxSpeedMPS)
	if targetSpeed < initialSpeed {
		targetSpeed = initialSpeed
	}

	brakeDist := action.BrakeStartMBeforeNext
	if brakeDist > seg.LengthM {
		brakeDist = seg.LengthM
	}
	accelCoastDist := seg.LengthM - brakeDist

	reqAccelDist := physics.CalculateDistance(initialSpeed, targetSpeed, actualAccel)
	var peakSpeed float64
	var accelDist float64

	if reqAccelDist <= accelCoastDist {
		accelDist = reqAccelDist
		peakSpeed = targetSpeed
	} else {
		accelDist = accelCoastDist
		peakSpeed = math.Sqrt(initialSpeed*initialSpeed + 2*actualAccel*accelDist)
	}

	accelTime := physics.TimeToAccelerate(initialSpeed, peakSpeed, actualAccel)
	accelFuel := physics.FuelUsed(initialSpeed, peakSpeed, accelDist)

	coastDist := accelCoastDist - accelDist
	coastTime := 0.0
	coastFuel := 0.0
	if coastDist > 0 {
		coastTime = coastDist / peakSpeed
		coastFuel = physics.FuelUsed(peakSpeed, peakSpeed, coastDist)
	}

	v2 := (peakSpeed * peakSpeed) - (2 * actualDecel * brakeDist)
	if v2 < 0 {
		v2 = 0
	}
	finalSpeed := math.Max(math.Sqrt(v2), config.Car.CrawlConstantMPS)
	brakeTime := (peakSpeed - finalSpeed) / actualDecel
	brakeFuel := physics.FuelUsed(peakSpeed, finalSpeed, brakeDist)

	straightWear := physics.DegradationStraight(degradeRate, seg.LengthM)
	brakeWear := physics.DegradationBraking(peakSpeed, finalSpeed, degradeRate)
	totalWear := straightWear + brakeWear

	state.TotalTimeSeconds += accelTime + coastTime + brakeTime
	state.CurrentFuelLitres -= accelFuel + coastFuel + brakeFuel
	state.TotalFuelUsedLitres += accelFuel + coastFuel + brakeFuel
	state.ActiveTyreDegrad += totalWear
	state.TotalTyreDegradation += totalWear
	state.CurrentSpeedMPS = finalSpeed

	return nil
}

func processCorner(state *CarState, seg models.Segment, config *models.LevelConfig) error {
	if state.IsCrawlMode {
		state.CurrentSpeedMPS = config.Car.CrawlConstantMPS
		state.TotalTimeSeconds += seg.LengthM / config.Car.CrawlConstantMPS
		return nil
	}

	weather := getCurrentWeather(config, state.TotalTimeSeconds)
	tyreProps, err := getActiveTyreProperties(config, state.ActiveTyreID)
	if err != nil {
		return err
	}
	frictionMulti, degradeRate := resolveTyreModifiers(tyreProps, weather.Condition)

	currentFriction := physics.CalculateTyreFriction(tyreProps.BaseFriction, state.ActiveTyreDegrad, frictionMulti)
	maxSafeSpeed := physics.CalculateMaxCornerSpeed(currentFriction, seg.RadiusM)

	if state.CurrentSpeedMPS > maxSafeSpeed {
		state.TotalTimeSeconds += config.Race.CornerCrashPenalty
		state.ActiveTyreDegrad += physics.CornerCrashPenaltyDegradation
		state.TotalTyreDegradation += physics.CornerCrashPenaltyDegradation
		state.IsCrawlMode = true
		state.CurrentSpeedMPS = config.Car.CrawlConstantMPS
	} else {
		state.TotalTimeSeconds += seg.LengthM / state.CurrentSpeedMPS
		cornerWear := physics.DegradationCorner(state.CurrentSpeedMPS, seg.RadiusM, degradeRate)
		state.ActiveTyreDegrad += cornerWear
		state.TotalTyreDegradation += cornerWear
	}

	cornerFuel := physics.FuelUsed(state.CurrentSpeedMPS, state.CurrentSpeedMPS, seg.LengthM)
	state.CurrentFuelLitres -= cornerFuel
	state.TotalFuelUsedLitres += cornerFuel

	return nil
}

func processLimpMode(state *CarState, seg models.Segment, config *models.LevelConfig) {
	state.CurrentSpeedMPS = config.Car.LimpConstantMPS
	state.TotalTimeSeconds += seg.LengthM / config.Car.LimpConstantMPS
}

func processPitStop(state *CarState, pit models.PitAction, config *models.LevelConfig) {
	pitTime := config.Race.BasePitStopTime

	if pit.TyreChangeSetID > 0 {
		state.ActiveTyreID = pit.TyreChangeSetID
		state.ActiveTyreDegrad = 0.0
		pitTime += config.Race.PitTyreSwapTime
	}

	if pit.FuelRefuelAmountL > 0 {
		fillAmount := math.Min(pit.FuelRefuelAmountL, config.Car.FuelTankCapacity-state.CurrentFuelLitres)
		state.CurrentFuelLitres += fillAmount
		pitTime += fillAmount / config.Race.PitRefuelRate
	}

	state.TotalTimeSeconds += pitTime

	state.CurrentSpeedMPS = config.Race.PitExitSpeedMPS
	state.IsLimpMode = false
	state.IsCrawlMode = false
}

func checkFailureModes(state *CarState) {
	if state.CurrentFuelLitres <= 0.0 {
		state.CurrentFuelLitres = 0.0
		state.IsLimpMode = true
	}

	if state.ActiveTyreDegrad >= 1.0 {
		state.NumberOfBlowouts++
		state.IsLimpMode = true
	}
}

func getCurrentWeather(config *models.LevelConfig, currentTime float64) models.WeatherCondition {
	if len(config.Weather.Conditions) == 0 {
		return models.WeatherCondition{Condition: "dry", AccelerationMulti: 1.0, DecelerationMulti: 1.0}
	}

	totalCycleTime := 0.0
	for _, w := range config.Weather.Conditions {
		totalCycleTime += w.DurationS
	}

	modTime := math.Mod(currentTime, totalCycleTime)

	accumulatedTime := 0.0
	for _, w := range config.Weather.Conditions {
		accumulatedTime += w.DurationS
		if modTime <= accumulatedTime {
			return w
		}
	}

	return config.Weather.Conditions[0]
}

func getActiveTyreProperties(config *models.LevelConfig, tyreID int) (*models.TyreProperty, error) {
	compoundName := ""
	for _, set := range config.AvailableSets {
		for _, id := range set.IDs {
			if id == tyreID {
				compoundName = set.Compound
				break
			}
		}
	}

	if compoundName == "" {
		return nil, fmt.Errorf("tyre ID %d not found in available sets", tyreID)
	}

	props, exists := config.Tyres.Properties[compoundName]
	if !exists {
		return nil, fmt.Errorf("properties for compound %s not found", compoundName)
	}

	return &props, nil
}

func resolveTyreModifiers(props *models.TyreProperty, weather string) (frictionMulti, degradeRate float64) {
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
