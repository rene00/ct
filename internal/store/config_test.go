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

func TestConfigStoreSelectLimit(t *testing.T) {
	dbFile, db, err := testtooling.CreateTmpDB()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	defer os.Remove(dbFile)

	configStore := ConfigStore{db}

	ctx := context.Background()
	ret, err := configStore.SelectLimit(ctx, 0)
	if err != nil {
		t.Error(err)
	}
	if len(ret) != 0 {
		t.Error()
	}

	for i := 0; i < 5; i++ {
		if _, err = testtooling.CreateMetric(ctx, db); err != nil {
			t.Error(err)
		}
	}

	ret, err = configStore.SelectLimit(ctx, 0)
	if err != nil {
		t.Error(err)
	}
	if len(ret) != 5 {
		t.Error()
	}

	ret, err = configStore.SelectLimit(ctx, 1)
	if err != nil {
		t.Error(err)
	}
	if len(ret) != 1 {
		t.Error()
	}
}
