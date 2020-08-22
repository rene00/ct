package store

import (
	"database/sql"
	"errors"
)

var (
	// ErrNotFound is row not found.
	ErrNotFound = errors.New("Not Found")
)

func NewStore(DB *sql.DB) *Store {
	return &Store{
		Config: ConfigStore{DB},
		Metric: MetricStore{DB},
		Log:    LogStore{DB},
	}
}

type Store struct {
	Log    LogStorer
	Metric MetricStorer
	Config ConfigStorer
}
