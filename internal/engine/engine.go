package engine

import (
	"fmt"
	"time"

	"github.com/akash202k/pulse/internal/model"
	"github.com/akash202k/pulse/internal/sli"
	"github.com/akash202k/pulse/internal/slo"
	"github.com/akash202k/pulse/internal/storage"
)

type Engine struct {
	scheduler   *Scheduler
	calculator  *sli.SLICalculator
	evaluator   *slo.SLOEvaluator
	store       storage.Storage
	services    []model.Service
	sloConfig   model.SLO
	stopCh      chan struct{}
	numWorkers  int
	metricsInterval time.Duration
}

func New(services []model.Service, sloConfig model.SLO, store storage.Storage, numWorkers int) *Engine {
	return &Engine{
		services:        services,
		sloConfig:       sloConfig,
		store:           store,
		numWorkers:      numWorkers,
		stopCh:          make(chan struct{}),
		calculator:      sli.NewSLICalculator(),
		metricsInterval: 10 * time.Second,
	}
}

func (e *Engine) Start() error {
	e.scheduler = NewScheduler(e.services, e.numWorkers)
	resultsCh := e.scheduler.Start(e.numWorkers)
	e.evaluator = slo.NewSLOEvaluator(e.sloConfig, e.calculator)

	go e.processResults(resultsCh)
	go e.periodicEvaluation()

	return nil
}

func (e *Engine) processResults(resultsCh <-chan model.ProbeResult) {
	for {
		select {
		case result, ok := <-resultsCh:
			if !ok {
				return
			}

			e.calculator.Record(result)

			if err := e.store.StoreProbeResult(result); err != nil {
				fmt.Printf("Error storing probe result: %v\n", err)
			}

			fmt.Printf("[Result] %s: success=%v status=%d latency=%v\n",
				result.Service, result.Success, result.Status, result.Latency)

		case <-e.stopCh:
			return
		}
	}
}

func (e *Engine) periodicEvaluation() {
	ticker := time.NewTicker(e.metricsInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			e.evaluateSLOs()
		case <-e.stopCh:
			return
		}
	}
}

func (e *Engine) evaluateSLOs() {
	for _, svc := range e.services {
		metrics := e.calculator.GetMetrics(svc.Name)
		status := e.evaluator.Evaluate(svc.Name)

		if err := e.store.StoreSLIMetrics(metrics); err != nil {
			fmt.Printf("Error storing SLI metrics for %s: %v\n", svc.Name, err)
		}

		if err := e.store.StoreSLOStatus(status); err != nil {
			fmt.Printf("Error storing SLO status for %s: %v\n", svc.Name, err)
		}

		met := "✓"
		if !status.AvailabilityMet || !status.LatencyMet {
			met = "✗"
		}

		fmt.Printf("[SLO %s] %s: availability=%.2f%% (target=%.2f%%), latency=%dms (target=%dms)\n",
			met, svc.Name, status.AvailabilityActual, status.AvailabilityTarget,
			status.LatencyActual, status.LatencyTarget)
	}
}

func (e *Engine) Stop() error {
	close(e.stopCh)
	e.scheduler.Stop()
	time.Sleep(1 * time.Second)
	return e.store.Close()
}
