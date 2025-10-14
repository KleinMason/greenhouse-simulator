package models

import (
	"math"
	"testing"
)

const floatTolerance = 0.0001

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < floatTolerance
}

func TestHealthDegradation(t *testing.T) {
	tests := []struct {
		name               string
		initialHealth      float64
		soilSaturation     float64
		expectedHealthLess float64
	}{
		{
			name:               "health degrades when soil too dry",
			initialHealth:      0.8,
			soilSaturation:     0.1, // below MinSaturation
			expectedHealthLess: 0.8,
		}, {
			name:               "health degrades when soil too wet",
			initialHealth:      0.8,
			soilSaturation:     0.9, // above MaxSaturation
			expectedHealthLess: 0.8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plant := &Plant{
				Type: PlantType{
					MinSaturation:         0.3,
					MaxSaturation:         0.7,
					HealthDegradationRate: 0.05,
				},
				Health:         tt.initialHealth,
				SoilSaturation: tt.soilSaturation,
				Alive:          true,
			}

			plant.OnTick()

			if plant.Health >= tt.expectedHealthLess {
				t.Errorf("expected health to degrade below %.2f, got %.2f",
					tt.expectedHealthLess, plant.Health)
			}
		})
	}
}

func TestHealthEnhancement(t *testing.T) {
	tests := []struct {
		name                  string
		initialHealth         float64
		soilSaturation        float64
		expectedHealthGreater float64
	}{
		{
			name:                  "health enhances when soil above MinSaturation and below MaxSaturation",
			initialHealth:         0.8,
			soilSaturation:        0.6, // above MinSaturation and below MaxSaturation
			expectedHealthGreater: 0.8,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plant := &Plant{
				Type: PlantType{
					MinSaturation:         0.3,
					MaxSaturation:         0.7,
					HealthEnhancementRate: 0.05,
				},
				Health:         tt.initialHealth,
				SoilSaturation: tt.soilSaturation,
				Alive:          true,
			}

			plant.OnTick()

			if plant.Health <= tt.expectedHealthGreater {
				t.Errorf("expected health to enhance above %.2f, got %.2f",
					tt.expectedHealthGreater, plant.Health)
			}
		})
	}
}

func TestGrowthStages_NoGrowthBelowThreshold(t *testing.T) {
	tests := []struct {
		name          string
		initialHealth float64
		initialGrowth float64
	}{
		{"health at 0.2", 0.2, 0.5},
		{"health at 0.29", 0.29, 0.3},
		{"health at 0.0", 0.0, 0.4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plant := &Plant{
				Type: PlantType{
					BaseGrowthRate: 0.1,
					MinSaturation:  0.3,
					MaxSaturation:  0.7,
				},
				Health:      tt.initialHealth,
				GrowthStage: tt.initialGrowth,
				Alive:       true,
			}

			plant.OnTick()

			if plant.GrowthStage != tt.initialGrowth {
				t.Errorf("expected no growth (%.2f), got %.2f",
					tt.initialGrowth, plant.GrowthStage)
			}
		})
	}
}

func TestGrowthStages_ReducedGrowthAtLowHealth(t *testing.T) {
	tests := []struct {
		name           string
		initialHealth  float64
		initialGrowth  float64
		baseGrowthRate float64
	}{
		{"health at 0.3", 0.3, 0.5, 0.1},
		{"health at 0.4", 0.4, 0.3, 0.15},
		{"health at 0.49", 0.49, 0.4, 0.2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plant := &Plant{
				Type: PlantType{
					BaseGrowthRate: tt.baseGrowthRate,
					MinSaturation:  0.3,
					MaxSaturation:  0.7,
				},
				Health:      tt.initialHealth,
				GrowthStage: tt.initialGrowth,
				Alive:       true,
			}

			plant.OnTick()

			actualIncrease := plant.GrowthStage - tt.initialGrowth

			if actualIncrease >= tt.baseGrowthRate {
				t.Errorf("expected slow growth (< %.2f), got %.2f",
					tt.baseGrowthRate, actualIncrease)
			}
		})
	}
}

func TestGrowthStages_NormalGrowth(t *testing.T) {
	tests := []struct {
		name              string
		initialHealth     float64
		initialGrowth     float64
		baseGrowthRate    float64
		soilSaturation    float64
		optimalSaturation float64
	}{
		{"health at 0.6", 0.6, 0.5, 0.05, 0.4, 0.7},
		{"health at 0.85", 0.85, 0.5, 0.05, 0.4, 0.7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plant := &Plant{
				Type: PlantType{
					BaseGrowthRate:    tt.baseGrowthRate,
					OptimalSaturation: tt.optimalSaturation,
					MinSaturation:     0.2,
					MaxSaturation:     0.8,
				},
				Health:         tt.initialHealth,
				GrowthStage:    tt.initialGrowth,
				SoilSaturation: tt.soilSaturation,
				Alive:          true,
			}

			plant.OnTick()

			actualIncrease := plant.GrowthStage - tt.initialGrowth

			if !almostEqual(actualIncrease, tt.baseGrowthRate) {
				t.Errorf("expected normal growth (~ %.2f), got %.2f",
					tt.baseGrowthRate, actualIncrease)
			}
		})
	}
}

