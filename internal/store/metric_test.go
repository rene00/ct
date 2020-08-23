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
