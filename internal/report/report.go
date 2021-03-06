package report

import (
	"bytes"
	"context"
	"ct/internal/store"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3" //nolint
	"github.com/olekukonko/tablewriter"
	"github.com/snabb/isoweek"
)

type reportDaily struct {
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

// Daily prints the daily report.
func Daily(ctx context.Context, db *sql.DB, metric *store.Metric) (string, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	rows, err := tx.QueryContext(ctx, `
	SELECT
	id,
	ROUND(value, 2),
	STRFTIME("%Y-%m-%d", timestamp)
	FROM log
	WHERE metric_id = ?
	ORDER BY log.timestamp
	ASC
	`, metric.MetricID)
	if err != nil {
		return "", err
	}

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetHeader([]string{"Month", "Value", "Comment"})

	s := store.NewStore(db)
	for rows.Next() {
		var id int64
		var value float64
		var timestamp string
		if err := rows.Scan(&id, &value, &timestamp); err != nil {
			return "", err
		}

		comment := ""
		logComment, err := s.LogComment.SelectOne(ctx, id)
		if err != nil && err != store.ErrNotFound {
			return "", err
		}
		if err == nil {
			comment = logComment.Comment
		}

		table.Append([]string{timestamp, strconv.FormatFloat(value, 'f', -1, 64), comment})
	}

	err = rows.Err()
	if err != nil {
		return "", err
	}

	table.Render()

	return tableString.String(), tx.Commit()
}

func stringInSlice(s string, sl []string) bool {
	for _, i := range sl {
		if i == s {
			return true
		}
	}
	return false
}

// Report is the main struct for reports.
type Report struct {
	db     *sql.DB
	metric *store.Metric
}

// NewReport returns a new Report.
func NewReport(db *sql.DB, metric *store.Metric) Report {
	return Report{db: db, metric: metric}
}

// QueryText is a struct that is intended to help build SQL strings.
type QueryText struct {
	slt         string
	from        string
	where       []string
	whereMarker bool
	groupBy     string
	orderBy     string
}

// QueryOption is a type that closes over QueryText, providing a way to mutate queries.
type QueryOption func(*QueryText)

// WithStartTimestamp sets the WHERE clause of the SQL for log timestamps greater than or equal to t.
func WithStartTimestamp(t time.Time) QueryOption {
	return func(q *QueryText) {
		q.where = append(q.where, fmt.Sprintf("log.timestamp >= \"%s\"", t.Format("2006-01-02")))
	}
}

// WithEndTimestamp sets the WHERE clause of the SQL for log timestamps less than or equal to t.
func WithEndTimestamp(t time.Time) QueryOption {
	return func(q *QueryText) {
		q.where = append(q.where, fmt.Sprintf("log.timestamp <= \"%s\"", t.Format("2006-01-02")))
	}
}

// NewQueryText returns a new QueryText.
func NewQueryText(slt, from, groupBy, orderBy string, where []string) *QueryText {
	return &QueryText{slt: slt, from: from, where: where, groupBy: groupBy, orderBy: orderBy}
}

func (q *QueryText) String() string {
	buf := &bytes.Buffer{}
	buf.WriteString(q.slt + " ")
	buf.WriteString(q.from + " ")
	buf.WriteString("WHERE " + strings.Join(q.where, " AND ") + " ")
	buf.WriteString(q.groupBy + " ")
	buf.WriteString(q.orderBy)
	return buf.String()
}

// MonthlyCounter generates the monthly report for counter metrics.
func (r Report) MonthlyCounter(ctx context.Context, options ...QueryOption) (output string, err error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	q := NewQueryText(
		"SELECT ROUND(AVG(value), 2) AS metric_average, ROUND(SUM(value), 2) AS metric_sum, COUNT(1) AS metric_count, STRFTIME(\"%Y-%m\", log.timestamp) AS month",
		"FROM log",
		"GROUP BY month",
		"ORDER BY log.timestamp",
		[]string{"log.metric_id = ?"},
	)

	for _, option := range options {
		option(q)
	}

	if len(q.where) == 1 {
		q.where = append(q.where, "log.timestamp >= DATE('now', '-1 year')")
	}

	rows, err := tx.QueryContext(ctx, fmt.Sprint(q), r.metric.MetricID)
	if err != nil {
		return "", err
	}

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetHeader([]string{"Month", "Average", "Sum", "Count"})

	for rows.Next() {
		var avg float64
		var sum float64
		var count int
		var month string
		if err := rows.Scan(&avg, &sum, &count, &month); err != nil {
			return "", err
		}
		table.Append([]string{month, strconv.FormatFloat(avg, 'f', -1, 64), strconv.FormatFloat(sum, 'f', -1, 64), strconv.Itoa(count)})
	}

	err = rows.Err()
	if err != nil {
		return "", err
	}

	tx.Commit()

	table.Render()

	return tableString.String(), nil
}

// MonthlyGuage generates the monthly report for gauge metrics.
func (r Report) MonthlyGuage(ctx context.Context, options ...QueryOption) (output string, err error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	q := NewQueryText(
		"SELECT ROUND(AVG(value), 2) AS metric_average, COUNT(1) AS metric_count, STRFTIME(\"%Y-%m\", log.timestamp) AS month",
		"FROM log",
		"GROUP BY month",
		"ORDER BY log.timestamp",
		[]string{"log.metric_id = ?"},
	)

	for _, option := range options {
		option(q)
	}

	if len(q.where) == 1 {
		q.where = append(q.where, "log.timestamp >= DATE('now', '-1 year')")
	}

	rows, err := tx.QueryContext(ctx, fmt.Sprint(q), r.metric.MetricID)
	if err != nil {
		return "", err
	}

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetHeader([]string{"Month", "Average", "Count"})

	for rows.Next() {
		var avg float64
		var count int
		var month string
		if err := rows.Scan(&avg, &count, &month); err != nil {
			return "", err
		}
		table.Append([]string{month, strconv.FormatFloat(avg, 'f', -1, 64), strconv.Itoa(count)})
	}

	err = rows.Err()
	if err != nil {
		return "", err
	}

	tx.Commit()

	table.Render()

	return tableString.String(), nil
}

// WeeklyGauge generates the weekly report for gauge metrics.
func WeeklyGauge(ctx context.Context, db *sql.DB, metric *store.Metric) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	rows, err := tx.QueryContext(ctx, `
	SELECT 
	ROUND(AVG(value), 2) AS metric_average,
	COUNT(1) AS metric_count,
	STRFTIME("%Y-%W", log.timestamp) AS year_week,
	CAST(STRFTIME("%Y", log.timestamp) AS INTEGER) AS year,
	CAST(STRFTIME("%W", log.timestamp) AS INTEGER) AS week
	FROM log
	WHERE log.timestamp >= DATE('now', '-1 year')
	AND log.metric_id = ?
	GROUP BY year_week
	ORDER BY log.timestamp
`, metric.MetricID)
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Week", "Average", "Count"})

	for rows.Next() {
		var avg float64
		var count int
		var yearWeek string
		var week int
		var year int
		if err := rows.Scan(&avg, &count, &yearWeek, &year, &week); err != nil {
			return err
		}

		table.Append([]string{isoweek.StartTime(year, week, time.UTC).Format("2006-01-02"), strconv.FormatFloat(avg, 'f', -1, 64), strconv.Itoa(count)})
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	table.Render()

	return tx.Commit()
}

// WeeklyCounter generates the weekly report for counter metrics.
func WeeklyCounter(ctx context.Context, db *sql.DB, metric *store.Metric) error {
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
	STRFTIME("%Y-%W", log.timestamp) AS year_week,
	CAST(STRFTIME("%Y", log.timestamp) AS INTEGER) AS year,
	CAST(STRFTIME("%W", log.timestamp) AS INTEGER) AS week
	FROM log
	WHERE log.timestamp >= DATE('now', '-1 year')
	AND log.metric_id = ?
	GROUP BY week
	ORDER BY log.timestamp
`, metric.MetricID)
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Week", "Average", "Sum", "Count"})

	for rows.Next() {
		var avg float64
		var sum float64
		var count int
		var yearWeek string
		var week int
		var year int
		if err := rows.Scan(&avg, &sum, &count, &yearWeek, &year, &week); err != nil {
			return err
		}

		table.Append([]string{isoweek.StartTime(year, week, time.UTC).Format("2006-01-02"), strconv.FormatFloat(avg, 'f', -1, 64), strconv.FormatFloat(sum, 'f', -1, 64), strconv.Itoa(count)})
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
