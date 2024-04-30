package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/parmaster/currency-api/internal/data"
	"github.com/parmaster/currency-api/internal/validator"
)

func (s *Server) router() http.Handler {

	router := httprouter.New()
	router.GET("/", s.Index)
	router.GET("/v1/status", s.Status)

	router.GET("/v1/rates", s.Rates)
	// date format: 2006-02-01
	router.GET("/v1/rates/:date", s.Rates)

	// pair format: USD-UAH (1 USD = x UAH)
	router.GET("/v1/pair/:pair", s.Pair)

	return router
}

type StatusResponse struct {
	Status  string  `json:"status"`
	Version string  `json:"version"`
	Config  Options `json:"config"`
}

func (s *Server) Status(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	status := StatusResponse{
		Status:  "ok",
		Version: version,
		Config:  s.cfg,
	}

	err := s.writeJSON(w, http.StatusOK, status, nil)
	if err != nil {
		http.Error(w, "failed to write response: "+err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Write([]byte("Welcome!\n"))
}

// Rates returns exchange rates, stored in the database
// GET /v1/rates[/date]
func (s *Server) Rates(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	dateStr := ps.ByName("date")
	date, err := time.Parse("2006-01-02", dateStr)

	validator := validator.New()
	validator.Check(dateStr == "" || err == nil, "date", "invalid date format, use 2006-01-02")

	if !validator.Valid() {
		errors := struct {
			Error   string            `json:"error"`
			Message map[string]string `json:"message"`
		}{Error: "validation errors", Message: validator.Errors}
		s.writeJSON(w, http.StatusBadRequest, errors, nil)
		return
	}

	// if date is empty, use today
	if date.IsZero() {
		date = time.Now()
	}
	// prefer data from the database
	// if not found, use the API and store in the database
	// historical data can be unavailable

	// TODO: log requests in database

	rates := data.RateResponse{
		Date:  data.Date(date),
		Base:  "USD",
		Rates: map[string]float32{"UAH": 39.5, "EUR": 0.85, "RON": 4.8},
	}

	err = s.writeJSON(w, http.StatusOK, rates, nil)
	if err != nil {
		http.Error(w, "failed to write response: "+err.Error(), http.StatusInternalServerError)
	}

}

func (s *Server) Pair(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pair := strings.Split(ps.ByName("pair"), "-")

	valid := validator.New()
	valid.Check(len(pair) == 2, "pair", "invalid pair format, use USD-UAH")
	permittedValue := validator.PermittedValue(pair[0], strings.Split(s.cfg.Currencies, ",")...)
	valid.Check(permittedValue, "pair", "invalid currency, use these: "+s.cfg.Currencies)

	if !valid.Valid() {
		errors := struct {
			Error   string            `json:"error"`
			Message map[string]string `json:"message"`
		}{Error: "validation errors", Message: valid.Errors}
		s.writeJSON(w, http.StatusBadRequest, errors, nil)
		return
	}

	err := s.writeJSON(w, http.StatusOK, pair, nil)
	if err != nil {
		http.Error(w, "failed to write response: "+err.Error(), http.StatusInternalServerError)
	}

}
