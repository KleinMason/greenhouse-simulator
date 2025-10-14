package sensors

import (
	"errors"
	"greenhouse-simulator/internal/models"
	"sync"
)

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

type sensorManager struct {
	sensorsBySection map[string][]*models.Sensor
	sensorsByID      map[string]*models.Sensor
	mu               sync.RWMutex
}

func NewSensorManager() SensorManager {
	return &sensorManager{
		sensorsBySection: make(map[string][]*models.Sensor),
		sensorsByID:      map[string]*models.Sensor{},
	}
}

func (s *sensorManager) AddSensor(sensor *models.Sensor) error {
	if sensor == nil {
		return errors.New("sensor cannot be nil")
	}
	if sensor.ID == "" {
		return errors.New("sensor ID cannot be empty")
	}
	if sensor.SectionID == "" {
		return errors.New("sensor section ID cannot be empty")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if exists := s.sensorsByID[sensor.ID]; exists != nil {
		return errors.New("sensor with ID already exists: " + sensor.ID)
	}

	s.sensorsByID[sensor.ID] = sensor
	s.sensorsBySection[sensor.SectionID] = append(s.sensorsBySection[sensor.SectionID], sensor)

	return nil
}

func (s *sensorManager) GetReading(sensorID string) (*models.SensorReading, error) {
	return nil, errors.New("not implemented")
}

func (s *sensorManager) GetSectionReadings(sectionID string) ([]*models.SensorReading, error) {
	return nil, errors.New("not implemented")
}

func (s *sensorManager) GetAverageSaturation(sectionID string) (float64, error) {
	return 0, errors.New("not implemented")
}
