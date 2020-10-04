package report

import (
	"context"
	"ct/internal/store"
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

	// The monthly average.
	MetricAverage float64 `json:"metric_average"`

	// The monthly sum.
	MetricSum float64 `json:"metric_sum"`

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
	ROUND(log.value, 2) AS metric_value,
	log.timestamp AS metric_timestamp
	FROM log
	INNER JOIN metric ON metric.id = log.metric_id
	ORDER BY log.timestamp
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

// MonthlyCounter generates the monthly report for counter metrics.
func MonthlyCounter(ctx context.Context, db *sql.DB, metric *store.Metric) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	rows, err := tx.QueryContext(ctx, `
	SELECT 
	ROUND(AVG(value), 2) AS metric_average,
	ROUND(SUM(value), 2) AS metric_sum,
	COUNT(1) AS metric_count,
	STRFTIME("%Y-%m", log.timestamp) AS month
	FROM log
	WHERE log.timestamp >= DATE('now', '-1 year')
	AND log.metric_id = ?
	GROUP BY month
	ORDER BY log.timestamp
`, metric.MetricID)
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Month", "Average", "Sum", "Count"})

	for rows.Next() {
		var avg float64
		var sum float64
		var count int
		var month string
		if err := rows.Scan(&avg, &sum, &count, &month); err != nil {
			return err
		}
		table.Append([]string{month, strconv.FormatFloat(avg, 'f', -1, 64), strconv.FormatFloat(sum, 'f', -1, 64), strconv.Itoa(count)})
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	table.Render()

	return tx.Commit()
}

// MonthlyGauge generates the monthly report for gauge metrics.
func MonthlyGauge(ctx context.Context, db *sql.DB, metric *store.Metric) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	rows, err := tx.QueryContext(ctx, `
	SELECT 
	ROUND(AVG(value), 2) AS metric_average,
	COUNT(1) AS metric_count,
	STRFTIME("%Y-%m", log.timestamp) AS month
	FROM log
	WHERE log.timestamp >= DATE('now', '-1 year')
	AND log.metric_id = ?
	GROUP BY month
	ORDER BY log.timestamp
`, metric.MetricID)
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Month", "Average", "Count"})

	for rows.Next() {
		var avg float64
		var count int
		var month string
		if err := rows.Scan(&avg, &count, &month); err != nil {
			return err
		}
		table.Append([]string{month, strconv.FormatFloat(avg, 'f', -1, 64), strconv.Itoa(count)})
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	table.Render()

	return tx.Commit()
}

// Streak prints the streak report.
func Streak(db *sql.DB, metricName string) error {
	sqlStmt := `
	SELECT 
	metric.id AS metric_id,
	log.value AS metric_value,
	log.id AS log_id
	FROM log
	INNER JOIN metric 
	ON metric.id = log.metric_id
	WHERE metric.name = ?
	ORDER BY log.timestamp
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
	var logID int
	err = stmt.QueryRow(metricName).Scan(&metricID, &lastValue, &logID)
	if err != nil {
		return err
	}

	// FIXME: This query will find all values for the metric and then walk back to tally up the streak. It will  skip a day if the metric does not exist. The code below should be updated to fill in the empty days when tallying up the streak.
	sqlStmt = `
	SELECT
	value AS metric_value
	FROM log
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

	rows, err := stmt.Query(metricID, logID)
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
