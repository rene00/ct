// Tools for testing.

package testtooling

import (
	"context"
	"ct/internal/storage"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

// CreateLog creates a single log.
func CreateLog(ctx context.Context, db *sql.DB, metricID int64) (*int64, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO log (id, timestamp, metric_id, value) VALUES (NULL, ?, ?, ?)")
	if err != nil {
		return nil, err
	}

	value := "1"
	ts, err := time.Parse("2006-01-02", "2020-01-01")
	if err != nil {
		return nil, err
	}

	res, err := stmt.ExecContext(ctx, ts.Format("2006-01-02"), metricID, value)
	if err != nil {
		return nil, err
	}

	logID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &logID, tx.Commit()
}

// CreateMetric creates a single metric with sane defaults.
func CreateMetric(ctx context.Context, db *sql.DB) (*int64, error) {

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO metric (id, name) VALUES (NULL, ?)")
	if err != nil {
		return nil, err
	}

	res, err := stmt.ExecContext(ctx, "test")
	if err != nil {
		return nil, err
	}

	metricID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	stmt, err = tx.PrepareContext(ctx, "INSERT INTO config (metric_id, opt, val) VALUES (?, ?, ?)")
	if err != nil {
		return nil, err
	}

	if _, err = stmt.ExecContext(ctx, metricID, "data_type", "float"); err != nil {
		return nil, err
	}

	return &metricID, tx.Commit()
}

// CreateTmpDB creates a tmp sqlite db.
func CreateTmpDB() (string, *sql.DB, error) {
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
