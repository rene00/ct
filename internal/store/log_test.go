package store

import (
	"context"
	"ct/internal/storage"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func testSetup() (string, *sql.DB, error) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "")
	if err != nil {
		return "", nil, err
	}

	dbFile := tmpFile.Name()
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return "", nil, err
	}
	if err := storage.DoMigrateDb(fmt.Sprintf("sqlite3://%s", dbFile)); err != nil {
		return "", nil, err
	}
	return dbFile, db, nil
}

func TestLogStoreCreate(t *testing.T) {
	dbFile, db, err := testSetup()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	defer os.Remove(dbFile)

	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Error(err)
	}

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO metric (id, name) VALUES (NULL, ?)")
	if err != nil {
		t.Error(err)
	}

	res, err := stmt.ExecContext(ctx, "test")
	if err != nil {
		t.Error(err)
	}

	metricID, err := res.LastInsertId()
	if err != nil {
		t.Error(err)
	}

	stmt, err = tx.PrepareContext(ctx, "INSERT INTO config (metric_id, opt, val) VALUES (?, ?, ?)")
	if err != nil {
		t.Error(err)
	}

	if _, err = stmt.ExecContext(ctx, metricID, "data_type", "float"); err != nil {
		t.Error(err)
	}

	if err = tx.Commit(); err != nil {
		t.Error(err)
	}

	ts, err := time.Parse("2006-01-02", "2020-01-01")
	if err != nil {
		t.Error(err)
	}

	logStore := LogStore{db}

	if err = logStore.Create(ctx, &Log{MetricID: metricID, Value: "1", Timestamp: ts}); err != nil {
		t.Error(err)
	}

	if err = logStore.Create(ctx, &Log{MetricID: metricID, Value: "1", Timestamp: ts}); err == nil {
		t.Error("Want err but err is nil when claling logstore.Create()")
	}

}
