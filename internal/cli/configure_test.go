package cli

import (
	"ct/internal/store"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigCmd(t *testing.T) {
	tests := []struct {
		createMetricArgs [][]string
		cmdArgs          []string
		expectMetrics    []store.Metric
		expectConfigs    []store.Config
		desc             string
	}{
		{
			[][]string{[]string{"test"}},
			[]string{"test", "--value-text", "test1", "--data-type", "float"},
			[]store.Metric{store.Metric{MetricID: 1, Name: "test"}},
			[]store.Config{store.Config{MetricID: 1, Opt: "data_type", Val: "float"}, store.Config{MetricID: 1, Opt: "metric_type", Val: "gauge"}, store.Config{MetricID: 1, Opt: "value_text", Val: "test1"}},
			"basic metric",
		},
		{
			[][]string{[]string{"test", "--data-type=int"}},
			[]string{"test", "--data-type=float"},
			[]store.Metric{store.Metric{MetricID: 1, Name: "test"}},
			[]store.Config{store.Config{MetricID: 1, Opt: "data_type", Val: "float"}, store.Config{MetricID: 1, Opt: "metric_type", Val: "gauge"}},
			"change data-type of metric from int to float",
		},
		{
			[][]string{[]string{"test"}},
			[]string{"test", "--data-type", "bool"},
			[]store.Metric{store.Metric{MetricID: 1, Name: "test"}},
			[]store.Config{store.Config{MetricID: 1, Opt: "data_type", Val: "bool"}, store.Config{MetricID: 1, Opt: "metric_type", Val: "gauge"}},
			"change data-type of metric to bool",
		},
		{
			[][]string{[]string{"test", "--value-text=test1"}},
			[]string{"test", "--value-text=test2"},
			[]store.Metric{store.Metric{MetricID: 1, Name: "test"}},
			[]store.Config{store.Config{MetricID: 1, Opt: "data_type", Val: "float"}, store.Config{MetricID: 1, Opt: "metric_type", Val: "gauge"}, store.Config{MetricID: 1, Opt: "value_text", Val: "test2"}},
			"change value-text from test1 to test2",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			cli := initTest(t, test.createMetricArgs, nil)
			cmd := configureCmd(cli)
			cmd.SetArgs(test.cmdArgs)
			require.NoError(t, cmd.Execute())
			d := dumpTest(t, cli)
			require.Equal(t, test.expectMetrics, d.Metrics, test.desc)
			require.Equal(t, test.expectConfigs, d.Configs, test.desc)
		})
	}

}
