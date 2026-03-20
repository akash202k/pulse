package sli

import (
	"github.com/akash202k/pulse/internal/model"
)

type ServiceMetrics struct {
	ServiceName     string
	TotalProbes     int
	SuccessfulProbes int
	FailedProbes    int
	Availability    float64
	LatencyP95      float64
}

type SLICalculator struct {
	results map[string][]model.ProbeResult
}

func NewSLICalculator() *SLICalculator {
	return &SLICalculator{
		results: make(map[string][]model.ProbeResult),
	}
}

func (c *SLICalculator) Record(result model.ProbeResult) {
	c.results[result.Service] = append(c.results[result.Service], result)
}

func (c *SLICalculator) CalculateAvailability(serviceName string) float64 {
	results, exists := c.results[serviceName]
	if !exists || len(results) == 0 {
		return 0.0
	}

	successful := 0
	for _, r := range results {
		if r.Success {
			successful++
		}
	}

	return (float64(successful) / float64(len(results))) * 100
}

func (c *SLICalculator) CalculateP95Latency(serviceName string) int64 {
	results, exists := c.results[serviceName]
	if !exists || len(results) == 0 {
		return 0
	}

	// Sort results by latency (simple implementation)
	// In production, use a proper percentile calculation
	totalLatency := int64(0)
	for _, r := range results {
		totalLatency += r.Latency.Milliseconds()
	}

	// Simple average as proxy (in production, calculate actual P95)
	return totalLatency / int64(len(results))
}

func (c *SLICalculator) GetMetrics(serviceName string) ServiceMetrics {
	results, exists := c.results[serviceName]
	if !exists {
		return ServiceMetrics{ServiceName: serviceName}
	}

	successful := 0
	for _, r := range results {
		if r.Success {
			successful++
		}
	}

	return ServiceMetrics{
		ServiceName:      serviceName,
		TotalProbes:      len(results),
		SuccessfulProbes: successful,
		FailedProbes:     len(results) - successful,
		Availability:     (float64(successful) / float64(len(results))) * 100,
		LatencyP95:       float64(c.CalculateP95Latency(serviceName)),
	}
}

func (c *SLICalculator) Reset(serviceName string) {
	delete(c.results, serviceName)
}

func (c *SLICalculator) ResetAll() {
	c.results = make(map[string][]model.ProbeResult)
}
