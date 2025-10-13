package models

import (
	"fmt"
	"math"
	"time"
)

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

type Plant struct {
	ID             string
	Type           PlantType
	SoilSaturation float64 // 0.0 to 1.0
	Health         float64 // 0.0 (dead) to 1.0 (perfect)
	GrowthStage    float64 // 0.0 (seed) to 1.0 (mature)
	Alive          bool
	CreatedAt      time.Time
}

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
		return
	}
	growthRate := p.Type.BaseGrowthRate
	if p.Health < 0.5 {
		growthRate /= 1.35
	}
	if math.Abs(p.SoilSaturation-p.Type.OptimalSaturation) < 0.15 {
		growthRate *= 1.25
	}
	p.GrowthStage = math.Min(p.GrowthStage+growthRate, 1)
}

func updateSoilSaturation(p *Plant) {
	p.SoilSaturation = math.Max(p.SoilSaturation-p.Type.SaturationDepletion, 0)
}
