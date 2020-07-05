package storage

import (
	"ct/internal/model"
	"database/sql"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type ReportMonth struct {
	// The month.
	Month time.Time `json:"month"`

	// The metric name.
	Metric string `json:"metric_name"`

	// The monthly value.
	Value float64 `json:"metric_value"`
}

// ReportLastYear produces a report for all metrics for the last year.
func ReportLastYear(db *sql.DB) error {
	sqlStmt := `
	SELECT id, name
	FROM metric
`
	rows, err := db.Query(sqlStmt)
	if err != nil {
		return err
	}
	defer rows.Close()

	metrics := []model.Metric{}

	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return err
		}
		metric := model.Metric{ID: id, Name: name}
		metrics = append(metrics, metric)
	}
	return nil
}

func UpsertConfig(db *sql.DB, metricID int, opt, val string) error {
	sqlStmt := `
		INSERT INTO config
			(
				metric_id,
				opt,
				val
			)
			VALUES
			(
				?,
				?,
				?
			)
			ON CONFLICT(metric_id, opt)
			DO UPDATE SET val=?
		`
	stmt, err := db.Prepare(sqlStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(metricID, opt, val, val); err != nil {
		return err
	}
	return nil
}

// SetMetricID inserts the metric into the metric table.
func SetMetricID(db *sql.DB, metric model.Metric) error {
	var sqlStmt string

	sqlStmt = `
	INSERT INTO metric
		(
			id,
			name
		)
		VALUES
		(
			NULL,
			?
		)
	`
	stmt, err := db.Prepare(sqlStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err = stmt.Exec(metric.Name); err != nil {
		return err
	}
	return nil
}

// GetMetricID will call setMetricID to create the metric within the table and then call itself to return the metric ID.
func GetMetricID(db *sql.DB, metric model.Metric) (int, error) {
	var sqlStmt string
	var metricID int

	sqlStmt = `
	SELECT id
		FROM metric
		WHERE name = ?
	`
	stmt, err := db.Prepare(sqlStmt)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	if err = stmt.QueryRow(metric.Name).Scan(&metricID); err != nil {
		if err = SetMetricID(db, metric); err != nil {
			return 0, err
		}
	}
	if metricID == 0 {
		return GetMetricID(db, metric)
	}
	return metricID, nil
}

func GetMetricConfig(db *sql.DB, metric model.Metric) (model.MetricConfig, error) {
	var sqlStmt string
	metricConfig := model.MetricConfig{}

	sqlStmt = `
	SELECT opt, val
		FROM CONFIG
		WHERE metric_id = ?
	`
	rows, err := db.Query(sqlStmt, metric.ID)
	if err != nil {
		return metricConfig, err
	}
	defer rows.Close()

	for rows.Next() {
		var option string
		var value string
		if err = rows.Scan(&option, &value); err != nil {
			return metricConfig, err
		}
		switch option {
		case "frequency":
			metricConfig.Frequency = value
		case "value_text":
			metricConfig.ValueText = value
		case "data_type":
			metricConfig.DataType = value
		default:
			return metricConfig, errors.New("Unsupported metric config")
		}
	}
	return metricConfig, nil
}