func TestGrowthStages_BonusGrowth(t *testing.T) {
	tests := []struct {
		name              string
		initialHealth     float64
		initialGrowth     float64
		baseGrowthRate    float64
		soilSaturation    float64
		optimalSaturation float64
	}{
		{"health at 0.6", 0.6, 0.5, 0.05, 0.7, 0.7},
		{"health at 0.85", 0.85, 0.5, 0.05, 0.67, 0.7},
		{"health at 0.75", 0.75, 0.5, 0.05, 0.76, 0.7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plant := &Plant{
				Type: PlantType{
					BaseGrowthRate:    tt.baseGrowthRate,
					OptimalSaturation: tt.optimalSaturation,
					MinSaturation:     0.2,
					MaxSaturation:     0.8,
				},
				Health:         tt.initialHealth,
				GrowthStage:    tt.initialGrowth,
				SoilSaturation: tt.soilSaturation,
				Alive:          true,
			}

			plant.OnTick()

			actualIncrease := plant.GrowthStage - tt.initialGrowth

			if actualIncrease <= tt.baseGrowthRate {
				t.Errorf("expected bonus growth (> %.2f), got %.2f",
					tt.baseGrowthRate, actualIncrease)
			}
		})
	}
}

func TestGrowthStages_CapAt1(t *testing.T) {
	tests := []struct {
		name              string
		initialHealth     float64
		initialGrowth     float64
		soilSaturation    float64
		optimalSaturation float64
	}{
		{"health at 0.6", 0.6, 0.80, 0.6, 0.7},
		{"health at 0.85", 0.85, 0.995, 0.65, 0.7},
		{"health at 0.75", 0.75, 0.82, 0.99, 0.7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plant := &Plant{
				Type: PlantType{
					BaseGrowthRate:    0.2,
					OptimalSaturation: tt.optimalSaturation,
					MinSaturation:     0.2,
					MaxSaturation:     0.8,
				},
				Health:         tt.initialHealth,
				GrowthStage:    tt.initialGrowth,
				SoilSaturation: tt.soilSaturation,
				Alive:          true,
			}

			plant.OnTick()

			if plant.GrowthStage > 1.0 {
				t.Errorf("expected growth to cap at 1.0, got %.2f",
					plant.GrowthStage)
			}

			if plant.GrowthStage != 1.0 {
				t.Errorf("expected growth to be exactly 1.0, got %.2f",
					plant.GrowthStage)
			}
		})
	}
}

func TestDeath_NoActionIfDead(t *testing.T) {
	plant := &Plant{
		Alive: false,
	}
	plant.OnTick()
	if plant.Alive {
		t.Errorf("expected plant to be dead, got alive")
	}
	if plant.Health != 0 {
		t.Errorf("expected health to be 0, got %.2f", plant.Health)
	}
	if plant.GrowthStage != 0 {
		t.Errorf("expected growth stage to be 0, got %.2f", plant.GrowthStage)
	}
	if plant.SoilSaturation != 0 {
		t.Errorf("expected soil saturation to be 0, got %.2f", plant.SoilSaturation)
	}
}

func TestDeath_HealthDegradesToZero(t *testing.T) {
	plant := &Plant{
		Health:         0.01,
		SoilSaturation: 0.0,
		Type: PlantType{
			OptimalSaturation:     0.6,
			MinSaturation:         0.3,
			MaxSaturation:         0.7,
			HealthDegradationRate: 0.08,
		},
		Alive: true,
	}
	plant.OnTick()
	if plant.Health != 0 {
		t.Errorf("expected health to be 0, got %.2f", plant.Health)
	}
	if plant.Alive {
		t.Errorf("expected plant to be dead (Alive=false), got Alive=%v", plant.Alive)
	}
}

func TestHealthClamps_ClampTo1(t *testing.T) {
	tests := []struct {
		name           string
		initialHealth  float64
		soilSaturation float64
	}{
		{
			name:           "health clamps to 1",
			initialHealth:  0.95,
			soilSaturation: 0.6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plant := &Plant{
				Type: PlantType{
					MinSaturation:         0.3,
					MaxSaturation:         0.7,
					HealthDegradationRate: 0.1,
					HealthEnhancementRate: 0.1,
				},
				Health:         tt.initialHealth,
				SoilSaturation: tt.soilSaturation,
				Alive:          true,
			}

			plant.OnTick()

			if plant.Health > 1.0 {
				t.Errorf("expected health to clamp to 1.0, got %.2f",
					plant.Health)
			}
			if plant.Health != 1.0 {
				t.Errorf("expected health to be exactly 1.0, got %.2f",
					plant.Health)
			}
		})
	}
}

func TestHealthClamps_ClampTo0(t *testing.T) {
	tests := []struct {
		name           string
		initialHealth  float64
		soilSaturation float64
	}{
		{
			name:           "health clamps to 0",
			initialHealth:  0.05,
			soilSaturation: 0.1, // below MinSaturation

		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plant := &Plant{
				Type: PlantType{
					MinSaturation:         0.3,
					MaxSaturation:         0.7,
					HealthDegradationRate: 0.1,
					HealthEnhancementRate: 0.1,
				},
				Health:         tt.initialHealth,
				SoilSaturation: tt.soilSaturation,
				Alive:          true,
			}

			plant.OnTick()

			if plant.Health < 0.0 {
				t.Errorf("expected health to clamp to 0.0, got %.2f",
					plant.Health)
			}
			if plant.Health != 0.0 {
				t.Errorf("expected health to be exactly 0.0, got %.2f",
					plant.Health)
			}
		})
	}
}

func TestSaturationClamps_ClampTo0(t *testing.T) {

	plant := &Plant{
		SoilSaturation: 0.02,
		Type: PlantType{
			SaturationDepletion: 0.05,
		},
		Health: 0.5,
		Alive:  true,
	}

	plant.OnTick()

	if plant.SoilSaturation < 0.0 {
		t.Errorf("expected soil saturation to clamp to 0.0, got %.2f",
			plant.SoilSaturation)
	}
	if plant.SoilSaturation != 0.0 {
		t.Errorf("expected soil saturation to be exactly 0.0, got %.2f",
			plant.SoilSaturation)
	}
}
