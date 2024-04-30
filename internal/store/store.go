package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/parmaster/currency-api/internal/data"
)

var (
	ErrNotFound = errors.New("no data found")
)

type Storer interface {
	// Read reads records for the given module from the database
	Read(time.Time) (data.Rates, error)
	// Write writes the data to the database
	Write(data.Rates) error
	// Log requests to the database
	Log(string, string) error
	// ReadLogs return 10 most recent logs from the database
	ReadLogs() ([]string, error)
}

func Load(ctx context.Context, path string, s *Storer) error {
	var err error
	*s, err = NewSQLite(ctx, path)
	if err != nil {
		return fmt.Errorf("failed to init SQLite storage: %e", err)
	}
	return err
}
