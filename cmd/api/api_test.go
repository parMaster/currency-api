package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestServer_Rates(t *testing.T) {
	s := NewServer(Options{}, context.Background())
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/v1/rates", nil)

	// Happy paths
	s.Rates(w, r, httprouter.Params{})
	assert.Equal(t, http.StatusOK, w.Code, "empty date should return 200 OK")

	s.Rates(w, r, httprouter.Params{httprouter.Param{Key: "date", Value: "2021-01-01"}})
	assert.Equal(t, http.StatusOK, w.Code, "valid date should return 200 OK")

	// Invalid date
	s.Rates(w, r, httprouter.Params{httprouter.Param{Key: "date", Value: "2021-01-"}})
	assert.Contains(t, w.Body.String(), "Validation errors!", "invalid date should return 400 Bad Request")
	assert.Contains(t, w.Body.String(), "date: invalid date format, use 2006-01-02", "invalid date should return 400 Bad Request")
}
