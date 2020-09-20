package store

import (
	"context"
	"database/sql"
)

// MetricStorer manages metrics.
type MetricStorer interface {
	Create(context.Context, string) (*Metric, error)
	SelectOne(context.Context, string) (*Metric, error)
	SelectLimit(context.Context, int64) ([]Metric, error)
	Delete(context.Context, int64) error
}

// MetricStore manages metrics.
type MetricStore struct {
	DB *sql.DB
}

// Delete deletes a metric by id.
func (s MetricStore) Delete(ctx context.Context, metricID int64) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, "DELETE FROM metric WHERE id = ?")
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, metricID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// Create creates a new metric.
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

// SelectLimit returns all metric rows up to limit. A limit of 0 is no limit.
func (s MetricStore) SelectLimit(ctx context.Context, limit int64) ([]Metric, error) {

	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var ret []Metric
	var rows *sql.Rows

	if limit == 0 {
		rows, err = tx.QueryContext(ctx, "SELECT id, name FROM metric ORDER BY name")
		if err != nil {
			return nil, err
		}
	} else {
		rows, err = tx.QueryContext(ctx, "SELECT id, name FROM metric ORDER BY name LIMIT ?", limit)
		if err != nil {
			return nil, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		var o Metric
		if err = rows.Scan(&o.MetricID, &o.Name); err != nil {
			return nil, err
		}
		ret = append(ret, o)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ret, tx.Commit()
}

// SelectOne returns a single metric.
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
