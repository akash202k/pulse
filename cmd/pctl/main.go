package main

import (
	"fmt"
	"log"

	"github.com/akash202k/pulse/internal/config"
	"github.com/akash202k/pulse/internal/probe"
)

func main() {
	fmt.Println("Pulse CLI starting...")
	cfg, err := config.Load("configs/pulse.yaml")
	if err != nil {
		log.Fatal(err)
	}

	for _, svc := range cfg.Services {
		result := probe.Check(svc)
		fmt.Printf("%s success=%v latency=%v\n", result.Service, result.Success, result.Latency)
	}
}
