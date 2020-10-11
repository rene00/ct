package store

import (
	"context"
	"ct/internal/testtooling"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogComment(t *testing.T) {
	dbFile, db, err := testtooling.CreateTmpDB()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	defer os.Remove(dbFile)

	ctx := context.Background()
	metricID, err := testtooling.CreateMetric(ctx, db)
	assert.Nil(t, err)

	logID, err := testtooling.CreateLog(ctx, db, *metricID)
	assert.Nil(t, err)

	logCommentStore := LogCommentStore{db}

	t.Run("Upsert", func(t *testing.T) {
		err = logCommentStore.Upsert(ctx, &Log{LogID: *logID}, "test")
		assert.Nil(t, err)

		tx, err := db.BeginTx(ctx, nil)
		assert.Nil(t, err)

		var comment string
		err = tx.QueryRowContext(ctx, "SELECT comment FROM log_comment WHERE log_id = ?", *logID).Scan(&comment)
		assert.Nil(t, err)
		assert.Equal(t, "test", comment)

		err = tx.Commit()
		assert.Nil(t, err)

		tx, err = db.BeginTx(ctx, nil)
		assert.Nil(t, err)

		err = logCommentStore.Upsert(ctx, &Log{LogID: *logID}, "more test")
		assert.Nil(t, err)

		err = tx.QueryRowContext(ctx, "SELECT comment FROM log_comment WHERE log_id = ?", *logID).Scan(&comment)
		assert.Nil(t, err)
		assert.Equal(t, "more test", comment)

		err = tx.Commit()
		assert.Nil(t, err)
	})
}
