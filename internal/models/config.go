package models

type LevelConfig struct {
	Car           Car     `json:"car"`
	Race          Race    `json:"race"`
	Track         Track   `json:"track"`
	Tyres         Tyres   `json:"tyres"`
	AvailableSets []Set   `json:"available_sets"`
	Weather       Weather `json:"weather"`
}

type Car struct {
	MaxSpeedMPS      float64 `json:"max_speed_m/s"`
	AccelMPS2        float64 `json:"accel_m/se2"`
	BrakeMPS2        float64 `json:"brake_m/se2"`
	LimpConstantMPS  float64 `json:"limp_constant_m/s"`
	CrawlConstantMPS float64 `json:"crawl_constant_m/s"`
	FuelTankCapacity float64 `json:"fuel_tank_capacity_l"`
	InitialFuel      float64 `json:"initial_fuel_l"`
	FuelConsumption  float64 `json:"fuel_consumption_l/m"`
}

type Race struct {
	Name                  string  `json:"name"`
	Laps                  int     `json:"laps"`
	BasePitStopTime       float64 `json:"base_pit_stop_time_s"`
	PitTyreSwapTime       float64 `json:"pit_tyre_swap_time_s"`
	PitRefuelRate         float64 `json:"pit_refuel_rate_l/s"`
	CornerCrashPenalty    float64 `json:"corner_crash_penalty_s"`
	PitExitSpeedMPS       float64 `json:"pit_exit_speed_m/s"`
	FuelSoftCapLimit      float64 `json:"fuel_soft_cap_limit_l"`
	StartingWeatherCondID int     `json:"starting_weather_condition_id"`
}

type Track struct {
	Name     string    `json:"name"`
	Segments []Segment `json:"segments"`
}

type Segment struct {
	ID      int     `json:"id"`
	Type    string  `json:"type"`
	LengthM float64 `json:"length_m"`
	RadiusM float64 `json:"radius_m,omitempty"`
}

type Tyres struct {
	Properties map[string]TyreProperty `json:"properties"`
}

type TyreProperty struct {
	LifeSpan     float64 `json:"life_span"`
	BaseFriction float64 `json:"base_friction"`

	// Friction Multipliers
	DryFrictionMulti       float64 `json:"dry_friction_multiplier"`
	ColdFrictionMulti      float64 `json:"cold_friction_multiplier"`
	LightRainFrictionMulti float64 `json:"light_rain_friction_multiplier"`
	HeavyRainFrictionMulti float64 `json:"heavy_rain_friction_multiplier"`

	// Degradation Rates
	DryDegradation       float64 `json:"dry_degradation"`
	ColdDegradation      float64 `json:"cold_degradation"`
	LightRainDegradation float64 `json:"light_rain_degradation"`
	HeavyRainDegradation float64 `json:"heavy_rain_degradation"`
}

type Set struct {
	IDs      []int  `json:"ids"`
	Compound string `json:"compound"`
}

type Weather struct {
	Conditions []WeatherCondition `json:"conditions"`
}

type WeatherCondition struct {
	ID                int     `json:"id"`
	Condition         string  `json:"condition"`
	DurationS         float64 `json:"duration_s"`
	AccelerationMulti float64 `json:"acceleration_multiplier"`
	DecelerationMulti float64 `json:"deceleration_multiplier"`
}
