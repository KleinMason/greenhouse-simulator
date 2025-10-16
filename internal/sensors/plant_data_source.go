package sensors

import "greenhouse-simulator/internal/models"

type PlantDataSource interface {
	GetPlantsBySectionID(sectionID string) []*models.Plant
	GetAllPlants() []*models.Plant
}
