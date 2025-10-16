package models

import "time"

// WateringEvent represents a single watering occurrence in the irrigation system.
// Water is applied gradually over the specified Duration to simulate realistic
// irrigation behavior. Events can be triggered either manually or by the automated
// watering schedule.
type WateringEvent struct {
	SectionID string
	Amount    float64
	StartTime time.Time
	Duration  time.Duration
	IsManual  bool
}

// WateringSchedule defines the automated watering configuration for a garden section.
// The schedule monitors soil saturation at regular intervals and triggers watering
// events when saturation drops below the target threshold.
type WateringSchedule struct {
	SectionID        string
	TargetSaturation float64
	CheckInterval    int // in ticks
	WaterAmount      float64
	Enabled          bool
}
