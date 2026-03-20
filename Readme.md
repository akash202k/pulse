# Pulse

SLO (Service Level Objective) monitoring system written in Go. Continuously probes HTTP services, calculates SLIs (Service Level Indicators), evaluates SLOs, and provides REST API for metrics.

## Features

- **Periodic Health Probes**: Configurable interval HTTP health checks for multiple services
- **Concurrent Probing**: Worker pool architecture for efficient parallel requests
- **SLI Calculation**: Availability and latency percentile metrics
- **SLO Evaluation**: Real-time validation against defined targets
- **Persistent Storage**: SQLite database for all probe results and metrics
- **REST API**: Query metrics, SLO status, and historical probe data
- **Beautiful Dashboard**: Real-time HTML dashboard with SLO compliance visualization
- **Reporting**: JSON API for compliance reports directly from database
- **Graceful Shutdown**: Signal handling for clean termination

## Quick Start

### Build
```bash
go build -o pctl ./cmd/pctl/main.go
```

### Configure
Edit `configs/pulse.yaml` to define services and SLO targets:
```yaml
services:
  - name: app1
    url: http://localhost:3001/api/app1/health
    interval: 20s
SLO:
  availability: 99.9
  latency_p95: 300ms
```

### Run
```bash
./pctl
```

The application will:
1. Start probing all configured services every N seconds
2. Store results in `pulse.db`
3. Calculate SLI metrics every 10 seconds
4. Evaluate against SLO targets
5. Expose metrics via HTTP API on `:9090`

## API Endpoints

- `GET /health` - Health check
- `GET /api/metrics/{serviceName}` - Get SLI metrics for a service
- `GET /api/slo` - Get all latest SLO statuses
- `GET /api/slo/{serviceName}` - Get SLO status for a service
- `GET /api/probes/{serviceName}` - Get recent probe results (last 100)
- `GET /api/report` - Get compliance report (JSON) with all service details
- `GET /api/report/json` - Download compliance report as JSON file
- `GET /dashboard` - View real-time HTML dashboard with compliance metrics

### Example Queries
```bash
# Health check
curl http://localhost:9090/health

# Get metrics for app1
curl http://localhost:9090/api/metrics/app1

# Get all SLO statuses
curl http://localhost:9090/api/slo | jq .

# Get compliance report
curl http://localhost:9090/api/report | jq .

# Get probe history
curl http://localhost:9090/api/probes/app1 | jq .

# View dashboard in browser
open http://localhost:9090/dashboard
```

### Dashboard Features
The dashboard (`/dashboard`) provides a real-time view of SLO compliance:
- **Summary cards** - Total services, compliance rate, met/failed counts
- **Service cards** - Individual availability and latency status for each service
- **Health indicators** - Color-coded (✓/✗) visual feedback
- **Auto-refresh** - Updates every 60 seconds automatically
- **Responsive design** - Works on desktop and mobile devices

## Architecture

### Components

- **Engine** - Orchestrates scheduler, SLI calculator, SLO evaluator
- **Scheduler** - Distributes probes to worker pool based on service intervals
- **Workers** - Execute HTTP health checks concurrently
- **SLI Calculator** - Computes availability and latency percentiles
- **SLO Evaluator** - Validates metrics against targets
- **Storage** - Interface for persisting metrics (SQLite implementation)
- **Report Generator** - Creates compliance reports from database metrics
- **API Server** - REST endpoints for metrics queries and dashboard
- **Dashboard** - Real-time HTML dashboard with visual SLO compliance

### Data Flow
```
Config → Engine → Scheduler → Workers → Probe Results
         ↓        ↓          ↓
      SLI Calc → SLO Eval → Storage → API
```

## 📁 Project Structure

```text
pulse/
├── cmd/
│   └── pctl/
│       └── main.go          # CLI entrypoint
│
├── internal/
│   ├── config/             # pulse.yaml parsing + validation
│   │   └── config.go
│   │
│   ├── engine/             # core runtime
│   │   ├── engine.go       # orchestrator
│   │   ├── scheduler.go    # periodic scheduling
│   │   └── worker.go       # probe execution
│   │
│   ├── probe/              # HTTP probing logic
│   │   └── probe.go
│   │
│   ├── model/              # shared structs
│   │   └── result.go       # ProbeResult, Service, SLO
│   │
│   ├── storage/            # DB interface + sqlite impl
│   │   ├── storage.go      # Storage interface
│   │   └── sqlite/
│   │       └── sqlite.go   # SQLite implementation
│   │
│   ├── sli/                # SLI calculations
│   │   └── sli.go
│   │
│   ├── slo/                # SLO evaluation
│   │   └── slo.go
│   │
│   ├── report/             # Report generation
│   │   └── report.go       # Compliance reports from DB
│   │
│   └── api/                # HTTP API & Dashboard
│       ├── server.go       # API endpoints
│       └── dashboard.go    # HTML dashboard generator
│
├── test/
│   └── test-server.js      # Dummy test servers
│
├── configs/
│   └── pulse.yaml          # example config
│
├── go.mod
└── README.md
```
