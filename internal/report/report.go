package report

import (
	"database/sql"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/pflag"
)

type ReportAll struct {
	// The metric name.
	MetricName string `json:"metric_name"`

	// The monthly value.
	MetricValue float64 `json:"metric_value"`

	// The month.
	Timestamp string `json:"month"`
}

type ReportMonth struct {
	// The metric name.
	MetricName string `json:"metric_name"`

	// The monthly value.
	MetricValue float64 `json:"metric_value"`

	// The count for the month.
	Count int `json:"count"`

	// The month.
	Month string `json:"month"`
}

func All(db *sql.DB, flags *pflag.FlagSet) error {
	metrics, err := flags.GetStringSlice("metrics")
	if err != nil {
		return err
	}

	sqlStmt := `
	SELECT metric.name AS metric_name,
	ROUND(ct.value, 2) AS metric_value,
	ct.timestamp AS metric_timestamp
	FROM ct
	INNER JOIN metric ON metric.id = ct.metric_id
	ORDER BY ct.timestamp
	ASC
	`
	rows, err := db.Query(sqlStmt)
	if err != nil {
		return err
	}
	defer rows.Close()

	report := []ReportAll{}

	for rows.Next() {
		var metricName string
		var metricValue float64
		var metricTimestamp string
		if err := rows.Scan(&metricName, &metricValue, &metricTimestamp); err != nil {
			return err
		}
		reportAll := ReportAll{metricName, metricValue, metricTimestamp}
		if len(metrics) != 0 {
			f := stringInSlice(metricName, metrics)
			if f {
				report = append(report, reportAll)
			}
			continue
		}
		report = append(report, reportAll)
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	reportData := [][]string{}

	for _, r := range report {
		rd := []string{r.Timestamp, r.MetricName, strconv.FormatFloat(r.MetricValue, 'f', -1, 64)}
		reportData = append(reportData, rd)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Timestamp", "Name", "Value"})

	for _, v := range reportData {
		table.Append(v)
	}

	table.Render()

	return nil
}

func stringInSlice(s string, sl []string) bool {
	for _, i := range sl {
		if i == s {
			return true
		}
	}
	return false
}

func MonthlyAverage(db *sql.DB, flags *pflag.FlagSet) error {
	metrics, err := flags.GetStringSlice("metrics")
	if err != nil {
		return err
	}

	sqlStmt := `
	SELECT metric.name AS metric_name,
	ROUND(AVG(ct.value), 2) AS metric_value,
	COUNT(1) AS metric_count,
	STRFTIME("%Y-%m", ct.timestamp) AS month
	FROM ct
	INNER JOIN metric ON metric.id = ct.metric_id
	WHERE ct.timestamp >= DATE('now', '-1 year')
	GROUP BY metric_name, month
	ORDER BY ct.timestamp
	`
	rows, err := db.Query(sqlStmt)
	if err != nil {
		return err
	}
	defer rows.Close()

	report := []ReportMonth{}

	for rows.Next() {
		var metricName string
		var metricValue float64
		var count int
		var month string
		if err := rows.Scan(&metricName, &metricValue, &count, &month); err != nil {
			return err
		}
		reportMonth := ReportMonth{metricName, metricValue, count, month}
		if len(metrics) != 0 {
			f := stringInSlice(metricName, metrics)
			if f {
				report = append(report, reportMonth)
			}
			continue
		}
		report = append(report, reportMonth)
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	reportData := [][]string{}

	for _, r := range report {
		rd := []string{r.Month, r.MetricName, strconv.FormatFloat(r.MetricValue, 'f', -1, 64), strconv.Itoa(r.Count)}
		reportData = append(reportData, rd)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Month", "Name", "Value", "Count"})

	for _, v := range reportData {
		table.Append(v)
	}

	table.Render()

	return nil
}
