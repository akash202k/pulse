# Pulse

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
│   │   ├── engine.go       # starts scheduler + workers
│   │   ├── scheduler.go    # tick logic
│   │   └── worker.go       # probe goroutines
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
│   │       └── sqlite.go
│   │
│   ├── sli/                # SLI calculations
│   │   └── sli.go
│   │
│   ├── slo/                # SLO evaluation
│   │   └── slo.go
│   │
│   ├── api/                # HTTP API
│   │   └── server.go
│   │
│   └── events/             # internal event bus (future plugins)
│       └── bus.go
│
├── web/                    # UI later
│
├── configs/
│   └── pulse.yaml          # example config
│
├── docs/
│
├── go.mod
├── README.md
└── Makefile
```
