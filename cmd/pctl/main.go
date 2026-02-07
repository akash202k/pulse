package main

import (
	"fmt"
	"log"

	"github.com/akash202k/pulse/internal/config"
)

func main() {
	fmt.Println("Pulse CLI starting...")
	cfg, err := config.Load("configs/pulse.yaml")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Loaded %d Services\n", len(cfg.Services))
	fmt.Printf("Config: %+v\n", cfg)
}
