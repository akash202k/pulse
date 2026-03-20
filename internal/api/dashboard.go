package api

import (
	"fmt"

	"github.com/akash202k/pulse/internal/report"
)

func generateDashboardHTML(rep *report.ComplianceReport) string {
	servicesHTML := ""
	for _, svc := range rep.Services {
		statusClass := "passed"
		statusIcon := "✓"
		if !svc.SLOMet {
			statusClass = "failed"
			statusIcon = "✗"
		}

		availabilityColor := getHealthColor(svc.AvailabilityActual, svc.AvailabilityTarget)
		latencyColor := getLatencyColor(float64(svc.LatencyActualMs), float64(svc.LatencyTargetMs))

		servicesHTML += fmt.Sprintf(`
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

	complianceClass := "good"
	if rep.ComplianceRate < 80 {
		complianceClass = "warning"
	}
	if rep.ComplianceRate < 50 {
		complianceClass = "critical"
	}

	html := `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Pulse - SLO Dashboard</title>
	<style>
		* {
			margin: 0;
			padding: 0;
			box-sizing: border-box;
		}

		body {
			font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
			background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
			min-height: 100vh;
			padding: 20px;
		}

		.container {
			max-width: 1400px;
			margin: 0 auto;
		}

		header {
			background: white;
			padding: 30px;
			border-radius: 12px;
			box-shadow: 0 10px 30px rgba(0,0,0,0.1);
			margin-bottom: 30px;
		}

		header h1 {
			color: #333;
			margin-bottom: 10px;
			font-size: 32px;
		}

		header p {
			color: #666;
			font-size: 14px;
		}

		.summary {
			display: grid;
			grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
			gap: 20px;
			margin-bottom: 30px;
		}

		.summary-card {
			background: white;
			padding: 20px;
			border-radius: 12px;
			box-shadow: 0 10px 30px rgba(0,0,0,0.1);
			text-align: center;
		}

		.summary-card h3 {
			color: #666;
			font-size: 14px;
			margin-bottom: 10px;
			text-transform: uppercase;
			letter-spacing: 0.5px;
		}

		.summary-card .value {
			font-size: 28px;
			font-weight: bold;
			margin-bottom: 5px;
		}

		.compliance-rate {
			font-size: 36px !important;
			color: #667eea;
		}

		.compliance-rate.good {
			color: #10b981 !important;
		}

		.compliance-rate.warning {
			color: #f59e0b !important;
		}

		.compliance-rate.critical {
			color: #ef4444 !important;
		}

		.services-grid {
			display: grid;
			grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
			gap: 20px;
		}

		.service-card {
			background: white;
			border-radius: 12px;
			padding: 20px;
			box-shadow: 0 10px 30px rgba(0,0,0,0.1);
			transition: all 0.3s ease;
			border-left: 5px solid #667eea;
		}

		.service-card:hover {
			transform: translateY(-5px);
			box-shadow: 0 15px 40px rgba(0,0,0,0.15);
		}

		.service-card.passed {
			border-left-color: #10b981;
		}

		.service-card.failed {
			border-left-color: #ef4444;
		}

		.service-header {
			display: flex;
			align-items: center;
			margin-bottom: 15px;
			gap: 10px;
		}

		.status-icon {
			width: 30px;
			height: 30px;
			border-radius: 50%;
			display: flex;
			align-items: center;
			justify-content: center;
			color: white;
			font-weight: bold;
			font-size: 16px;
		}

		.status-icon.passed {
			background: #10b981;
		}

		.status-icon.failed {
			background: #ef4444;
		}

		.service-header h3 {
			color: #333;
			font-size: 18px;
		}

		.service-metrics {
			display: flex;
			flex-direction: column;
			gap: 12px;
		}

		.metric-row {
			display: flex;
			justify-content: space-between;
			align-items: center;
			padding: 8px 0;
			border-bottom: 1px solid #f0f0f0;
			font-size: 14px;
		}

		.metric-row:last-child {
			border-bottom: none;
		}

		.metric-label {
			color: #666;
			font-weight: 500;
		}

		.metric-value {
			color: #333;
			font-weight: bold;
			font-family: 'Courier New', monospace;
		}

		.metric-value.positive {
			color: #10b981;
		}

		.metric-value.negative {
			color: #ef4444;
		}

		.metric-status {
			font-size: 12px;
			font-weight: bold;
			padding: 2px 8px;
			border-radius: 4px;
			background: #f0f0f0;
		}

		.metric-status.pass {
			background: #d1fae5;
			color: #065f46;
		}

		.metric-status.fail {
			background: #fee2e2;
			color: #7f1d1d;
		}

		footer {
			text-align: center;
			color: white;
			margin-top: 40px;
			font-size: 12px;
		}

		@media (max-width: 768px) {
			header {
				padding: 20px;
			}

			header h1 {
				font-size: 24px;
			}

			.services-grid {
				grid-template-columns: 1fr;
			}

			.summary {
				grid-template-columns: repeat(2, 1fr);
			}
		}
	</style>
</head>
<body>
	<div class="container">
		<header>
			<h1>📊 Pulse - SLO Dashboard</h1>
			<p>Generated: ` + rep.GeneratedAt.Format("2006-01-02 15:04:05") + `</p>
		</header>

		<div class="summary">
			<div class="summary-card">
				<h3>Total Services</h3>
				<div class="value">` + fmt.Sprintf("%d", rep.TotalServices) + `</div>
			</div>
			<div class="summary-card">
				<h3>Services Meeting SLO</h3>
				<div class="value" style="color: #10b981;">` + fmt.Sprintf("%d", rep.ServicesMetSLO) + `</div>
			</div>
			<div class="summary-card">
				<h3>Services Failed SLO</h3>
				<div class="value" style="color: #ef4444;">` + fmt.Sprintf("%d", rep.ServicesFailedSLO) + `</div>
			</div>
			<div class="summary-card">
				<h3>Compliance Rate</h3>
				<div class="value compliance-rate ` + complianceClass + `">` + fmt.Sprintf("%.1f%%", rep.ComplianceRate) + `</div>
			</div>
		</div>

		<div class="services-grid">
			` + servicesHTML + `
		</div>

		<footer>
			<p>Pulse SLO Monitoring System | Auto-refreshes every 60 seconds</p>
		</footer>
	</div>

	<script>
		// Auto-refresh every 60 seconds
		setTimeout(() => location.reload(), 60000);
	</script>
</body>
</html>`

	return html
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
