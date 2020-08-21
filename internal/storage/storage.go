package storage

import (
	"context"
	"ct/internal/model"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3" //nolint
)

var (
	// ErrNotFound is row not found.
	ErrNotFound = errors.New("Not Found")
)

// UpsertConfig inserts or updates a config row.
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

// CreateMetric creates a metric row.
func CreateMetric(db *sql.DB, metric model.Metric) (*int64, error) {
	var sqlStmt string

	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			err = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	sqlStmt = `
	INSERT INTO metric (id, name)
	VALUES (NULL, ?)
	`
	stmt, err := tx.PrepareContext(ctx, sqlStmt)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, metric.Name)
	if err != nil {
		return nil, err
	}

	metricID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	config := metric.Config
	dataType := config.DataType
	switch config.DataType {
	case "float":
	case "int":
	case "bool":
	case "":
		dataType = "float"
	default:
		return nil, fmt.Errorf("metric config datatype not supported: %s", config.DataType)
	}

	sqlStmt = `
	INSERT INTO config (metric_id, opt, val)
	VALUES (?, ?, ?)
	`
	stmt, err = tx.Prepare(sqlStmt)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx, metricID, "data_type", dataType); err != nil {
		return nil, err
	}

	sqlStmt = `
	INSERT INTO config (metric_id, opt, val)
	VALUES (?, "frequency", "daily")
	`
	stmt, err = tx.Prepare(sqlStmt)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx, metricID); err != nil {
		return nil, err
	}

	return &metricID, nil
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

// GetMetric returns a metric row.
func GetMetric(db *sql.DB, metricName string) (*model.Metric, error) {
	getMetricSQL := `
	SELECT id, name
	FROM metric
	WHERE name = ?
	`
	stmt, err := db.Prepare(getMetricSQL)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var id int
	var name string
	err = stmt.QueryRow(metricName).Scan(&id, &name)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err != nil && err == sql.ErrNoRows {
		return nil, ErrNotFound
	}

	metric := &model.Metric{ID: id, Name: name}
	if metric.Config, err = GetMetricConfig(db, *metric); err != nil {
		return nil, err
	}

	return metric, nil
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

// GetMetricConfig returns the config for a metric.
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
