package models

type Strategy struct {
	InitialTyreID int   `json:"initial_tyre_id"`
	Laps          []Lap `json:"laps"`
}

type Lap struct {
	LapNumber int             `json:"lap"`
	Segments  []SegmentAction `json:"segments"`
	Pit       PitAction       `json:"pit"`
}

type SegmentAction struct {
	ID                    int     `json:"id"`
	Type                  string  `json:"type"`
	TargetMPS             float64 `json:"target_m/s,omitempty"`
	BrakeStartMBeforeNext float64 `json:"brake_start_m_before_next,omitempty"`
}

type PitAction struct {
	Enter             bool    `json:"enter"`
	TyreChangeSetID   int     `json:"tyre_change_set_id,omitempty"`
	FuelRefuelAmountL float64 `json:"fuel_refuel_amount_l,omitempty"`
}
