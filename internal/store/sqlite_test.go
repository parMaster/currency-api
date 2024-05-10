package store

import (
	"context"
	"testing"
	"time"

	"github.com/parmaster/currency-api/internal/data"
	"github.com/stretchr/testify/assert"
)

// Test_Sqlite_InitData tests the initialization of the SQLite storage
// with the initial data
func Test_Sqlite_InitData(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var err error
	store, err := NewSQLite(ctx, ":memory:")
	assert.Nil(t, err, "Failed to open SQLite storage: %e", err)

	// Check if there are 6 initial rows in the rates table
	cnt := 0
	row := store.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM rates")
	err = row.Scan(&cnt)
	assert.Nil(t, err)
	assert.Equal(t, 6, cnt)

	// Tear down the storage
	store.cleanup()
}
func Test_Sqlite_Full(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var err error
	store, err := NewSQLite(ctx, ":memory:")
	assert.Nil(t, err, "Failed to open SQLite storage: %e", err)

	ratesInit := data.Rates{
		Date: data.Date{Time: time.Now()},
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
	rates.Date = data.Date{Time: time.Now().AddDate(0, 0, -1)}
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
