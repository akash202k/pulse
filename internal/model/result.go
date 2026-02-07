package model

import "time"

type Service struct {
	Name     string
	Url      string
	Interval time.Duration
}

type Probesult struct {
	Service   string
	Timestamp time.Time
	Success   bool
	Status    int
	Latency   time.Duration
}

type SLO struct {
	Availability float64
	LatencyP95   time.Duration
}
