package sensors

import "greenhouse-simulator/internal/models"

// SensorManager manages all sensors in the greenhouse and provides
// real-time readings grouped by plant sections.
type SensorManager interface {
	// AddSensor registers a new sensor in the system.
	AddSensor(sensor *models.Sensor) error
	// GetReading returns the current reading for a specific sensor.
	GetReading(sensorID string) (*models.SensorReading, error)
	// GetSectionReadings returns all sensor readings for a plant section.
	GetSectionReadings(sectionID string) ([]*models.SensorReading, error)
	// GetAverageSaturation calculates the average soil moisture for all sensors in a section.
	GetAverageSaturation(sectionID string) (float64, error)
}
