package storage

import (
	"fmt"
	"testing"
)

func TestDoMigrateDb(t *testing.T) {
	if err := DoMigrateDb(fmt.Sprintf("sqlite3://:memory:")); err != nil {
		t.Error(err)
	}
}
