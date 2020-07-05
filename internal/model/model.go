package model

import "time"

type Metric struct {
	ID     int          `json:"metric_id"`
	Name   string       `json:"metric_name"`
	Config MetricConfig `json:"metric_config"`
	Data   []MetricData `json:"metric_data"`
}

type MetricConfig struct {
	Frequency string `json:"frequency,omitempty"`
	ValueText string `json:"value_text,omitempty"`
	DataType  string `json:"data_type,omitempty"`
}

type MetricData struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
}
