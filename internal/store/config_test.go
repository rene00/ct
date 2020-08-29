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

func TestConfigStoreSelectOne(t *testing.T) {
	dbFile, db, err := testtooling.CreateTmpDB()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	defer os.Remove(dbFile)

	configStore := ConfigStore{db}

	ctx := context.Background()
	if _, err = configStore.SelectOne(ctx, 1, "doesNotExist"); err == nil {
		t.Error("No error raised when calling configStore.SelectOne() with dummy option")
	}

	metricID, err := testtooling.CreateMetric(ctx, db)
	if err != nil {
		t.Error(err)
	}

	ret, err := configStore.SelectOne(ctx, *metricID, "data_type")
	if err != nil {
		t.Error(err)
	}
	if ret != "float" {
		t.Errorf("want float but got %s", ret)
	}
}
