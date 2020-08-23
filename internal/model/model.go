package model

import "time"

// Metric is the main struct which holds the config and data of a metric.
type Metric struct {
	ID     int          `json:"metric_id"`
	Name   string       `json:"metric_name"`
	Config MetricConfig `json:"metric_config"`
	Data   []MetricData `json:"metric_data"`
}

// MetricConfig has the config of a metric.
type MetricConfig struct {
	ValueText string `json:"value_text,omitempty"`
	DataType  string `json:"data_type,omitempty"`
}

// MetricData has the data of a metric.
type MetricData struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
}
