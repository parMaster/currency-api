package store

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/parmaster/currency-api/internal/data"
	"github.com/stretchr/testify/assert"
)

// Returns system temp dir (i.e. /tmp on Linux, no trailing slash).
// If TEMP_DIR environment variable is set, it is returned instead
func tempDir() string {

	if os.Getenv("TEMP_DIR") != "" {
		return os.Getenv("TEMP_DIR")
	}

	return os.TempDir()
}

func Test_Sqlite_Full(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var err error
	store, err := NewSQLite(ctx, fmt.Sprintf("file:%s/test.db?cache=shared&mode=rwc", tempDir()))
	if err != nil {
		log.Printf("[ERROR] Failed to open SQLite storage: %e", err)
	}

	ratesInit := data.Rates{
		Date: data.Date(time.Now()),
		Base: "USD",
		Rates: map[string]data.FloatRate{
			"UAH": 27.5,
			"EUR": 0.8,
			"RON": 4.5,
		},
	}

	err = store.Write(ratesInit)
	assert.Nil(t, err)

	rates, err := store.Read(time.Now())
	assert.Nil(t, err)
	assert.Equal(t, ratesInit.Base, rates.Base)
	assert.Equal(t, ratesInit.Date.String(), rates.Date.String())
	assert.Equal(t, ratesInit.Rates, rates.Rates)

	// New rates for the previous day
	rates.Date = data.Date(time.Now().AddDate(0, 0, -1))
	rates.Rates["UAH"] = 27.6
	err = store.Write(rates)
	assert.Nil(t, err)

	// Read rates for the previous day
	ratesChanged, err := store.Read(time.Now().AddDate(0, 0, -1))
	assert.Nil(t, err)
	assert.Equal(t, rates.Base, ratesChanged.Base)
	assert.Equal(t, rates.Date.String(), time.Now().AddDate(0, 0, -1).Format("2006-01-02"))
	assert.Equal(t, rates.Rates, ratesChanged.Rates)

	// Rates for the current day should not be changed
	ratesCurrent, err := store.Read(time.Now())
	assert.Nil(t, err)
	assert.Equal(t, ratesInit.Base, ratesCurrent.Base)
	assert.Equal(t, ratesInit.Date.String(), ratesCurrent.Date.String())
	assert.Equal(t, ratesInit.Rates, ratesCurrent.Rates)

	// Change rates for the current day and check
	ratesCurrent.Rates["UAH"] = 27.7
	err = store.Write(ratesCurrent)
	assert.Nil(t, err)

	ratesChanged, err = store.Read(time.Now())
	assert.Nil(t, err)
	assert.Equal(t, ratesCurrent.Base, ratesChanged.Base)
	assert.Equal(t, ratesCurrent.Date.String(), ratesChanged.Date.String())
	assert.Equal(t, ratesCurrent.Rates, ratesChanged.Rates)

	// Teardown
	store.cleanup()
}
