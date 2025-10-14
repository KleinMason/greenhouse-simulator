package models

import "time"

// SensorType represents the different types of sensors available in the system.
type SensorType string

const (
	// SoilMoisture sensors measure soil water content (0.0 to 1.0).
	SoilMoisture SensorType = "soil_moisture"
	// Temperature sensors measure ambient temperature in Celsius.
	Temperature SensorType = "temperature"
	// Light sensor measures light/luminosity intensity.
	Light SensorType = "light"
	// Humidity sensors measure relative humidity (0.0 to 1.0).
	Humidity SensorType = "humidity"
)

// Sensor represents a physical sensor device in the greenhouse.
// Each sensor monitors a specific section and measures one environmental factor.
type Sensor struct {
	ID        string
	Type      SensorType
	SectionID string
}

// SensorReading represents a single measurement taken by a sensor.
type SensorReading struct {
	SensorID  string
	Timestamp time.Time
	Value     float64
}
