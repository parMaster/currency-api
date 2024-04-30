package client

import (
	"testing"
	"time"

	"github.com/parmaster/currency-api/internal/data"
	"github.com/stretchr/testify/assert"
)

func Test_GetRates(t *testing.T) {

	mock := []byte(`{"date":"2024-04-29 00:00:00+00","base":"USD","rates":{"RON":"4.648225","EUR":"0.9340084415835656","USD":"1.0","UAH":"39.65869122539661"}}`)
	client := New("test")
	rates, err := client.GetRates([]string{"USD", "UAH", "EUR", "RON"}, mock)

	assert.Nil(t, err)
	assert.NotEmpty(t, rates)
	assert.NotEmpty(t, rates.Date)
	date, err := time.Parse("2006-01-02 15:04:05+00", "2024-04-29 00:00:00+00")
	assert.Nil(t, err)
	assert.Equal(t, data.Date(date), rates.Date)
	assert.NotEmpty(t, rates.Base)
	assert.Equal(t, "USD", rates.Base)
	assert.NotEmpty(t, rates.Rates)
	assert.Equal(t, 4, len(rates.Rates))

	_, ok := rates.Rates["USD"]
	assert.True(t, ok)
	_, ok = rates.Rates["UAH"]
	assert.True(t, ok)
	_, ok = rates.Rates["EUR"]
	assert.True(t, ok)
	_, ok = rates.Rates["RON"]
	assert.True(t, ok)

	assert.Equal(t, float64(1.0), float64(rates.Rates["USD"]))
	assert.Equal(t, float64(39.65869122539661), float64(rates.Rates["UAH"]))
	assert.Equal(t, float64(0.9340084415835656), float64(rates.Rates["EUR"]))
	assert.Equal(t, float64(4.648225), float64(rates.Rates["RON"]))
}
