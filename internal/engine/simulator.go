package engine

import (
	"greenhouse-simulator/internal/models"
	"log"
	"sync"
	"time"
)

type Simulator interface {
	Start()
	Pause()
	Resume()
	Stop()
}

type simulator struct {
	ticker       *time.Ticker
	pause        chan struct{}
	resume       chan struct{}
	stop         chan struct{}
	tickInterval time.Duration
	currentTick  int
	isPaused     bool
	mu           sync.Mutex
}

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

func (s *simulator) Start() {
	for {
		select {
		case <-s.ticker.C:
			log.Print("-------------------------")
			log.Printf("Tick %d\n", s.currentTick)

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

func (s *simulator) Stop() {
	s.stop <- struct{}{}
}

func (s *simulator) IsPaused() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.isPaused
}

func testPlants() []models.Plant {
	// TODO: working here
	return nil
}
