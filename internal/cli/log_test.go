package cli

import (
	"bytes"
	"ct/internal/store"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func TestLogCreateCmdWithWrongDataType(t *testing.T) {
	cli := initTest(t, [][]string{[]string{"test", "--data-type=int"}}, nil)
	cmd := createLogCmd(cli)
	cmd.SetArgs([]string{"test", "1.1", "--update"})
	require.Error(t, cmd.Execute())
}

func TestLogCreateCmdWithUpdate(t *testing.T) {
	cli := initTest(t, [][]string{[]string{"test"}}, nil)
	cmd := createLogCmd(cli)
	cmd.SetArgs([]string{"test", "1"})
	require.NoError(t, cmd.Execute())
	cmd.SetArgs([]string{"test", "2", "--update"})
	require.NoError(t, cmd.Execute())
	now, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	expect := []store.Log{store.Log{LogID: 1, MetricID: 1, Value: "2", Timestamp: now}}
	d := dumpTest(t, cli)
	require.Equal(t, expect, d.Logs)
}

func executeCmdTest(cmd *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)
	c, err = cmd.ExecuteC()
	return c, buf.String(), err
}

func TestLogCreateCmd(t *testing.T) {
	tests := []struct {
		createMetricArgs [][]string
		cmdArgs          [][]string
		expectLogs       func() []store.Log
		desc             string
	}{
		{
			[][]string{[]string{"test"}},
			[][]string{[]string{"test", "1"}},
			func() []store.Log {
				now, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
				return []store.Log{store.Log{LogID: 1, MetricID: 1, Value: "1", Timestamp: now}}
			},
			"basic metric",
		},
		{
			[][]string{[]string{"test1"}, []string{"test2"}},
			[][]string{[]string{"test1", "1"}, []string{"test2", "2"}},
			func() []store.Log {
				now, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
				return []store.Log{store.Log{LogID: 1, MetricID: 1, Value: "1", Timestamp: now}, store.Log{LogID: 2, MetricID: 2, Value: "2", Timestamp: now}}
			},
			"basic multiple metrics",
		},
		{
			[][]string{[]string{"test"}},
			[][]string{[]string{"test", "1", "--timestamp=2020-01-01"}},
			func() []store.Log {
				t, _ := time.Parse("2006-01-02", "2020-01-01")
				return []store.Log{store.Log{LogID: 1, MetricID: 1, Value: "1", Timestamp: t}}
			},
			"single metric with timestamp",
		},
		{
			[][]string{[]string{"test", "--data-type=float"}},
			[][]string{[]string{"test", "1.1"}},
			func() []store.Log {
				now, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
				return []store.Log{store.Log{LogID: 1, MetricID: 1, Value: "1.1", Timestamp: now}}
			},
			"metric with float",
		},
		{
			[][]string{[]string{"test", "--data-type=float"}},
			[][]string{[]string{"test", "1.1", "--timestamp=2020-01-01"}},
			func() []store.Log {
				t, _ := time.Parse("2006-01-02", "2020-01-01")
				return []store.Log{store.Log{LogID: 1, MetricID: 1, Value: "1.1", Timestamp: t}}
			},
			"metric with float and timestamp",
		},
		{
			[][]string{[]string{"test", "--data-type=bool"}},
			[][]string{[]string{"test", "yes"}},
			func() []store.Log {
				now, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
				return []store.Log{store.Log{LogID: 1, MetricID: 1, Value: "1", Timestamp: now}}
			},
			"metric with bool true yes",
		},
		{
			[][]string{[]string{"test", "--data-type=bool"}},
			[][]string{[]string{"test", "no"}},
			func() []store.Log {
				now, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
				return []store.Log{store.Log{LogID: 1, MetricID: 1, Value: "0", Timestamp: now}}
			},
			"metric with bool false no",
		},
		{
			[][]string{[]string{"test", "--data-type=bool"}},
			[][]string{[]string{"test", "1"}},
			func() []store.Log {
				now, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
				return []store.Log{store.Log{LogID: 1, MetricID: 1, Value: "1", Timestamp: now}}
			},
			"metric with bool true 1",
		},
		{
			[][]string{[]string{"test", "--data-type=bool"}},
			[][]string{[]string{"test", "0"}},
			func() []store.Log {
				now, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
				return []store.Log{store.Log{LogID: 1, MetricID: 1, Value: "0", Timestamp: now}}
			},
			"metric with bool false 0",
		},
		{
			[][]string{[]string{"test", "--data-type=bool"}},
			[][]string{[]string{"test", "0", "--timestamp=2020-01-01"}},
			func() []store.Log {
				t, _ := time.Parse("2006-01-02", "2020-01-01")
				return []store.Log{store.Log{LogID: 1, MetricID: 1, Value: "0", Timestamp: t}}
			},
			"metric with bool and timestamp",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			cli := initTest(t, test.createMetricArgs, nil)
			for _, i := range test.cmdArgs {
				cmd := createLogCmd(cli)
				cmd.SetArgs(i)
				require.NoError(t, cmd.Execute())
			}
			d := dumpTest(t, cli)
			require.Equal(t, test.expectLogs(), d.Logs, test.desc)
		})
	}
}
