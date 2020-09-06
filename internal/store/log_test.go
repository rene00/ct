package store

import (
	"context"
	"ct/internal/testtooling"
	"os"
	"testing"
	"time"
)

func TestLogStoreUpsert(t *testing.T) {
	dbFile, db, err := testtooling.CreateTmpDB()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	defer os.Remove(dbFile)

	ctx := context.Background()
	metricID, err := testtooling.CreateMetric(ctx, db)
	if err != nil {
		t.Error(err)
	}

	ts, err := time.Parse("2006-01-02", "2020-01-01")
	if err != nil {
		t.Error(err)
	}

	logStore := LogStore{db}

	if err = logStore.Upsert(ctx, &Log{MetricID: *metricID, Value: "1", Timestamp: ts}); err != nil {
		t.Error(err)
	}

	if err = logStore.Upsert(ctx, &Log{MetricID: *metricID, Value: "2", Timestamp: ts}); err != nil {
		t.Error(err)
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Error(err)
	}

	var ret int
	err = tx.QueryRowContext(ctx, "SELECT COUNT(1) FROM log WHERE metric_id = ? AND value = 2 AND timestamp = ?", metricID, ts.Format("2006-01-02")).Scan(&ret)
	if err != nil {
		t.Errorf("Failed to select log: %v", err)
	}
	if ret != 1 {
		t.Errorf("Incorrect number of log entries created: %v", ret)
	}

	if err = tx.Commit(); err != nil {
		t.Error(err)
	}
}

func TestLogStoreCreate(t *testing.T) {
	dbFile, db, err := testtooling.CreateTmpDB()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	defer os.Remove(dbFile)

	ctx := context.Background()
	metricID, err := testtooling.CreateMetric(ctx, db)
	if err != nil {
		t.Error(err)
	}

	ts, err := time.Parse("2006-01-02", "2020-01-01")
	if err != nil {
		t.Error(err)
	}

	logStore := LogStore{db}

	if err = logStore.Create(ctx, &Log{MetricID: *metricID, Value: "1", Timestamp: ts}); err != nil {
		t.Error(err)
	}

	if err = logStore.Create(ctx, &Log{MetricID: *metricID, Value: "1", Timestamp: ts}); err == nil {
		t.Error("No error raised when calling logStore.Create() on duplicate log")
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Error(err)
	}

	var ret int
	err = tx.QueryRowContext(ctx, "SELECT COUNT(1) FROM log WHERE metric_id = ? AND value = 1 AND timestamp = ?", metricID, ts.Format("2006-01-02")).Scan(&ret)
	if err != nil {
		t.Errorf("Failed to select log: %v", err)
	}
	if ret != 1 {
		t.Errorf("Incorrect number of log entries created: %v", ret)
	}

	if err = tx.Commit(); err != nil {
		t.Error(err)
	}
}

func TestLogSelectOne(t *testing.T) {
	dbFile, db, err := testtooling.CreateTmpDB()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	defer os.Remove(dbFile)

	ctx := context.Background()
	metricID, err := testtooling.CreateMetric(ctx, db)
	if err != nil {
		t.Error(err)
	}

	logID, err := testtooling.CreateLog(ctx, db, *metricID)
	if err != nil {
		t.Error(err)
	}

	ts, err := time.Parse("2006-01-02", "2020-01-01")
	if err != nil {
		t.Error(err)
	}

	logStore := LogStore{db}
	log, err := logStore.SelectOne(ctx, *metricID, ts)
	if err != nil {
		t.Error(err)
	}
	if log.LogID != *logID {
		t.Errorf("want log.LogID to equal %d but got %d", logID, log.LogID)
	}
}

func TestLogStoreSelectLimit(t *testing.T) {
	dbFile, db, err := testtooling.CreateTmpDB()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	defer os.Remove(dbFile)

	logStore := LogStore{db}

	ctx := context.Background()
	ret, err := logStore.SelectLimit(ctx, 0)
	if err != nil {
		t.Error(err)
	}
	if len(ret) != 0 {
		t.Error()
	}

	for i := 0; i < 5; i++ {
		metricID, err := testtooling.CreateMetric(ctx, db)
		if err != nil {
			t.Error(err)
		}
		if _, err := testtooling.CreateLog(ctx, db, *metricID); err != nil {
			t.Error(err)
		}
	}

	ret, err = logStore.SelectLimit(ctx, 0)
	if err != nil {
		t.Error(err)
	}
	if len(ret) != 5 {
		t.Error()
	}

	ret, err = logStore.SelectLimit(ctx, 1)
	if err != nil {
		t.Error(err)
	}
	if len(ret) != 1 {
		t.Error()
	}
}
