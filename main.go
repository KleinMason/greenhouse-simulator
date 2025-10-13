package main

import (
	"greenhouse-simulator/internal/engine"
	"greenhouse-simulator/internal/models"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	sim := engine.NewSimulator(4 * time.Second)

	testPlants := getTestPlants()

	for i := range testPlants {
		sim.AddPlant(&testPlants[i])
	}

	go sim.Start()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	log.Println("Shutdown signal received, stopping simulator...")
	sim.Stop()

	time.Sleep(100 * time.Millisecond)
	log.Println("Shutdown complete")
}

func getTestPlants() []models.Plant {
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

	return []models.Plant{
		{
			ID:             "tomato-1",
			Type:           tomato,
			SoilSaturation: 0.6,
			Health:         1.0,
			GrowthStage:    0.0,
			Alive:          true,
			CreatedAt:      time.Now(),
		},
		{
			ID:             "lettuce-1",
			Type:           lettuce,
			SoilSaturation: 0.7,
			Health:         1.0,
			GrowthStage:    0.0,
			Alive:          true,
			CreatedAt:      time.Now(),
		},
	}
}
