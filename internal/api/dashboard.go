package api

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"

	"github.com/akash202k/pulse/internal/report"
)

type DashboardData struct {
	GeneratedAt              string
	TotalServices            int
	ServicesMetSLO           int
	ServicesFailedSLO        int
	ComplianceClass          string
	ComplianceRateFormatted  string
	ServicesHTML             template.HTML
}

func generateDashboardHTML(rep *report.ComplianceReport) string {
	// Generate service cards HTML
	servicesHTML := generateServiceCards(rep.Services)

	// Determine compliance class
	complianceClass := "good"
	if rep.ComplianceRate < 80 {
		complianceClass = "warning"
	}
	if rep.ComplianceRate < 50 {
		complianceClass = "critical"
	}

	// Prepare data for template
	data := DashboardData{
		GeneratedAt:             rep.GeneratedAt.Format("2006-01-02 15:04:05"),
		TotalServices:           rep.TotalServices,
		ServicesMetSLO:          rep.ServicesMetSLO,
		ServicesFailedSLO:       rep.ServicesFailedSLO,
		ComplianceClass:         complianceClass,
		ComplianceRateFormatted: fmt.Sprintf("%.1f%%", rep.ComplianceRate),
		ServicesHTML:            template.HTML(servicesHTML),
	}

	// Load and parse template
	templatePath := filepath.Join("internal", "api", "dashboard.html")
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		// Fallback if template file not found
		return generateFallbackHTML(data)
	}

	// Execute template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return generateFallbackHTML(data)
	}

	return buf.String()
}

func generateServiceCards(services []report.ServiceReport) string {
	var html string
	for _, svc := range services {
		statusClass := "passed"
		statusIcon := "✓"
		if !svc.SLOMet {
			statusClass = "failed"
			statusIcon = "✗"
		}

		availabilityColor := getHealthColor(svc.AvailabilityActual, svc.AvailabilityTarget)
		latencyColor := getLatencyColor(float64(svc.LatencyActualMs), float64(svc.LatencyTargetMs))

		html += fmt.Sprintf(`
		<div class="service-card %s">
			<div class="service-header">
				<span class="status-icon %s">%s</span>
				<h3>%s</h3>
			</div>
			<div class="service-metrics">
				<div class="metric-row">
					<span class="metric-label">Availability:</span>
					<span class="metric-value %s">%.2f%% / %.2f%%</span>
					<span class="metric-status">%s</span>
				</div>
				<div class="metric-row">
					<span class="metric-label">Latency (P95):</span>
					<span class="metric-value %s">%dms / %dms</span>
					<span class="metric-status">%s</span>
				</div>
				<div class="metric-row">
					<span class="metric-label">Success Rate:</span>
					<span class="metric-value">%d / %d probes</span>
				</div>
			</div>
		</div>`,
			statusClass,
			statusClass, statusIcon,
			svc.ServiceName,
			availabilityColor, svc.AvailabilityActual, svc.AvailabilityTarget,
			statusFromBool(svc.AvailabilityMet),
			latencyColor, svc.LatencyActualMs, svc.LatencyTargetMs,
			statusFromBool(svc.LatencyMet),
			svc.SuccessfulProbes, svc.TotalProbes,
		)
	}
	return html
}

func generateFallbackHTML(data DashboardData) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
	<title>Pulse - SLO Dashboard</title>
	<style>
		body { font-family: sans-serif; padding: 20px; background: #f5f5f5; }
		.container { max-width: 1200px; margin: 0 auto; background: white; padding: 20px; border-radius: 8px; }
		h1 { color: #333; }
		.summary { display: grid; grid-template-columns: repeat(4, 1fr); gap: 20px; margin: 20px 0; }
		.card { background: #f9f9f9; padding: 15px; border-radius: 4px; }
		.card h3 { margin: 0 0 10px 0; color: #666; }
		.value { font-size: 24px; font-weight: bold; color: #667eea; }
		footer { text-align: center; color: #999; margin-top: 40px; }
	</style>
</head>
<body>
	<div class="container">
		<h1>📊 Pulse - SLO Dashboard</h1>
		<p>Generated: %s</p>
		<div class="summary">
			<div class="card">
				<h3>Total Services</h3>
				<div class="value">%d</div>
			</div>
			<div class="card">
				<h3>Services Meeting SLO</h3>
				<div class="value" style="color: #10b981;">%d</div>
			</div>
			<div class="card">
				<h3>Services Failed SLO</h3>
				<div class="value" style="color: #ef4444;">%d</div>
			</div>
			<div class="card">
				<h3>Compliance Rate</h3>
				<div class="value">%s</div>
			</div>
		</div>
		<footer>
			<p>Pulse SLO Monitoring System | Auto-refreshes every 10 seconds</p>
		</footer>
	</div>
	<script>
		setTimeout(() => location.reload(), 10000);
	</script>
</body>
</html>`, data.GeneratedAt, data.TotalServices, data.ServicesMetSLO, data.ServicesFailedSLO, data.ComplianceRateFormatted)
}

func statusFromBool(b bool) string {
	if b {
		return `<span class="metric-status pass">✓ PASS</span>`
	}
	return `<span class="metric-status fail">✗ FAIL</span>`
}

func getHealthColor(actual, target float64) string {
	if actual >= target {
		return "positive"
	}
	return "negative"
}

func getLatencyColor(actual, target float64) string {
	if actual <= target {
		return "positive"
	}
	return "negative"
}
