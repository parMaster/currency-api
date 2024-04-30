package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/parmaster/currency-api/internal/data"
	"github.com/stretchr/testify/assert"
)

func TestServer_Status(t *testing.T) {
	cfg := Options{
		Port:       8080,
		Currencies: "USD,UAH,EUR",
		Interval:   420,
		Debug:      true,
	}
	s := NewServer(cfg, context.Background())
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/v1/status", nil)

	s.Status(w, r, httprouter.Params{})
	assert.Equal(t, http.StatusOK, w.Code, "status should return 200 OK")

	status := StatusResponse{}
	json.Unmarshal(w.Body.Bytes(), &status)

	assert.Equal(t, "ok", status.Status, "status should be ok")
	assert.NotEmpty(t, status.Version, "version should not be empty")
	assert.Equal(t, cfg, status.Config, "config should match")
}

func TestServer_Rates(t *testing.T) {
	s := NewServer(Options{}, context.Background())
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/v1/rates", nil)

	// Invalid date
	s.Rates(w, r, httprouter.Params{httprouter.Param{Key: "date", Value: "2021-01-"}})
	assert.Contains(t, w.Body.String(), "validation errors", "invalid date should return 400 Bad Request")
	assert.Contains(t, w.Body.String(), "invalid date format, use 2006-01-02", "invalid date should return 400 Bad Request")
	var resp any
	json.Unmarshal(w.Body.Bytes(), &resp)

	// Happy paths
	w = httptest.NewRecorder()
	s.Rates(w, r, httprouter.Params{})
	assert.Equal(t, http.StatusOK, w.Code, "empty date should return 200 OK")

	w = httptest.NewRecorder()
	s.Rates(w, r, httprouter.Params{httprouter.Param{Key: "date", Value: "2021-04-29"}})
	assert.Equal(t, http.StatusOK, w.Code, "valid date should return 200 OK")

	// Check response
	rates := data.RateResponse{}
	json.Unmarshal(w.Body.Bytes(), &rates)

	assert.Equal(t, "USD", rates.Base, "base currency should be USD")
	date, err := time.Parse("2006-01-02 15:04:05", "2021-04-29 00:00:00")
	assert.Nil(t, err, "date should be valid")
	assert.Equal(t, data.Date(date), rates.Date, "base currency should be USD")
	assert.NotEmpty(t, rates.Rates, "rates should not be empty")

}
