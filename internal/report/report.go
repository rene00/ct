package report

import (
	"database/sql"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3" //nolint
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/pflag"
)

type reportAll struct {
	// The metric name.
	MetricName string `json:"metric_name"`

	// The monthly value.
	MetricValue float64 `json:"metric_value"`

	// The month.
	Timestamp string `json:"month"`
}

type reportMonth struct {
	// The metric name.
	MetricName string `json:"metric_name"`

	// The monthly value.
	MetricValue float64 `json:"metric_value"`

	// The count for the month.
	Count int `json:"count"`

	// The month.
	Month string `json:"month"`
}

// All prints the all report.
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

	report := []reportAll{}

	for rows.Next() {
		var metricName string
		var metricValue float64
		var metricTimestamp string
		if err := rows.Scan(&metricName, &metricValue, &metricTimestamp); err != nil {
			return err
		}
		_reportAll := reportAll{metricName, metricValue, metricTimestamp}
		if len(metrics) != 0 {
			f := stringInSlice(metricName, metrics)
			if f {
				report = append(report, _reportAll)
			}
			continue
		}
		report = append(report, _reportAll)
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

// MonthlyAverage generates the monthly average report.
// BUG(rene): exclude bool metrics from this report.
func MonthlyAverage(db *sql.DB, metrics []string) error {
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

	report := []reportMonth{}

	for rows.Next() {
		var metricName string
		var metricValue float64
		var count int
		var month string
		if err := rows.Scan(&metricName, &metricValue, &count, &month); err != nil {
			return err
		}
		_reportMonth := reportMonth{metricName, metricValue, count, month}
		if len(metrics) != 0 {
			f := stringInSlice(metricName, metrics)
			if f {
				report = append(report, _reportMonth)
			}
			continue
		}
		report = append(report, _reportMonth)
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

// Streak prints the streak report.
func Streak(db *sql.DB, metricName string) error {
	sqlStmt := `
	SELECT 
	metric.id AS metric_id,
	ct.value AS metric_value,
	ct.id AS ct_id
	FROM ct
	INNER JOIN metric 
	ON metric.id = ct.metric_id
	WHERE metric.name = ?
	ORDER BY ct.timestamp
	DESC
	LIMIT 1
	`

	stmt, err := db.Prepare(sqlStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var metricID int
	var lastValue string
	var ctID int
	err = stmt.QueryRow(metricName).Scan(&metricID, &lastValue, &ctID)
	if err != nil {
		return err
	}

	// FIXME: This query will find all values for the metric and then walk back to tally up the streak. It will  skip a day if the metric does not exist. The code below should be updated to fill in the empty days when tallying up the streak.
	sqlStmt = `
	SELECT
	value AS metric_value
	FROM ct
	WHERE metric_id = ?
	AND id != ?
	ORDER BY timestamp
	ASC
	`

	stmt, err = db.Prepare(sqlStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(metricID, ctID)
	if err != nil {
		return err
	}
	defer rows.Close()

	streakTally := 1
	for rows.Next() {
		var i string
		if err := rows.Scan(&i); err != nil {
			return err
		}
		if i == lastValue {
			streakTally++
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Current Streak"})
	table.Append([]string{metricName, strconv.Itoa(streakTally)})
	table.Render()

	return nil
}
