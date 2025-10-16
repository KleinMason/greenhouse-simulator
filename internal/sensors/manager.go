package sensors

import (
	"errors"
	"greenhouse-simulator/internal/models"
	"sync"
	"time"
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
	plantData        PlantDataSource
	mu               sync.RWMutex
}

// NewSensorManager creates and returns a new SensorManager instance.
// The returned manager is initialized with empty maps for tracking sensors
// by section and by ID, and is safe for concurrent use.
func NewSensorManager(plantData PlantDataSource) SensorManager {
	return &sensorManager{
		sensorsBySection: make(map[string][]*models.Sensor),
		sensorsByID:      map[string]*models.Sensor{},
		plantData:        plantData,
	}
}

// AddSensor registers a new sensor in the system and associates it with a plant section.
// The sensor must have a valid ID and SectionID. Returns an error if:
// - sensor is nil
// - sensor ID is empty
// - sensor section ID is empty
// - a sensor with the same ID already exists
//
// This method is safe for concurrent use.
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

// GetReading retrieves the current sensor reading for the specified sensor ID.
// It calculates the reading value by averaging the soil saturation of all plants
// in the sensor's associated section.
//
// Parameters:
//   - sensorID: The unique identifier of the sensor to get a reading from
//
// Returns:
//   - *models.SensorReading: A reading containing the sensor ID, current timestamp,
//     and the calculated average soil saturation value
//   - error: An error if the sensor ID is not found or if there are no plants
//     in the sensor's section
//
// The method is safe for concurrent use as it acquires a read lock during execution.
// The returned reading's Value field represents the average soil saturation percentage
// across all plants in the sensor's section.
func (s *sensorManager) GetReading(sensorID string) (*models.SensorReading, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sensor := s.sensorsByID[sensorID]
	if sensor == nil {
		return nil, errors.New("no sensor found for the provided ID: " + sensorID)
	}

	plants := s.plantData.GetPlantsBySectionID(sensor.SectionID)
	if len(plants) == 0 {
		return nil, errors.New("no plants in section: " + sensor.SectionID)
	}
	total := 0.0
	for _, plant := range plants {
		total += plant.SoilSaturation
	}
	average := total / float64(len(plants))

	return &models.SensorReading{
		SensorID:  sensor.ID,
		Timestamp: time.Now(),
		Value:     average,
	}, nil
}

func (s *sensorManager) GetSectionReadings(sectionID string) ([]*models.SensorReading, error) {
	return nil, errors.New("not implemented")
}

func (s *sensorManager) GetAverageSaturation(sectionID string) (float64, error) {
	return 0, errors.New("not implemented")
}
