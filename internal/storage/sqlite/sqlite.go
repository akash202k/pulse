package sqlite

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/akash202k/pulse/internal/model"
	"github.com/akash202k/pulse/internal/sli"
	"github.com/akash202k/pulse/internal/slo"
	"github.com/akash202k/pulse/internal/storage"
)

type SQLiteStore struct {
	db *sql.DB
}

func New(dbPath string) (storage.Storage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	store := &SQLiteStore{db: db}
	if err := store.init(); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *SQLiteStore) init() error {
	schema := `
	CREATE TABLE IF NOT EXISTS probe_results (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		service_name TEXT NOT NULL,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		success BOOLEAN,
		status INTEGER,
		latency INTEGER
	);

	CREATE TABLE IF NOT EXISTS sli_metrics (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		service_name TEXT NOT NULL UNIQUE,
		total_probes INTEGER,
		successful_probes INTEGER,
		failed_probes INTEGER,
		availability REAL,
		latency_p95 REAL,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS slo_status (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		service_name TEXT NOT NULL,
		availability_met BOOLEAN,
		latency_met BOOLEAN,
		availability_target REAL,
		availability_actual REAL,
		latency_target INTEGER,
		latency_actual INTEGER,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_probe_service ON probe_results(service_name);
	CREATE INDEX IF NOT EXISTS idx_probe_timestamp ON probe_results(timestamp);
	CREATE INDEX IF NOT EXISTS idx_slo_service ON slo_status(service_name);
	CREATE INDEX IF NOT EXISTS idx_slo_timestamp ON slo_status(timestamp);
	`

	_, err := s.db.Exec(schema)
	return err
}

func (s *SQLiteStore) StoreProbeResult(result model.ProbeResult) error {
	query := `
	INSERT INTO probe_results (service_name, timestamp, success, status, latency)
	VALUES (?, ?, ?, ?, ?)
	`

	_, err := s.db.Exec(
		query,
		result.Service,
		result.Timestamp,
		result.Success,
		result.Status,
		result.Latency.Milliseconds(),
	)
	return err
}

func (s *SQLiteStore) GetProbeResults(serviceName string, limit int) ([]model.ProbeResult, error) {
	query := `
	SELECT service_name, timestamp, success, status, latency
	FROM probe_results
	WHERE service_name = ?
	ORDER BY timestamp DESC
	LIMIT ?
	`

	rows, err := s.db.Query(query, serviceName, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []model.ProbeResult
	for rows.Next() {
		var r model.ProbeResult
		var latencyMs int64
		var ts time.Time

		if err := rows.Scan(&r.Service, &ts, &r.Success, &r.Status, &latencyMs); err != nil {
			return nil, err
		}
		r.Timestamp = ts
		r.Latency = time.Duration(latencyMs) * time.Millisecond
		results = append(results, r)
	}

	return results, rows.Err()
}

func (s *SQLiteStore) StoreSLIMetrics(metrics sli.ServiceMetrics) error {
	query := `
	INSERT OR REPLACE INTO sli_metrics 
	(service_name, total_probes, successful_probes, failed_probes, availability, latency_p95)
	VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.Exec(
		query,
		metrics.ServiceName,
		metrics.TotalProbes,
		metrics.SuccessfulProbes,
		metrics.FailedProbes,
		metrics.Availability,
		metrics.LatencyP95,
	)
	return err
}

func (s *SQLiteStore) GetSLIMetrics(serviceName string) (sli.ServiceMetrics, error) {
	query := `
	SELECT service_name, total_probes, successful_probes, failed_probes, availability, latency_p95
	FROM sli_metrics
	WHERE service_name = ?
	`

	var metrics sli.ServiceMetrics
	err := s.db.QueryRow(query, serviceName).Scan(
		&metrics.ServiceName,
		&metrics.TotalProbes,
		&metrics.SuccessfulProbes,
		&metrics.FailedProbes,
		&metrics.Availability,
		&metrics.LatencyP95,
	)

	if err == sql.ErrNoRows {
		return sli.ServiceMetrics{ServiceName: serviceName}, nil
	}
	return metrics, err
}

func (s *SQLiteStore) StoreSLOStatus(status slo.SLOStatus) error {
	query := `
	INSERT INTO slo_status 
	(service_name, availability_met, latency_met, availability_target, availability_actual, latency_target, latency_actual)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.Exec(
		query,
		status.ServiceName,
		status.AvailabilityMet,
		status.LatencyMet,
		status.AvailabilityTarget,
		status.AvailabilityActual,
		status.LatencyTarget,
		status.LatencyActual,
	)
	return err
}

func (s *SQLiteStore) GetLatestSLOStatus(serviceName string) (slo.SLOStatus, error) {
	query := `
	SELECT service_name, availability_met, latency_met, availability_target, availability_actual, latency_target, latency_actual
	FROM slo_status
	WHERE service_name = ?
	ORDER BY timestamp DESC
	LIMIT 1
	`

	var status slo.SLOStatus
	err := s.db.QueryRow(query, serviceName).Scan(
		&status.ServiceName,
		&status.AvailabilityMet,
		&status.LatencyMet,
		&status.AvailabilityTarget,
		&status.AvailabilityActual,
		&status.LatencyTarget,
		&status.LatencyActual,
	)

	if err == sql.ErrNoRows {
		return slo.SLOStatus{ServiceName: serviceName}, nil
	}
	return status, err
}

func (s *SQLiteStore) GetAllSLOStatuses() ([]slo.SLOStatus, error) {
	query := `
	SELECT service_name, availability_met, latency_met, availability_target, availability_actual, latency_target, latency_actual
	FROM (
		SELECT *, ROW_NUMBER() OVER (PARTITION BY service_name ORDER BY timestamp DESC) as rn
		FROM slo_status
	)
	WHERE rn = 1
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var statuses []slo.SLOStatus
	for rows.Next() {
		var status slo.SLOStatus
		if err := rows.Scan(
			&status.ServiceName,
			&status.AvailabilityMet,
			&status.LatencyMet,
			&status.AvailabilityTarget,
			&status.AvailabilityActual,
			&status.LatencyTarget,
			&status.LatencyActual,
		); err != nil {
			return nil, err
		}
		statuses = append(statuses, status)
	}

	return statuses, rows.Err()
}

func (s *SQLiteStore) Close() error {
	return s.db.Close()
}
