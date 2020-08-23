package store

import (
	"context"
	"database/sql"
)

type ConfigStorer interface {
	Create(context.Context, int64) error
}

type ConfigStore struct {
	DB *sql.DB
}

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
