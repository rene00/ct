package store

import (
	"context"
	"ct/internal/testtooling"
	"os"
	"testing"
)

func TestMetricStoreCreate(t *testing.T) {
	dbFile, db, err := testtooling.CreateTmpDB()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	defer os.Remove(dbFile)

	metricStore := MetricStore{db}

	ctx := context.Background()
	metric, err := metricStore.Create(ctx, "test")
	if err != nil {
		t.Error(err)
	}
	if metric.Name != "test" {
		t.Errorf("want metric.Name to equal test but got %s", metric.Name)
	}
}

func TestMetricStoreSelectOne(t *testing.T) {
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

	metricStore := MetricStore{db}

	metric, err := metricStore.SelectOne(ctx, "test")
	if err != nil {
		t.Error(err)
	}
	if metric.MetricID != *metricID {
		t.Errorf("want metric.MetricID to equal %d but got %d", metricID, metric.MetricID)
	}

}

func TestMetricStoreSelectLimit(t *testing.T) {
	dbFile, db, err := testtooling.CreateTmpDB()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	defer os.Remove(dbFile)

	metricStore := MetricStore{db}

	ctx := context.Background()
	ret, err := metricStore.SelectLimit(ctx, 0)
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

	ret, err = metricStore.SelectLimit(ctx, 0)
	if err != nil {
		t.Error(err)
	}
	if len(ret) != 5 {
		t.Error()
	}

	ret, err = metricStore.SelectLimit(ctx, 1)
	if err != nil {
		t.Error(err)
	}
	if len(ret) != 1 {
		t.Error()
	}
}
