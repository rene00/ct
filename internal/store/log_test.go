package store

import (
	"context"
	"ct/internal/testtooling"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLogStore(t *testing.T) {

	ts, err := time.Parse("2006-01-02", "2020-01-01")
	assert.Nil(t, err)

	t.Run("Create", func(t *testing.T) {
		dbFile, db, err := testtooling.CreateTmpDB()
		if err != nil {
			t.Error(err)
		}
		defer db.Close()
		defer os.Remove(dbFile)

		ctx := context.Background()
		metricID, err := testtooling.CreateMetric(ctx, db)
		assert.Nil(t, err)

		logStore := LogStore{db}
		_, err = logStore.Create(ctx, &Log{MetricID: *metricID, Value: "1", Timestamp: ts})
		assert.Nil(t, err)

		tx, err := db.BeginTx(ctx, nil)
		assert.Nil(t, err)

		var ret int
		err = tx.QueryRowContext(ctx, "SELECT COUNT(1) FROM log WHERE metric_id = ? AND value = 1 AND timestamp = ?", metricID, ts.Format("2006-01-02")).Scan(&ret)
		assert.Nil(t, err)
		assert.Equal(t, 1, ret)

		// Create a duplicate log which will return err.
		_, err = logStore.Create(ctx, &Log{MetricID: *metricID, Value: "1", Timestamp: ts})
		assert.Error(t, err)
	})

	t.Run("Update", func(t *testing.T) {
		dbFile, db, err := testtooling.CreateTmpDB()
		if err != nil {
			t.Error(err)
		}
		defer db.Close()
		defer os.Remove(dbFile)

		ctx := context.Background()
		metricID, err := testtooling.CreateMetric(ctx, db)
		assert.Nil(t, err)

		logStore := LogStore{db}
		_, err = logStore.Upsert(ctx, &Log{MetricID: *metricID, Value: "1", Timestamp: ts})
		assert.Nil(t, err)
		_, err = logStore.Upsert(ctx, &Log{MetricID: *metricID, Value: "2", Timestamp: ts})
		assert.Nil(t, err)
		tx, err := db.BeginTx(ctx, nil)
		assert.Nil(t, err)

		var ret int
		err = tx.QueryRowContext(ctx, "SELECT COUNT(1) FROM log WHERE metric_id = ? AND value = 2 AND timestamp = ?", metricID, ts.Format("2006-01-02")).Scan(&ret)
		assert.Nil(t, err)
		assert.Equal(t, 1, ret)
	})

	t.Run("SelectOne", func(t *testing.T) {
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

		logStore := LogStore{db}
		log, err := logStore.SelectOne(ctx, *metricID, ts)
		assert.Nil(t, err)

		assert.Equal(t, log.LogID, *logID)
	})

	t.Run("SelectLimit", func(t *testing.T) {
		dbFile, db, err := testtooling.CreateTmpDB()
		assert.Nil(t, err)
		defer db.Close()
		defer os.Remove(dbFile)

		logStore := LogStore{db}

		ctx := context.Background()
		ret, err := logStore.SelectLimit(ctx, 0)
		assert.Nil(t, err)
		assert.Len(t, ret, 0)

		for i := 0; i < 5; i++ {
			metricID, err := testtooling.CreateMetric(ctx, db)
			assert.Nil(t, err)
			_, err = testtooling.CreateLog(ctx, db, *metricID)
			assert.Nil(t, err)
		}

		ret, err = logStore.SelectLimit(ctx, 0)
		assert.Nil(t, err)
		assert.Len(t, ret, 5)

		ret, err = logStore.SelectLimit(ctx, 1)
		assert.Nil(t, err)
		assert.Len(t, ret, 1)
	})

	t.Run("SelectWithTimestamp", func(t *testing.T) {
		dbFile, db, err := testtooling.CreateTmpDB()
		assert.NoError(t, err)
		defer db.Close()
		defer os.Remove(dbFile)

		logStore := LogStore{db}

		ctx := context.Background()
		ret, err := logStore.SelectWithTimestamp(ctx, 0, time.Now())
		assert.NoError(t, err)
		assert.Equal(t, len(ret), 0)
	})
}
