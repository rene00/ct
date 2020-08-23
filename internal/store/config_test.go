package store

import (
	"context"
	"ct/internal/testtooling"
	"os"
	"testing"
)

func TestConfigStoreCreate(t *testing.T) {
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

	configStore := ConfigStore{db}
	if err = configStore.Create(ctx, *metricID); err != nil {
		t.Error(err)
	}
}
