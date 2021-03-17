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

// Config is a metric config.
type Config struct {
	MetricID int64  `json:"metric_id" db:"metric_id"`
	Opt      string `json:"opt" db:"opt"`
	Val      string `json:"val" db:"val"`
}

// NewConfig creates a Config.
func NewConfig(metricID int64, opt, val string) *Config {
	return &Config{MetricID: metricID, Opt: opt, Val: val}
}

// LogComment is a log comment.
type LogComment struct {
	LogID   int64  `json:"log_id" db:"log_id"`
	Comment string `json:"comment" db:"comment"`
}

// IsDataTypeSupported returns bool on whether data-type value is supported.
func (c Config) IsDataTypeSupported() bool {
	supported := false
	for _, v := range []string{"int", "float", "bool"} {
		if v == c.Val {
			supported = true
			break
		}
	}
	return supported
}

// IsMetricTypeSupported returns bool on whether metric-type value is supported.
func (c Config) IsMetricTypeSupported() bool {
	supported := false
	for _, v := range []string{"gauge", "counter", "bool"} {
		if v == c.Val {
			supported = true
			break
		}
	}
	return supported
}
