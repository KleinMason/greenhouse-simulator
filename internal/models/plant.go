package models

import (
	"errors"
	"fmt"
	"math"
	"time"
)

// PlantType defines the characteristics and behavior parameters for a specific type of plant.
// It contains all the configuration values that determine how a plant of this type will
// grow, consume resources, and respond to environmental conditions.
type PlantType struct {
	Name                  string
	OptimalSaturation     float64
	MinSaturation         float64 // below this, health degrades
	MaxSaturation         float64 // above this, health degrades
	BaseGrowthRate        float64 // per tick
	SaturationDepletion   float64 // per tick
	HealthDegradationRate float64 // per tick if not in optimal saturation range
	HealthEnhancementRate float64 // per tick if in the optimal saturation range
}

// Plant represents an individual plant instance in the simulation.
// Each plant has its own state that changes over time based on environmental
// conditions and the characteristics defined by its PlantType.
type Plant struct {
	ID             string
	Type           PlantType
	SectionID      string
	SoilSaturation float64 // 0.0 to 1.0
	Health         float64 // 0.0 (dead) to 1.0 (perfect)
	GrowthStage    float64 // 0.0 (seed) to 1.0 (mature)
	Alive          bool
	CreatedAt      time.Time
}

// NewPlant creates a new Plant instance with the specified parameters and validates all inputs.
//
// This constructor function performs comprehensive validation on all input parameters to ensure
// the plant is created in a valid state. It validates that IDs are not empty, saturation values
// are within the valid range (0.0 to 1.0), and all PlantType configuration values are valid.
//
// Parameters:
//   - id: A unique identifier for the plant instance (cannot be empty)
//   - plantType: The PlantType configuration that defines the plant's characteristics and behavior
//   - sectionID: The identifier of the garden section where this plant is located (cannot be empty)
//   - initialSaturation: The starting soil saturation level (must be between 0.0 and 1.0)
//
// Returns:
//   - *Plant: A pointer to the newly created Plant instance with default starting values:
//   - Health: 1.0 (perfect health)
//   - GrowthStage: 0.0 (seed stage)
//   - Alive: true
//   - CreatedAt: current timestamp
//   - error: An error if any validation fails, including:
//   - Empty id or sectionID
//   - Invalid saturation values (outside 0.0-1.0 range)
//   - Invalid PlantType configuration values
//
// Example usage:
//
//	plantType := PlantType{Name: "Tomato", OptimalSaturation: 0.7, ...}
//	plant, err := NewPlant("plant-001", plantType, "section-A", 0.5)
//	if err != nil {
//	    // handle validation error
//	}
func NewPlant(id string, plantType PlantType, sectionID string, initialSaturation float64) (*Plant, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	if sectionID == "" {
		return nil, errors.New("sectionID cannot be empty")
	}
	if initialSaturation < 0 || initialSaturation > 1 {
		return nil, errors.New("initial saturation must be between 0.0 and 1.0")
	}
	if plantType.Name == "" {
		return nil, errors.New("plant type must have a name")
	}
	if plantType.OptimalSaturation < 0 || plantType.OptimalSaturation > 1 {
		return nil, errors.New("plant type optimal saturation must be between 0.0 and 1.0")
	}
	if plantType.MinSaturation < 0 || plantType.MinSaturation > 1 {
		return nil, errors.New("plant type min saturation must be between 0.0 and 1.0")
	}
	if plantType.MaxSaturation < 0 || plantType.MaxSaturation > 1 {
		return nil, errors.New("plant type max saturation must be between 0.0 and 1.0")
	}
	if plantType.BaseGrowthRate < 0 || plantType.BaseGrowthRate > 1 {
		return nil, errors.New("plant type base growth rate must be between 0.0 and 1.0")
	}
	if plantType.SaturationDepletion < 0 || plantType.SaturationDepletion > 1 {
		return nil, errors.New("plant type saturation depletion rate must be between 0.0 and 1.0")
	}
	if plantType.HealthDegradationRate < 0 || plantType.HealthDegradationRate > 1 {
		return nil, errors.New("plant type health degradation rate must be between 0.0 and 1.0")
	}
	if plantType.HealthEnhancementRate < 0 || plantType.HealthEnhancementRate > 1 {
		return nil, errors.New("plant type health enhancement rate must be between 0.0 and 1.0")
	}

	plant := Plant{
		ID:             id,
		Type:           plantType,
		SectionID:      sectionID,
		SoilSaturation: initialSaturation,
		Health:         1.0,
		GrowthStage:    0.0,
		Alive:          true,
		CreatedAt:      time.Now(),
	}

	return &plant, nil
}

const GROWTH_SLOW_FACTOR = 1.35
const GROWTH_OPTIMAL_FACTOR = 1.25

// OnTick simulates one time step in the plant's lifecycle.
// This method is called periodically to update the plant's state based on its current conditions.
//
// The tick process follows this sequence:
// 1. Skip processing if the plant is already dead
// 2. Update health based on soil saturation:
//   - Degrades health if soil saturation is outside the optimal range (MinSaturation to MaxSaturation)
//   - Enhances health if soil saturation is within the optimal range
//
// 3. Check if plant dies (health <= 0) and mark as not alive if so
// 4. Update growth stage based on health and soil conditions
// 5. Deplete soil saturation based on the plant's consumption rate
//
// This method modifies the plant's Health, GrowthStage, SoilSaturation, and potentially Alive fields.
func (p *Plant) OnTick() {
	if !p.Alive {
		return
	}
	if outOfOptimalSaturationRange(p) {
		degradeHealth(p)
	} else {
		enhanceHealth(p)
	}

	if p.Health <= 0 {
		p.Alive = false
		return
	}
	updateGrowthStage(p)
	updateSoilSaturation(p)
}

func (p *Plant) String() string {
	return fmt.Sprintf("[%s] Health:%.2f Growth:%.2f Sat:%.2f Alive:%v",
		p.ID, p.Health, p.GrowthStage, p.SoilSaturation, p.Alive)

}

func outOfOptimalSaturationRange(p *Plant) bool {
	return (p.SoilSaturation < p.Type.MinSaturation) || (p.SoilSaturation > p.Type.MaxSaturation)
}

func degradeHealth(p *Plant) {
	p.Health = math.Max(p.Health-p.Type.HealthDegradationRate, 0)
}

func enhanceHealth(p *Plant) {
	p.Health = math.Min(p.Health+p.Type.HealthEnhancementRate, 1)
}

func updateGrowthStage(p *Plant) {
	if p.Health < 0.3 {
		return // no growth
	}

	growthRate := p.Type.BaseGrowthRate // start with base rate
	if p.Health < 0.5 {
		growthRate /= GROWTH_SLOW_FACTOR // SLOWER growth
	}
	if math.Abs(p.SoilSaturation-p.Type.OptimalSaturation) < 0.15 {
		growthRate *= GROWTH_OPTIMAL_FACTOR // BONUS growth (near optimal)
	}
	p.GrowthStage = math.Min(p.GrowthStage+growthRate, 1) // Cap at 1.0
}

func updateSoilSaturation(p *Plant) {
	p.SoilSaturation = math.Max(p.SoilSaturation-p.Type.SaturationDepletion, 0)
}
