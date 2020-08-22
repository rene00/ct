package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type LogStorer interface {
	Create(context.Context, *Log) error
	Upsert(context.Context, *Log) error
}

type LogStore struct {
	DB *sql.DB
}

func (s LogStore) Create(ctx context.Context, o *Log) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var metricConfigDataType string
	if err = tx.QueryRowContext(ctx, "SELECT val FROM config WHERE metric_id = ? AND opt = ?", o.MetricID, "data_type").Scan(&metricConfigDataType); err != nil {
		return err
	}

	switch metricConfigDataType {
	case "int":
		_, err := strconv.ParseInt(o.Value, 0, 0)
		if err != nil {
			return errors.New("Value not an int")
		}
	case "float":
		_, err := strconv.ParseFloat(o.Value, 0)
		if err != nil {
			return errors.New("Value not a float")
		}
	case "bool":
		value, err := getBoolFromValue(o.Value)
		if err != nil {
			return err
		}
		o.Value = value
	default:
		return errors.New("Missing data_type for metric")
	}

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO log (id, timestamp, metric_id, value) VALUES (NULL, ?, ?, ?)")
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, o.Timestamp.Format("2006-01-02"), o.MetricID, o.Value)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s LogStore) Upsert(ctx context.Context, o *Log) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var metricConfigDataType string
	if err = tx.QueryRowContext(ctx, "SELECT val FROM config WHERE metric_id = ? AND opt = ?", o.MetricID, "data_type").Scan(&metricConfigDataType); err != nil {
		return err
	}

	switch metricConfigDataType {
	case "int":
		_, err := strconv.ParseInt(o.Value, 0, 0)
		if err != nil {
			return errors.New("Value not an int")
		}
	case "float":
		_, err := strconv.ParseFloat(o.Value, 0)
		if err != nil {
			return errors.New("Value not a float")
		}
	case "bool":
		value, err := getBoolFromValue(o.Value)
		if err != nil {
			return err
		}
		o.Value = value
	default:
		return errors.New("Missing data_type for metric")
	}

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO log (id, timestamp, metric_id, value) VALUES (NULL, ?, ?, ?) ON CONFLICT(metric_id, timestamp) DO UPDATE SET value = ?")
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, o.Timestamp.Format("2006-01-02"), o.MetricID, o.Value, o.Value)
	if err != nil {
		return err
	}

	return tx.Commit()
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
