package sensors

import (
	"greenhouse-simulator/internal/models"
	"testing"
	"time"
)

// mockPlantDataSource is a test mock implementation of PlantDataSource
type mockPlantDataSource struct {
	plantsBySectionID map[string][]*models.Plant
	allPlants         []*models.Plant
}

func (m *mockPlantDataSource) GetPlantsBySectionID(sectionID string) []*models.Plant {
	return m.plantsBySectionID[sectionID]
}

func (m *mockPlantDataSource) GetAllPlants() []*models.Plant {
	return m.allPlants
}

// Helper function to create a test plant
func createTestPlant(id, sectionID string, soilSaturation float64) *models.Plant {
	plantType := models.PlantType{
		Name:                  "TestPlant",
		OptimalSaturation:     0.7,
		MinSaturation:         0.3,
		MaxSaturation:         0.9,
		BaseGrowthRate:        0.01,
		SaturationDepletion:   0.02,
		HealthDegradationRate: 0.05,
		HealthEnhancementRate: 0.03,
	}
	plant, _ := models.NewPlant(id, plantType, sectionID, soilSaturation)
	return plant
}

func TestNewSensorManager(t *testing.T) {
	mockData := &mockPlantDataSource{
		plantsBySectionID: make(map[string][]*models.Plant),
	}

	manager := NewSensorManager(mockData)

	if manager == nil {
		t.Fatal("NewSensorManager returned nil")
	}
}

func TestAddSensor(t *testing.T) {
	tests := []struct {
		name        string
		sensor      *models.Sensor
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid sensor",
			sensor: &models.Sensor{
				ID:        "sensor-1",
				Type:      models.SoilMoisture,
				SectionID: "section-A",
			},
			expectError: false,
		},
		{
			name:        "nil sensor",
			sensor:      nil,
			expectError: true,
			errorMsg:    "sensor cannot be nil",
		},
		{
			name: "empty sensor ID",
			sensor: &models.Sensor{
				ID:        "",
				Type:      models.SoilMoisture,
				SectionID: "section-A",
			},
			expectError: true,
			errorMsg:    "sensor ID cannot be empty",
		},
		{
			name: "empty section ID",
			sensor: &models.Sensor{
				ID:        "sensor-1",
				Type:      models.SoilMoisture,
				SectionID: "",
			},
			expectError: true,
			errorMsg:    "sensor section ID cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockData := &mockPlantDataSource{
				plantsBySectionID: make(map[string][]*models.Plant),
			}
			manager := NewSensorManager(mockData)

			err := manager.AddSensor(tt.sensor)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got nil")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("expected error message '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestAddSensor_DuplicateID(t *testing.T) {
	mockData := &mockPlantDataSource{
		plantsBySectionID: make(map[string][]*models.Plant),
	}
	manager := NewSensorManager(mockData)

	sensor1 := &models.Sensor{
		ID:        "sensor-1",
		Type:      models.SoilMoisture,
		SectionID: "section-A",
	}

	// Add first sensor
	err := manager.AddSensor(sensor1)
	if err != nil {
		t.Fatalf("failed to add first sensor: %v", err)
	}

	// Try to add sensor with same ID
	sensor2 := &models.Sensor{
		ID:        "sensor-1",
		Type:      models.Temperature,
		SectionID: "section-B",
	}

	err = manager.AddSensor(sensor2)
	if err == nil {
		t.Error("expected error when adding duplicate sensor ID, got nil")
	}
	// TODO: You can add more specific error message validation here
}

func TestGetReading(t *testing.T) {
	// Create test plants with different saturation levels
	plant1 := createTestPlant("plant-1", "section-A", 0.6)
	plant2 := createTestPlant("plant-2", "section-A", 0.8)
	plant3 := createTestPlant("plant-3", "section-A", 0.7)

	mockData := &mockPlantDataSource{
		plantsBySectionID: map[string][]*models.Plant{
			"section-A": {plant1, plant2, plant3},
		},
	}

	manager := NewSensorManager(mockData)

	sensor := &models.Sensor{
		ID:        "sensor-1",
		Type:      models.SoilMoisture,
		SectionID: "section-A",
	}

	err := manager.AddSensor(sensor)
	if err != nil {
		t.Fatalf("failed to add sensor: %v", err)
	}

	reading, err := manager.GetReading("sensor-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if reading == nil {
		t.Fatal("reading is nil")
	}

	// Check sensor ID
	if reading.SensorID != "sensor-1" {
		t.Errorf("expected SensorID 'sensor-1', got '%s'", reading.SensorID)
	}

	// Check timestamp is recent (within last second)
	timeDiff := time.Since(reading.Timestamp)
	if timeDiff > time.Second {
		t.Errorf("timestamp is too old: %v", timeDiff)
	}

	// Check average calculation: (0.6 + 0.8 + 0.7) / 3 = 0.7
	// Note: We use an epsilon for float comparison to handle floating-point precision
	expectedAverage := 0.7
	epsilon := 0.0001
	if reading.Value < expectedAverage-epsilon || reading.Value > expectedAverage+epsilon {
		t.Errorf("expected average %f, got %f", expectedAverage, reading.Value)
	}
}

func TestGetReading_SensorNotFound(t *testing.T) {
	mockData := &mockPlantDataSource{
		plantsBySectionID: make(map[string][]*models.Plant),
	}
	manager := NewSensorManager(mockData)

	_, err := manager.GetReading("nonexistent-sensor")
	if err == nil {
		t.Error("expected error for nonexistent sensor, got nil")
	}
	// TODO: Add more specific error message check
}

func TestGetReading_NoPlants(t *testing.T) {
	// Create a section with no plants
	mockData := &mockPlantDataSource{
		plantsBySectionID: map[string][]*models.Plant{
			"empty-section": {},
		},
	}

	manager := NewSensorManager(mockData)

	sensor := &models.Sensor{
		ID:        "sensor-1",
		Type:      models.SoilMoisture,
		SectionID: "empty-section",
	}

	err := manager.AddSensor(sensor)
	if err != nil {
		t.Fatalf("failed to add sensor: %v", err)
	}

	_, err = manager.GetReading("sensor-1")
	if err == nil {
		t.Error("expected error when section has no plants, got nil")
	}
}

// TODO: Add tests for GetSectionReadings once implemented
// TODO: Add tests for GetAverageSaturation once implemented
// TODO: Consider adding concurrent access tests to verify thread-safety
