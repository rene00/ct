package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// LogStorer manages metric logs.
type LogStorer interface {
	Create(context.Context, *Log) (*Log, error)
	Upsert(context.Context, *Log) (*Log, error)
	SelectOne(context.Context, int64, time.Time) (*Log, error)
	SelectLimit(context.Context, int64) ([]Log, error)
	SelectWithTimestamp(context.Context, int64, time.Time) ([]Log, error)
	SelectLast(context.Context, int64) (*Log, error)
}

// LogStore manages metric logs.
type LogStore struct {
	DB *sql.DB
}

// SelectOne returns a single log.
func (s LogStore) SelectOne(ctx context.Context, metricID int64, ts time.Time) (*Log, error) {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var ret Log
	err = tx.QueryRowContext(ctx, "SELECT id, timestamp, metric_id, value FROM log WHERE metric_id = ? AND timestamp = ?", metricID, ts.Format("2006-01-02")).Scan(&ret.LogID, &ret.Timestamp, &ret.MetricID, &ret.Value)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err != nil && err == sql.ErrNoRows {
		return nil, ErrNotFound
	}

	return &ret, tx.Commit()
}

// SelectLast returns the last inserted log.
func (s LogStore) SelectLast(ctx context.Context, metricID int64) (*Log, error) {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var ret Log
	err = tx.QueryRowContext(ctx, "SELECT id, timestamp, metric_id, value FROM log WHERE metric_id = ? ORDER BY timestamp DESC LIMIT 1", metricID).Scan(&ret.LogID, &ret.Timestamp, &ret.MetricID, &ret.Value)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err != nil && err == sql.ErrNoRows {
		return nil, ErrNotFound
	}

	return &ret, tx.Commit()
}

// SelectWithTimestamp returns a slice of logs since timestamp.
func (s LogStore) SelectWithTimestamp(ctx context.Context, metricID int64, ts time.Time) ([]Log, error) {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var ret []Log
	var rows *sql.Rows

	rows, err = tx.QueryContext(ctx, "SELECT id, metric_id, value, timestamp FROM log WHERE metric_id = ? AND timestamp >= ?", metricID, ts.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var o Log
		if err = rows.Scan(&o.LogID, &o.MetricID, &o.Value, &o.Timestamp); err != nil {
			return nil, err
		}
		ret = append(ret, o)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ret, tx.Commit()
}

// Create creates a new log item.
func (s LogStore) Create(ctx context.Context, log *Log) (*Log, error) {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var metricConfigDataType string
	if err = tx.QueryRowContext(ctx, "SELECT val FROM config WHERE metric_id = ? AND opt = ?", log.MetricID, "data_type").Scan(&metricConfigDataType); err != nil {
		return nil, err
	}

	switch metricConfigDataType {
	case "int":
		_, err := strconv.ParseInt(log.Value, 0, 0)
		if err != nil {
			return nil, errors.New("Value not an int")
		}
	case "float":
		_, err := strconv.ParseFloat(log.Value, 0)
		if err != nil {
			return nil, errors.New("Value not a float")
		}
	case "bool":
		value, err := getBoolFromValue(log.Value)
		if err != nil {
			return nil, err
		}
		log.Value = value
	default:
		return nil, errors.New("Missing data_type for metric")
	}

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO log (id, timestamp, metric_id, value) VALUES (NULL, ?, ?, ?)")
	if err != nil {
		return nil, err
	}

	res, err := stmt.ExecContext(ctx, log.Timestamp.Format("2006-01-02"), log.MetricID, log.Value)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	log.LogID, err = res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return log, nil
}

// Upsert inserts a new metric log or updates an existing log if it exists.
func (s LogStore) Upsert(ctx context.Context, log *Log) (*Log, error) {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var metricConfigDataType string
	if err = tx.QueryRowContext(ctx, "SELECT val FROM config WHERE metric_id = ? AND opt = ?", log.MetricID, "data_type").Scan(&metricConfigDataType); err != nil {
		return nil, err
	}

	switch metricConfigDataType {
	case "int":
		_, err := strconv.ParseInt(log.Value, 0, 0)
		if err != nil {
			return nil, errors.New("Value not an int")
		}
	case "float":
		_, err := strconv.ParseFloat(log.Value, 0)
		if err != nil {
			return nil, errors.New("Value not a float")
		}
	case "bool":
		value, err := getBoolFromValue(log.Value)
		if err != nil {
			return nil, err
		}
		log.Value = value
	default:
		return nil, errors.New("Missing data_type for metric")
	}

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO log (id, timestamp, metric_id, value) VALUES (NULL, ?, ?, ?) ON CONFLICT(metric_id, timestamp) DO UPDATE SET value = ?")
	if err != nil {
		return nil, err
	}

	res, err := stmt.ExecContext(ctx, log.Timestamp.Format("2006-01-02"), log.MetricID, log.Value, log.Value)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	log.LogID, err = res.LastInsertId()
	if err != nil {
		return nil, err
	}

	// log.LogID will be 0 on UPDATE so do SELECT to get the ID and assign.
	if log.LogID == 0 {
		tx, err := s.DB.BeginTx(ctx, nil)
		if err != nil {
			return nil, err
		}
		defer tx.Rollback()

		err = tx.QueryRowContext(ctx, "SELECT id FROM log WHERE metric_id = ? AND timestamp = ?", log.MetricID, log.Timestamp.Format("2006-01-02")).Scan(&log.LogID)
		if err != nil {
			return nil, err
		}
		if err = tx.Commit(); err != nil {
			return nil, err
		}
	}

	return log, nil
}

// SelectLimit returns all log rows up to limit. A limit of 0 is no limit.
func (s LogStore) SelectLimit(ctx context.Context, limit int64) ([]Log, error) {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var ret []Log
	var rows *sql.Rows

	if limit == 0 {
		rows, err = tx.QueryContext(ctx, "SELECT id, metric_id, value, timestamp FROM log")
		if err != nil {
			return nil, err
		}
	} else {
		rows, err = tx.QueryContext(ctx, "SELECT id, metric_id, value, timestamp FROM log LIMIT ?", limit)
		if err != nil {
			return nil, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		var o Log
		if err = rows.Scan(&o.LogID, &o.MetricID, &o.Value, &o.Timestamp); err != nil {
			return nil, err
		}
		ret = append(ret, o)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ret, tx.Commit()
}

func getBoolFromValue(value string) (string, error) {
	for k, v := range map[string][]string{
		"0": []string{"n", "no"},
		"1": []string{"y", "yes"},
	} {
		for _, i := range v {
			if strings.ToLower(value) == i {
				return k, nil
			}
		}
	}

	_, err := strconv.ParseBool(value)
	if err != nil {
		return "", fmt.Errorf("Value is not a bool")
	}

	return value, nil
}
