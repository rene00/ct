package store

import (
	"context"
	"database/sql"
)

// ConfigStorer manages metric config.
type ConfigStorer interface {
	Create(context.Context, int64) error
	SelectOne(context.Context, int64, string) (string, error)
	Upsert(context.Context, *Config) error
	SelectLimit(context.Context, int64) ([]Config, error)
}

// ConfigStore manages metric config.
type ConfigStore struct {
	DB *sql.DB
}

// Create creates a new config item.
func (s ConfigStore) Create(ctx context.Context, metricID int64) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO config (metric_id, opt, val) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}

	// FIXME: Default data_type to float for now but support WithDataType() in the future.
	_, err = stmt.ExecContext(ctx, metricID, "data_type", "float")
	if err != nil {
		return err
	}

	return tx.Commit()
}

// SelectOne selects a single config value from a metric ID and config option.
func (s ConfigStore) SelectOne(ctx context.Context, metricID int64, configOpt string) (string, error) {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	var ret string
	err = tx.QueryRowContext(ctx, "SELECT val FROM config WHERE metric_id = ? AND opt = ?", metricID, configOpt).Scan(&ret)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}
	if err != nil && err == sql.ErrNoRows {
		return "", ErrNotFound
	}

	return ret, tx.Commit()
}

// Upsert inserts a new config or updates an existing config if it exists.
func (s ConfigStore) Upsert(ctx context.Context, o *Config) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO config (metric_id, opt, val) VALUES (?, ?, ?) ON CONFLICT(metric_id, opt) DO UPDATE SET val = ?")
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, o.MetricID, o.Opt, o.Val, o.Val)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// SelectLimit returns all config rows up to limit. A limit of 0 is no limit.
func (s ConfigStore) SelectLimit(ctx context.Context, limit int64) ([]Config, error) {

	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var ret []Config
	var rows *sql.Rows

	if limit == 0 {
		rows, err = tx.QueryContext(ctx, "SELECT metric_id, opt, val FROM config")
		if err != nil {
			return nil, err
		}
	} else {
		rows, err = tx.QueryContext(ctx, "SELECT metric_id, opt, val FROM config LIMIT ?", limit)
		if err != nil {
			return nil, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		var o Config
		if err = rows.Scan(&o.MetricID, &o.Opt, &o.Val); err != nil {
			return nil, err
		}
		ret = append(ret, o)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ret, tx.Commit()
}
