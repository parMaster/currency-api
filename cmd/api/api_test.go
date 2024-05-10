package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/parmaster/currency-api/internal/data"
	"github.com/parmaster/currency-api/internal/store"
	"github.com/stretchr/testify/assert"
)

func TestServer_Status(t *testing.T) {
	cfg := Options{
		Port:       8080,
		Currencies: "USD,UAH,EUR",
		Interval:   420,
		Debug:      true,
	}
	db, _ := store.NewSQLite(context.Background(), ":memory:")
	s := NewServer(cfg, db, context.Background())
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

func getApiKey() string {
	apiKey := os.Getenv("APIKEY")
	if apiKey == "" {
		key, err := os.ReadFile("../../apikey.env")
		if err != nil {
			return ""
		}
		apiKey = strings.TrimSpace(string(key))
	}
	return apiKey
}

func TestServer_Rates(t *testing.T) {
	db, err := store.NewSQLite(context.Background(), fmt.Sprintf("file:%s/test.db?cache=shared&mode=rwc", os.TempDir()))
	assert.Nil(t, err, "Failed to open SQLite storage: %e", err)

	apiKey := getApiKey()
	if apiKey == "" {
		t.Skip("APIKEY not set")
	}
	s := NewServer(Options{ApiKey: apiKey, Currencies: "UAH,USD,EUR,RON"}, db, context.Background())
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
	s.Rates(w, r, httprouter.Params{httprouter.Param{Key: "date", Value: "2024-04-20"}})
	assert.Equal(t, http.StatusOK, w.Code, "valid date should return 200 OK")

	// Check response
	rates := data.RateResponse{}
	err = json.Unmarshal(w.Body.Bytes(), &rates)
	assert.Nil(t, err, "response should be valid")

	assert.Equal(t, "USD", rates.Base, "base currency should be USD")
	date, err := time.Parse("2006-01-02 15:04:05", "2024-04-20 00:00:00")
	assert.Nil(t, err, "date should be valid")
	assert.Equal(t, data.Date{Time: date}, rates.Date, "date should match")
	assert.NotEmpty(t, rates.Rates, "rates should not be empty")
}

func Test_ServerPairs(t *testing.T) {
	db, err := store.NewSQLite(context.Background(), ":memory:")
	assert.Nil(t, err, "Failed to open SQLite storage: %e", err)
	apiKey := getApiKey()
	if apiKey == "" {
		t.Skip("APIKEY not set")
	}
	s := NewServer(Options{ApiKey: apiKey, Currencies: "USD,UAH,RON,EUR"}, db, context.Background())
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/v1/pair/USD-UAH", nil)

	s.Pair(w, r, httprouter.Params{httprouter.Param{Key: "pair", Value: "USD-UAH"}})
	assert.Equal(t, http.StatusOK, w.Code, "pair should return 200 OK")

	pair := data.PairResponse{}
	json.Unmarshal(w.Body.Bytes(), &pair)

	assert.Equal(t, "USD-UAH", pair.Pair, "base currency should be USD")
	assert.NotEmpty(t, pair.Rate, "rate should not be empty")
}
