package engine

import (
	"net/http"
	"time"

	"github.com/akash202k/pulse/internal/model"
)

func probeService(svc model.Service) model.ProbeResult {
	start := time.Now()

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(svc.Url)
	latency := time.Since(start)

	result := model.ProbeResult{
		Service:   svc.Name,
		Timestamp: time.Now(),
		Latency:   latency,
	}

	if err != nil {
		result.Success = false
		return result
	}

	defer resp.Body.Close()

	result.Status = resp.StatusCode
	result.Success = resp.StatusCode >= 200 && resp.StatusCode < 300

	return result
}
