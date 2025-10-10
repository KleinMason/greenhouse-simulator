package main

import (
	"greenhouse-simulator/internal/engine"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	sim := engine.NewSimulator(2 * time.Second)

	go sim.Start()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	log.Println("Shutdown signal received, stopping simulator...")
	sim.Stop()

	time.Sleep(100 * time.Millisecond)
	log.Println("Shutdown complete")
}
