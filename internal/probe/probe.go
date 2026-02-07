package probe

import (
	"net/http"
	"time"

	"github.com/akash202k/pulse/internal/model"
)

func Check(service model.Service) model.ProbeResult {
	start := time.Now()

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(service.Url)
	latency := time.Since(start)

	result := model.ProbeResult{
		Service:   service.Name,
		Timestamp: time.Now(),
		Latency:   latency,
	}

	if err != nil {
		result.Success = false
		return result
	}

	result.Status = resp.StatusCode
	result.Success = resp.StatusCode >= 200 && resp.StatusCode <= 300

	return result
}
