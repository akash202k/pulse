package engine

import (
	"fmt"
	"time"

	"github.com/akash202k/pulse/internal/model"
)

type Scheduler struct {
	services []model.Service
	workers  chan model.Service
	results  chan model.ProbeResult
	done     chan struct{}
}

func NewScheduler(services []model.Service, workerCount int) *Scheduler {
	return &Scheduler{
		services: services,
		workers:  make(chan model.Service, workerCount),
		results:  make(chan model.ProbeResult, workerCount*2),
		done:     make(chan struct{}),
	}
}

func (s *Scheduler) Start(numWorkers int) <-chan model.ProbeResult {
	// Start workers
	for i := 0; i < numWorkers; i++ {
		go s.worker(i)
	}

	// Schedule probes
	go s.schedule()

	return s.results
}

func (s *Scheduler) schedule() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	lastProbe := make(map[string]time.Time)

	for {
		select {
		case <-s.done:
			close(s.workers)
			return
		case <-ticker.C:
			now := time.Now()
			for _, svc := range s.services {
				lastRun, exists := lastProbe[svc.Name]
				if !exists || now.Sub(lastRun) >= svc.Interval {
					select {
					case s.workers <- svc:
						lastProbe[svc.Name] = now
					case <-s.done:
						return
					default:
						// Worker queue full, skip this tick
					}
				}
			}
		}
	}
}

func (s *Scheduler) worker(id int) {
	for svc := range s.workers {
		fmt.Printf("[Worker %d] Probing %s\n", id, svc.Name)
		result := probeService(svc)
		s.results <- result
	}
}

func (s *Scheduler) Stop() {
	close(s.done)
}

func (s *Scheduler) Results() <-chan model.ProbeResult {
	return s.results
}
