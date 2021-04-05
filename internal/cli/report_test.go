package cli

import (
	"testing"
)

func TestDailyReportCmd(t *testing.T) {
	tests := []struct {
		createMetricArgs [][]string
		createLogArgs    [][]string
		cmdArgs          [][]string
		desc             string
	}{
		{[][]string{[]string{"test1"}}, [][]string{[]string{"test1", "1"}}, [][]string{[]string{"test1"}}, "basic report"},
		{[][]string{[]string{"test1"}, []string{"test2"}}, [][]string{[]string{"test1", "1"}, []string{"test2", "2"}}, [][]string{[]string{"test1"}, []string{"test2"}}, "report with 2 metrics"},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			cli := initTest(t, test.createMetricArgs, test.createLogArgs)
			for _, i := range test.cmdArgs {
				cmd := dailyReportCmd(cli)
				cmd.SetArgs(i)
				if _, err := cmd.ExecuteC(); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestMonthlyReportCmd(t *testing.T) {
	tests := []struct {
		createMetricArgs [][]string
		createLogArgs    [][]string
		cmdArgs          [][]string
		desc             string
	}{
		{[][]string{[]string{"test1"}}, [][]string{[]string{"test1", "1"}}, [][]string{[]string{"test1"}}, "basic report"},
		{[][]string{[]string{"test1"}, []string{"test2"}}, [][]string{[]string{"test1", "1"}, []string{"test2", "2"}}, [][]string{[]string{"test1"}, []string{"test2"}}, "report with 2 metrics"},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			cli := initTest(t, test.createMetricArgs, test.createLogArgs)
			for _, i := range test.cmdArgs {
				cmd := monthlyReportCmd(cli)
				cmd.SetArgs(i)
				if _, err := cmd.ExecuteC(); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}

func TestStreakReportCmd(t *testing.T) {
	tests := []struct {
		createMetricArgs [][]string
		createLogArgs    [][]string
		cmdArgs          [][]string
		desc             string
	}{
		{[][]string{[]string{"test1", "--data-type=bool"}}, [][]string{[]string{"test1", "yes"}}, [][]string{[]string{"test1"}}, "basic report"},
		{[][]string{[]string{"test1", "--data-type=bool"}, []string{"test2", "--data-type=bool"}}, [][]string{[]string{"test1", "yes"}, []string{"test2", "yes"}}, [][]string{[]string{"test1"}, []string{"test2"}}, "report with 2 metrics"},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			cli := initTest(t, test.createMetricArgs, test.createLogArgs)
			for _, i := range test.cmdArgs {
				cmd := streakReportCmd(cli)
				cmd.SetArgs(i)
				if _, err := cmd.ExecuteC(); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}
