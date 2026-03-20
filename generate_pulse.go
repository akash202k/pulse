// package main

// import (
// 	"fmt"
// 	"os"

// 	"gopkg.in/yaml.v3"
// )

// type Service struct {
// 	Name     string `yaml:"name"`
// 	URL      string `yaml:"url"`
// 	Interval string `yaml:"interval"`
// }

// type SLO struct {
// 	Availability float64 `yaml:"availability"`
// 	LatencyP95   string  `yaml:"latency_p95"`
// }

// type Config struct {
// 	Services []Service `yaml:"services"`
// 	SLO      SLO       `yaml:"SLO"`
// }

// func main() {
// 	const (
// 		appCount  = 30
// 		basePort  = 3000
// 		serverCnt = 20
// 	)

// 	var services []Service

// 	for i := 1; i <= appCount; i++ {
// 		port := basePort + (i % serverCnt)

// 		svc := Service{
// 			Name:     fmt.Sprintf("app%d", i),
// 			URL:      fmt.Sprintf("http://localhost:%d/api/app%d/health", port, i),
// 			Interval: "20s",
// 		}

// 		services = append(services, svc)
// 	}

// 	cfg := Config{
// 		Services: services,
// 		SLO: SLO{
// 			Availability: 99.9,
// 			LatencyP95:   "300ms",
// 		},
// 	}

// 	out, err := yaml.Marshal(cfg)
// 	if err != nil {
// 		panic(err)
// 	}

// 	if err := os.WriteFile("configs/pulse.yaml", out, 0644); err != nil {
// 		panic(err)
// 	}

// 	fmt.Println("configs/pulse.yaml generated successfully")
// }
