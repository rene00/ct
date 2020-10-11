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

	t.Run("Create", func(t *testing.T) {
		err = logCommentStore.Create(ctx, &Log{LogID: *logID}, "test")
		assert.Nil(t, err)
	})
}
