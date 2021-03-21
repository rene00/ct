package cli

import (
	"ct/internal/store"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListMetricCmd(t *testing.T) {
	cli := initTest(t, [][]string{[]string{"test"}}, nil)
	cmd := listMetricCmd(cli)
	require.NoError(t, cmd.Execute())
}

func TestDeleteMetricCmd(t *testing.T) {
	tests := []struct {
		createMetricArgs [][]string
		cmdArgs          []string
		desc             string
	}{
		{
			[][]string{[]string{"test"}},
			[]string{"test"},
			"delete single metric",
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			cli := initTest(t, test.createMetricArgs, nil)
			cmd := deleteMetricCmd(cli)
			cmd.SetArgs(test.cmdArgs)
			require.NoError(t, cmd.Execute())
			d := dumpTest(t, cli)
			require.Nil(t, d.Metrics, test.desc)
		})
	}
}

func TestCreateMetricCmd(t *testing.T) {
	tests := []struct {
		cmdArgs       [][]string
		expectMetrics []store.Metric
		expectConfigs []store.Config
		desc          string
	}{
		{
			[][]string{[]string{"test"}},
			[]store.Metric{store.Metric{MetricID: 1, Name: "test"}},
			[]store.Config{store.Config{MetricID: 1, Opt: "data_type", Val: "float"}, store.Config{MetricID: 1, Opt: "metric_type", Val: "gauge"}},
			"basic metric",
		},
		{
			[][]string{[]string{"test1"}, []string{"test2"}},
			[]store.Metric{store.Metric{MetricID: 1, Name: "test1"}, store.Metric{MetricID: 2, Name: "test2"}},
			[]store.Config{store.Config{MetricID: 1, Opt: "data_type", Val: "float"}, store.Config{MetricID: 1, Opt: "metric_type", Val: "gauge"}, store.Config{MetricID: 2, Opt: "data_type", Val: "float"}, store.Config{MetricID: 2, Opt: "metric_type", Val: "gauge"}},
			"multiple basic metrics",
		},
		{
			[][]string{[]string{"test", "--value-text=testValueText"}},
			[]store.Metric{store.Metric{MetricID: 1, Name: "test"}},
			[]store.Config{store.Config{MetricID: 1, Opt: "data_type", Val: "float"}, store.Config{MetricID: 1, Opt: "metric_type", Val: "gauge"}, store.Config{MetricID: 1, Opt: "value_text", Val: "testValueText"}},
			"create metric with value-text option",
		},
		{
			[][]string{[]string{"test", "--data-type=bool"}},
			[]store.Metric{store.Metric{MetricID: 1, Name: "test"}},
			[]store.Config{store.Config{MetricID: 1, Opt: "data_type", Val: "bool"}, store.Config{MetricID: 1, Opt: "metric_type", Val: "gauge"}},
			"create metrics with data-type=bool option",
		},
		{
			[][]string{[]string{"test", "--data-type=int"}},
			[]store.Metric{store.Metric{MetricID: 1, Name: "test"}},
			[]store.Config{store.Config{MetricID: 1, Opt: "data_type", Val: "int"}, store.Config{MetricID: 1, Opt: "metric_type", Val: "gauge"}},
			"create metrics with data-type=int option",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			cli := initTest(t, test.cmdArgs, nil)
			d := dumpTest(t, cli)
			require.Equal(t, test.expectMetrics, d.Metrics, test.desc)
			require.Equal(t, test.expectConfigs, d.Configs, test.desc)
		})
	}
}
