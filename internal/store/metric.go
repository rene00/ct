package store

import (
	"context"
	"database/sql"
)

type MetricStorer interface {
	Create(context.Context, string) (*Metric, error)
	SelectOne(context.Context, string) (*Metric, error)
}

type MetricStore struct {
	DB *sql.DB
}

func (s MetricStore) Create(ctx context.Context, metricName string) (*Metric, error) {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO metric (id, name) VALUES (NULL, ?)")
	if err != nil {
		return nil, err
	}

	res, err := stmt.ExecContext(ctx, metricName)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	metric := &Metric{Name: metricName}

	metric.MetricID, err = res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return metric, nil
}

func (s MetricStore) SelectOne(ctx context.Context, metricName string) (*Metric, error) {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var ret Metric
	err = tx.QueryRowContext(ctx, "SELECT id, name FROM metric WHERE name = ?", metricName).Scan(&ret.MetricID, &ret.Name)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err != nil && err == sql.ErrNoRows {
		return nil, ErrNotFound
	}

	return &ret, tx.Commit()
}
