package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/akash202k/pulse/internal/storage"
)

type Server struct {
	addr    string
	storage storage.Storage
	server  *http.Server
}

type HealthResponse struct {
	Status string `json:"status"`
}

type MetricsResponse struct {
	Service     string  `json:"service"`
	Availability float64 `json:"availability"`
	LatencyP95  float64 `json:"latency_p95_ms"`
	TotalProbes int     `json:"total_probes"`
	FailedProbes int     `json:"failed_probes"`
}

type SLOResponse struct {
	Service            string `json:"service"`
	AvailabilityMet    bool   `json:"availability_met"`
	LatencyMet         bool   `json:"latency_met"`
	AvailabilityTarget float64 `json:"availability_target"`
	AvailabilityActual float64 `json:"availability_actual"`
	LatencyTarget      int64  `json:"latency_target_ms"`
	LatencyActual      int64  `json:"latency_actual_ms"`
}

func New(addr string, store storage.Storage) *Server {
	return &Server{
		addr:    addr,
		storage: store,
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	
	mux.HandleFunc("/health", s.health)
	mux.HandleFunc("/api/metrics", s.metrics)
	mux.HandleFunc("/api/metrics/", s.metricsService)
	mux.HandleFunc("/api/slo", s.sloStatus)
	mux.HandleFunc("/api/slo/", s.sloStatusService)
	mux.HandleFunc("/api/probes/", s.probeResults)

	s.server = &http.Server{
		Addr:    s.addr,
		Handler: mux,
	}

	fmt.Printf("Starting API server on %s\n", s.addr)
	return s.server.ListenAndServe()
}

func (s *Server) Stop() error {
	if s.server != nil {
		return s.server.Close()
	}
	return nil
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(HealthResponse{Status: "ok"})
}

func (s *Server) metrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Use /api/metrics/{serviceName} to get metrics for a specific service",
	})
}

func (s *Server) metricsService(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	serviceName := r.URL.Path[len("/api/metrics/"):]
	if serviceName == "" {
		http.Error(w, "Service name required", http.StatusBadRequest)
		return
	}

	metrics, err := s.storage.GetSLIMetrics(serviceName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := MetricsResponse{
		Service:      metrics.ServiceName,
		Availability: metrics.Availability,
		LatencyP95:   metrics.LatencyP95,
		TotalProbes:  metrics.TotalProbes,
		FailedProbes: metrics.FailedProbes,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) sloStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	statuses, err := s.storage.GetAllSLOStatuses()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var responses []SLOResponse
	for _, status := range statuses {
		responses = append(responses, SLOResponse{
			Service:            status.ServiceName,
			AvailabilityMet:    status.AvailabilityMet,
			LatencyMet:         status.LatencyMet,
			AvailabilityTarget: status.AvailabilityTarget,
			AvailabilityActual: status.AvailabilityActual,
			LatencyTarget:      status.LatencyTarget,
			LatencyActual:      status.LatencyActual,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

func (s *Server) sloStatusService(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	serviceName := r.URL.Path[len("/api/slo/"):]
	if serviceName == "" {
		http.Error(w, "Service name required", http.StatusBadRequest)
		return
	}

	status, err := s.storage.GetLatestSLOStatus(serviceName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := SLOResponse{
		Service:            status.ServiceName,
		AvailabilityMet:    status.AvailabilityMet,
		LatencyMet:         status.LatencyMet,
		AvailabilityTarget: status.AvailabilityTarget,
		AvailabilityActual: status.AvailabilityActual,
		LatencyTarget:      status.LatencyTarget,
		LatencyActual:      status.LatencyActual,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) probeResults(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	serviceName := r.URL.Path[len("/api/probes/"):]
	if serviceName == "" {
		http.Error(w, "Service name required", http.StatusBadRequest)
		return
	}

	limit := 100
	results, err := s.storage.GetProbeResults(serviceName, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
