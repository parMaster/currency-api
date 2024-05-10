package client

import (
	"os"
	"testing"
	"time"

	"github.com/parmaster/currency-api/internal/data"
	"github.com/stretchr/testify/assert"
)

func Test_ParseJSONResponse(t *testing.T) {

	// Happy path
	mock := []byte(`{"date":"2024-04-29 12:34:56+00","base":"USD","rates":{"RON":"4.648225","EUR":"0.9340084415835656","USD":"1.0","UAH":"39.65869122539661"}}`)
	client := New(os.Getenv("APIKEY"))
	rates, err := client.parseResponse(mock)

	assert.Nil(t, err)
	assert.NotEmpty(t, rates)
	assert.NotEmpty(t, rates.Date)
	date, err := time.Parse("2006-01-02 15:04:05+00", "2024-04-29 12:34:56+00")
	assert.Nil(t, err)
	assert.Equal(t, data.Date{Time: date}, rates.Date)
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

	// Invalid JSON
	mock = []byte(`{"date":"2024-04-29 00:00:00+00","base":"USD","rates":{"RON":"4.648225","EUR":"0.9340084415835656","USD":"1.0","UAH":"39.65869122539661"`)
	rates, err = client.parseResponse(mock)
	assert.NotNil(t, err)
	assert.Empty(t, rates)

}

func Test_RequestLatestRates(t *testing.T) {
	if os.Getenv("APIKEY") == "" {
		t.Skip("APIKEY not set")
	}
	client := New(os.Getenv("APIKEY"))
	rates, err := client.GetLatest("USD,UAH,EUR,RON")
	assert.Nil(t, err)
	assert.NotEmpty(t, rates)
	assert.NotEmpty(t, rates.Date)
	assert.NotEmpty(t, rates.Base)
	assert.NotEmpty(t, rates.Rates)
}

func Test_RequestHistoricalRates(t *testing.T) {
	t.Skip("Skipping historical rates test - requires a paid plan")

	if os.Getenv("APIKEY") == "" {
		t.Skip("APIKEY not set")
	}
	client := New(os.Getenv("APIKEY"))
	rates, err := client.GetHistorical("USD,UAH,EUR,RON", time.Now().AddDate(0, 0, -1))
	assert.Nil(t, err)
	assert.NotEmpty(t, rates)
	assert.NotEmpty(t, rates.Date)
	assert.NotEmpty(t, rates.Base)
	assert.NotEmpty(t, rates.Rates)
}
