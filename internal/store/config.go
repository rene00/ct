package store

import (
	"context"
	"database/sql"
)

// ConfigStorer manages metric config.
type ConfigStorer interface {
	Create(context.Context, int64) error
	SelectOne(context.Context, int64, string) (string, error)
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
