package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/akash202k/pulse/internal/api"
	"github.com/akash202k/pulse/internal/config"
	"github.com/akash202k/pulse/internal/engine"
	"github.com/akash202k/pulse/internal/storage/sqlite"
)

func main() {
	fmt.Println("Pulse CLI starting...")
	
	cfg, err := config.Load("configs/pulse.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// Initialize storage
	store, err := sqlite.New("pulse.db")
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	defer store.Close()

	// Create engine
	eng := engine.New(cfg.Services, cfg.SLO, store, 5)
	if err := eng.Start(); err != nil {
		log.Fatalf("Failed to start engine: %v", err)
	}

	// Start API server
	apiServer := api.New(":9090", store)
	go func() {
		if err := apiServer.Start(); err != nil {
			log.Printf("API server error: %v", err)
		}
	}()

	// Wait for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Pulse is running. Press Ctrl+C to stop...")
	<-sigCh

	fmt.Println("\nShutting down...")
	apiServer.Stop()
	eng.Stop()
	time.Sleep(500 * time.Millisecond)
	fmt.Println("Pulse stopped.")
}
