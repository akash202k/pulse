package slo

import (
	"github.com/akash202k/pulse/internal/model"
	"github.com/akash202k/pulse/internal/sli"
)

type SLOStatus struct {
	ServiceName        string
	AvailabilityMet    bool
	LatencyMet         bool
	AvailabilityTarget float64
	AvailabilityActual float64
	LatencyTarget      int64
	LatencyActual      int64
}

type SLOEvaluator struct {
	slo       model.SLO
	calculator *sli.SLICalculator
}

func NewSLOEvaluator(slo model.SLO, calc *sli.SLICalculator) *SLOEvaluator {
	return &SLOEvaluator{
		slo:        slo,
		calculator: calc,
	}
}

func (e *SLOEvaluator) Evaluate(serviceName string) SLOStatus {
	metrics := e.calculator.GetMetrics(serviceName)
	
	availabilityMet := metrics.Availability >= e.slo.Availability
	latencyMet := metrics.LatencyP95 <= float64(e.slo.LatencyP95.Milliseconds())

	return SLOStatus{
		ServiceName:        serviceName,
		AvailabilityMet:    availabilityMet,
		LatencyMet:         latencyMet,
		AvailabilityTarget: e.slo.Availability,
		AvailabilityActual: metrics.Availability,
		LatencyTarget:      e.slo.LatencyP95.Milliseconds(),
		LatencyActual:      int64(metrics.LatencyP95),
	}
}

func (e *SLOEvaluator) EvaluateAll(services []model.Service) []SLOStatus {
	var statuses []SLOStatus
	for _, svc := range services {
		statuses = append(statuses, e.Evaluate(svc.Name))
	}
	return statuses
}

func (e *SLOEvaluator) IsMet(serviceName string) bool {
	status := e.Evaluate(serviceName)
	return status.AvailabilityMet && status.LatencyMet
}
