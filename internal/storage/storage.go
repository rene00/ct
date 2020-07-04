package storage

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func UpsertConfig(db *sql.DB, metricID int, opt, val string) error {
	sqlStmt := `
		INSERT INTO config
			(
				metric_id,
				opt,
				val
			)
			VALUES
			(
				?,
				?,
				?
			)
			ON CONFLICT(metric_id, opt)
			DO UPDATE SET val=?
		`
	stmt, err := db.Prepare(sqlStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(metricID, opt, val, val); err != nil {
		return err
	}
	return nil
}
