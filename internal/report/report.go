package report

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/akash202k/pulse/internal/model"
	"github.com/akash202k/pulse/internal/slo"
	"github.com/akash202k/pulse/internal/storage"
)

type ServiceReport struct {
	ServiceName        string              `json:"service_name"`
	AvailabilityTarget float64             `json:"availability_target"`
	AvailabilityActual float64             `json:"availability_actual"`
	AvailabilityMet    bool                `json:"availability_met"`
	LatencyTargetMs    int64               `json:"latency_target_ms"`
	LatencyActualMs    int64               `json:"latency_actual_ms"`
	LatencyMet         bool                `json:"latency_met"`
	TotalProbes        int                 `json:"total_probes"`
	SuccessfulProbes   int                 `json:"successful_probes"`
	FailedProbes       int                 `json:"failed_probes"`
	SLOMet             bool                `json:"slo_met"`
}

type ComplianceReport struct {
	GeneratedAt      time.Time        `json:"generated_at"`
	ReportPeriod     string           `json:"report_period"`
	TotalServices    int              `json:"total_services"`
	ServicesMetSLO   int              `json:"services_met_slo"`
	ServicesFailedSLO int             `json:"services_failed_slo"`
	ComplianceRate   float64          `json:"compliance_rate"`
	Services         []ServiceReport  `json:"services"`
}

type ReportGenerator struct {
	store storage.Storage
}

func New(store storage.Storage) *ReportGenerator {
	return &ReportGenerator{
		store: store,
	}
}

func (rg *ReportGenerator) GenerateComplianceReport(services []model.Service) (*ComplianceReport, error) {
	report := &ComplianceReport{
		GeneratedAt: time.Now(),
		ReportPeriod: fmt.Sprintf("%s", time.Now().Format("2006-01-02 15:04:05")),
		Services:    []ServiceReport{},
	}

	sloStatuses, err := rg.store.GetAllSLOStatuses()
	if err != nil {
		return nil, err
	}

	statusMap := make(map[string]slo.SLOStatus)
	for _, status := range sloStatuses {
		statusMap[status.ServiceName] = status
	}

	metCount := 0
	for _, svc := range services {
		status, exists := statusMap[svc.Name]
		if !exists {
			continue
		}

		metrics, err := rg.store.GetSLIMetrics(svc.Name)
		if err != nil {
			continue
		}

		sloMet := status.AvailabilityMet && status.LatencyMet
		if sloMet {
			metCount++
		}

		serviceReport := ServiceReport{
			ServiceName:        svc.Name,
			AvailabilityTarget: status.AvailabilityTarget,
			AvailabilityActual: status.AvailabilityActual,
			AvailabilityMet:    status.AvailabilityMet,
			LatencyTargetMs:    status.LatencyTarget,
			LatencyActualMs:    status.LatencyActual,
			LatencyMet:         status.LatencyMet,
			TotalProbes:        metrics.TotalProbes,
			SuccessfulProbes:   metrics.SuccessfulProbes,
			FailedProbes:       metrics.FailedProbes,
			SLOMet:             sloMet,
		}

		report.Services = append(report.Services, serviceReport)
	}

	report.TotalServices = len(services)
	report.ServicesMetSLO = metCount
	report.ServicesFailedSLO = len(services) - metCount

	if report.TotalServices > 0 {
		report.ComplianceRate = (float64(metCount) / float64(len(services))) * 100
	}

	return report, nil
}

func (rg *ReportGenerator) ExportJSON(report *ComplianceReport) ([]byte, error) {
	return json.MarshalIndent(report, "", "  ")
}

func (rg *ReportGenerator) GenerateAndExportJSON(services []model.Service) ([]byte, error) {
	report, err := rg.GenerateComplianceReport(services)
	if err != nil {
		return nil, err
	}
	return rg.ExportJSON(report)
}

func (rg *ReportGenerator) GetServiceReport(serviceName string) (*ServiceReport, error) {
	status, err := rg.store.GetLatestSLOStatus(serviceName)
	if err != nil {
		return nil, err
	}

	metrics, err := rg.store.GetSLIMetrics(serviceName)
	if err != nil {
		return nil, err
	}

	sloMet := status.AvailabilityMet && status.LatencyMet

	return &ServiceReport{
		ServiceName:        serviceName,
		AvailabilityTarget: status.AvailabilityTarget,
		AvailabilityActual: status.AvailabilityActual,
		AvailabilityMet:    status.AvailabilityMet,
		LatencyTargetMs:    status.LatencyTarget,
		LatencyActualMs:    status.LatencyActual,
		LatencyMet:         status.LatencyMet,
		TotalProbes:        metrics.TotalProbes,
		SuccessfulProbes:   metrics.SuccessfulProbes,
		FailedProbes:       metrics.FailedProbes,
		SLOMet:             sloMet,
	}, nil
}
