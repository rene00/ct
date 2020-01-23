package cmd

type Metric struct {
	ID     int
	Name   string
	Config MetricConfig
}

type MetricConfig struct {
	Frequency string
	ValueText string
	DataType  string
}
