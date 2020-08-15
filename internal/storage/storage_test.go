package storage

import (
	"ct/internal/model"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
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
	if err := DoMigrateDb(fmt.Sprintf("sqlite3://%s", dbFile)); err != nil {
		return "", nil, err
	}
	return dbFile, db, nil
}

func TestCreateMetric(t *testing.T) {
	dbFile, db, err := testSetup()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	defer os.Remove(dbFile)

	m := model.Metric{Name: "test"}
	mID, err := CreateMetric(db, m)
	if err != nil {
		t.Error(err)
	}

	sqlStmt := `SELECT id, name FROM metric WHERE name = ?`
	stmt, err := db.Prepare(sqlStmt)
	if err != nil {
		t.Error(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query("test")
	if err != nil {
		t.Error(err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		count++
	}

	if count != 1 {
		t.Errorf("Want 1 row but got %d", count)
	}

	sqlStmt = `SELECT val FROM config WHERE metric_id = ? AND opt = "data_type"`
	stmt, err = db.Prepare(sqlStmt)
	if err != nil {
		t.Error(err)
	}
	defer stmt.Close()

	var dataType string
	if err = stmt.QueryRow(&mID).Scan(&dataType); err != nil {
		t.Error(err)
	}
	if dataType != "float" {
		t.Errorf("Want float but got %s", dataType)
	}

}
