package storage

import (
	"github.com/akash202k/pulse/internal/model"
	"github.com/akash202k/pulse/internal/sli"
	"github.com/akash202k/pulse/internal/slo"
)

type Storage interface {
	// ProbeResults
	StoreProbeResult(result model.ProbeResult) error
	GetProbeResults(serviceName string, limit int) ([]model.ProbeResult, error)
	
	// SLI Metrics
	StoreSLIMetrics(metrics sli.ServiceMetrics) error
	GetSLIMetrics(serviceName string) (sli.ServiceMetrics, error)
	
	// SLO Status
	StoreSLOStatus(status slo.SLOStatus) error
	GetLatestSLOStatus(serviceName string) (slo.SLOStatus, error)
	GetAllSLOStatuses() ([]slo.SLOStatus, error)
	
	// Cleanup
	Close() error
}
