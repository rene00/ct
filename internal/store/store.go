package store

import (
	"database/sql"
	"errors"
)

var (
	// ErrNotFound is row not found.
	ErrNotFound = errors.New("Not Found")
)

// NewStore returns a Store.
func NewStore(DB *sql.DB) *Store {
	return &Store{
		Config: ConfigStore{DB},
		Metric: MetricStore{DB},
		Log:    LogStore{DB},
	}
}

// Store is the main struct which is used to access store methods.
type Store struct {
	Log    LogStorer
	Metric MetricStorer
	Config ConfigStorer
}
