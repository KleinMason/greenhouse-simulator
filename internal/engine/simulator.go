package engine

import (
	"greenhouse-simulator/internal/models"
	"log"
	"sync"
	"time"
)

// Simulator defines the interface for controlling a greenhouse simulation.
// It provides methods to start, pause, resume, and stop the simulation,
// as well as manage plants within the greenhouse.
type Simulator interface {
	Start()
	Pause()
	Resume()
	Stop()
	AddPlant(p *models.Plant)
	GetPlants() []*models.Plant
	GetCurrentTick() int
}

type simulator struct {
	ticker       *time.Ticker
	pause        chan struct{}
	resume       chan struct{}
	stop         chan struct{}
	tickInterval time.Duration
	currentTick  int
	isPaused     bool
	mu           sync.RWMutex
	plants       []*models.Plant
}

// NewSimulator creates a new simulator instance with the specified tick interval.
// The tick interval determines how frequently the simulation updates.
func NewSimulator(tickInterval time.Duration) Simulator {
	return &simulator{
		ticker:       time.NewTicker(tickInterval),
		pause:        make(chan struct{}),
		resume:       make(chan struct{}),
		stop:         make(chan struct{}),
		tickInterval: tickInterval,
		currentTick:  0,
		isPaused:     false,
	}
}

// Start begins the simulation loop and runs until Stop is called.
// The simulation will process ticks at the configured interval,
// updating all plants and handling pause/resume/stop signals.
func (s *simulator) Start() {
	log.Println("Starting...")
	for {
		select {
		case <-s.ticker.C:
			log.Print("\n---------------------------------------------------------------------------\n")
			log.Printf("Tick %d\n", s.currentTick)
			s.mu.RLock()
			for _, plant := range s.plants {
				plant.OnTick()
				log.Println(plant)
			}
			s.mu.RUnlock()

			s.currentTick++
		case <-s.pause:
			log.Println("Pausing...")
			<-s.resume
			s.mu.Lock()
			s.isPaused = false
			s.mu.Unlock()
			log.Println("Resumed!")
		case <-s.stop:
			log.Println("Stopping...")
			return
		}
	}
}

// Pause temporarily halts the simulation.
// If the simulation is already paused, this method does nothing.
// The simulation can be resumed using the Resume method.
func (s *simulator) Pause() {
	s.mu.Lock()
	if s.isPaused {
		s.mu.Unlock()
		log.Println("Already paused, ignoring")
		return
	}
	s.isPaused = true
	s.mu.Unlock()
	s.pause <- struct{}{}
}

// Resume continues a paused simulation.
// If the simulation is not paused, this method does nothing.
func (s *simulator) Resume() {
	s.mu.Lock()
	if !s.isPaused {
		s.mu.Unlock()
		log.Println("Already running, ignoring")
		return
	}
	s.mu.Unlock()
	s.resume <- struct{}{}
}

// Stop terminates the simulation.
// Once stopped, the simulation cannot be resumed and must be restarted.
func (s *simulator) Stop() {
	s.stop <- struct{}{}
}

// IsPaused returns true if the simulation is currently paused, false otherwise.
// This method is safe for concurrent use.
func (s *simulator) IsPaused() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isPaused
}

// AddPlant adds a new plant to the greenhouse simulator.
// The plant will be included in the simulation starting from the next tick.
// This method is safe for concurrent use.
func (s *simulator) AddPlant(p *models.Plant) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.plants = append(s.plants, p)
}

// GetPlants returns a snapshot of all plants in the greenhouse.
// The returned slice is a copy and safe to iterate, but the plants
// themselves are shared with the simulator.
func (s *simulator) GetPlants() []*models.Plant {
	s.mu.RLock()
	defer s.mu.RUnlock()
	plantCopy := make([]*models.Plant, len(s.plants))
	copy(plantCopy, s.plants)
	return plantCopy
}

// GetCurrentTick returns the current simulation tick count.
// This is thread-safe and can be called while the simulation is running.
func (s *simulator) GetCurrentTick() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.currentTick
}
