package main

import (
	"greenhouse-simulator/internal/engine"
	"greenhouse-simulator/internal/models"
	"greenhouse-simulator/internal/sensors"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	sim := engine.NewSimulator(4 * time.Second)
	sensorMgr := sensors.NewSensorManager(sim)

	sensor := &models.Sensor{
		ID:        "sensor-1",
		Type:      models.SoilMoisture,
		SectionID: "section-B",
	}
	sensorMgr.AddSensor(sensor)

	testPlants := getTestPlants()

	for i := range testPlants {
		err := sim.AddPlant(testPlants[i])
		if err != nil {
			slog.Warn(err.Error())
		}
	}

	go sim.Start()

	reading, err := sensorMgr.GetReading("sensor-1")
	if err != nil {
		slog.Error("failed to get reading for sensor", "error", err)
	}
	slog.Info("sensor reading", "SensorID", reading.SensorID, "Timestamp", reading.Timestamp, "Value", reading.Value)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	slog.Info("Shutdown signal received, stopping simulator...")
	sim.Stop()

	time.Sleep(100 * time.Millisecond)
	slog.Info("Shutdown complete")
}

func getTestPlants() []*models.Plant {
	tomato := models.PlantType{
		Name:                  "Tomato",
		OptimalSaturation:     0.6,
		MinSaturation:         0.3,
		MaxSaturation:         0.8,
		BaseGrowthRate:        0.05,
		SaturationDepletion:   0.04,
		HealthDegradationRate: 0.08,
		HealthEnhancementRate: 0.03,
	}

	lettuce := models.PlantType{
		Name:                  "Lettuce",
		OptimalSaturation:     0.7,
		MinSaturation:         0.4,
		MaxSaturation:         0.9,
		BaseGrowthRate:        0.08,
		SaturationDepletion:   0.05,
		HealthDegradationRate: 0.06,
		HealthEnhancementRate: 0.04,
	}

	var plants []*models.Plant

	tomatoPlant1, err := models.NewPlant("tomato-1", tomato, "section-A", 0.5)
	if err != nil {
		slog.Error("Error creating tomato plant", "error", err)
	} else {
		plants = append(plants, tomatoPlant1)
	}
	tomatoPlant2, err := models.NewPlant("tomato-2", tomato, "section-A", 0.3)
	if err != nil {
		slog.Error("Error creating tomato plant", "error", err)
	} else {
		plants = append(plants, tomatoPlant2)
	}

	lettucePlant1, err := models.NewPlant("lettuce-1", lettuce, "section-B", 0.6)
	if err != nil {
		slog.Error("Error creating lettuce plant", "error", err)
	} else {
		plants = append(plants, lettucePlant1)
	}

	return plants
}
