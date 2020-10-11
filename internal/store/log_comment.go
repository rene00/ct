package store

import (
	"context"
	"database/sql"
)

// LogCommentStorer manages log comments.
type LogCommentStorer interface {
	Upsert(context.Context, *Log, string) error
	SelectOne(context.Context, int64) (*LogComment, error)
}

// LogCommentStore manages log comments.
type LogCommentStore struct {
	DB *sql.DB
}

// Upsert creates a new log comment.
func (s LogCommentStore) Upsert(ctx context.Context, o *Log, comment string) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO log_comment (log_id, comment) VALUES (?, ?) ON CONFLICT(log_id) DO UPDATE SET comment = ?")
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, o.LogID, comment, comment)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s LogCommentStore) SelectOne(ctx context.Context, logID int64) (*LogComment, error) {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	ret := LogComment{LogID: logID}
	err = tx.QueryRowContext(ctx, "SELECT comment FROM log_comment WHERE log_id = ?", logID).Scan(&ret.Comment)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err != nil && err == sql.ErrNoRows {
		return nil, ErrNotFound
	}

	return &ret, tx.Commit()
}
