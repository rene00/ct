package store

import "time"

// Metric is a metric.
type Metric struct {
	MetricID int64  `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
}

// Log is a metric log.
type Log struct {
	LogID     int64     `json:"id" db:"id"`
	MetricID  int64     `json:"metric_id" db:"metric_id"`
	Value     string    `json:"value" db:"value"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
}
