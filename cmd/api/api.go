package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/parmaster/currency-api/internal/client"
	"github.com/parmaster/currency-api/internal/data"
	"github.com/parmaster/currency-api/internal/store"
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
	Status  string   `json:"status"`
	Version string   `json:"version"`
	Config  Options  `json:"config"`
	Logs    []string `json:"logs"`
}

func (s *Server) Status(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	status := StatusResponse{
		Status:  "ok",
		Version: version,
		Config:  s.cfg,
	}
	var err error
	status.Logs, err = s.db.ReadLogs()
	if err != nil {
		http.Error(w, "failed to read logs: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = s.writeJSON(w, http.StatusOK, status, nil)
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

	err = s.db.Log("rates", fmt.Sprintf("date: %s", dateStr))
	if err != nil {
		log.Printf("[ERROR] failed to log request: %v", err)
	}

	rates, err := s.GetUpdateRates(date)
	if err != nil && err != ErrNoContent {
		date = time.Now().AddDate(0, 0, -1)
		rates, err = s.GetUpdateRates(date)
	}
	if err != nil && errors.Is(err, ErrNoContent) {
		http.Error(w, "no rates available", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "failed to get rates: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = s.writeJSON(w, http.StatusOK, rates, nil)
	if err != nil {
		http.Error(w, "failed to write response: "+err.Error(), http.StatusInternalServerError)
	}

}

var ErrNoContent = errors.New("no rates available")

func (s *Server) GetUpdateRates(date time.Time) (data.Rates, error) {
	// prefer data from the database
	rates, err := s.db.Read(date)
	if err == store.ErrNotFound {
		// if not found, use the API
		client := client.New(s.cfg.ApiKey)
		if date.IsZero() {
			rates, err = client.GetLatest(s.cfg.Currencies)
		} else {
			rates, err = client.GetHistorical(s.cfg.Currencies, date)
		}
		if err == nil {
			// and store in the database
			err = s.db.Write(rates)
			if err != nil {
				log.Printf("[ERROR] failed to write rates: %v", err)
			}
		} else {
			log.Printf("[ERROR] failed to get rates: %v", err)
			return data.Rates{}, err
		}
	} else if err != nil {
		log.Printf("[ERROR] failed to read rates: %v", err)
		return data.Rates{}, err
	}
	// historical data can be unavailable - /historical endpoint requires a paid subscription

	if len(rates.Rates) == 0 {
		return data.Rates{}, ErrNoContent
	}
	return rates, nil
}

func (s *Server) Pair(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pair := strings.Split(ps.ByName("pair"), "-")

	valid := validator.New()

	// basic validation
	validPair := len(pair) == 2 && len(pair[0]) == 3 && len(pair[1]) == 3
	valid.Check(validPair, "pair", "invalid pair format, use USD-UAH")

	// check if the pair is in the list of supported currencies
	permittedValue := validator.PermittedValue(pair[0], strings.Split(s.cfg.Currencies, ",")...) &&
		validator.PermittedValue(pair[1], strings.Split(s.cfg.Currencies, ",")...)
	valid.Check(permittedValue, "pair", "invalid currency, use these: "+s.cfg.Currencies)

	if !valid.Valid() {
		errors := struct {
			Error   string            `json:"error"`
			Message map[string]string `json:"message"`
		}{Error: "validation errors", Message: valid.Errors}
		s.writeJSON(w, http.StatusBadRequest, errors, nil)
		return
	}

	err := s.db.Log("pair", fmt.Sprintf("pair: %s", strings.Join(pair, "-")))
	if err != nil {
		log.Printf("[ERROR] failed to log request: %v", err)
	}

	date := time.Time{}
	rates, err := s.GetUpdateRates(date)
	if err != nil && err != ErrNoContent {
		date = time.Now().AddDate(0, 0, -1)
		rates, err = s.GetUpdateRates(date)
	}
	if err != nil && errors.Is(err, ErrNoContent) {
		http.Error(w, "no rates available", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "failed to get rates: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rate := data.FloatRate(rates.Rates[pair[1]] / rates.Rates[pair[0]])

	pairResponse := data.PairResponse{
		Date: rates.Date.String(),
		Pair: strings.Join(pair, "-"),
		Rate: rate,
	}

	err = s.writeJSON(w, http.StatusOK, pairResponse, nil)
	if err != nil {
		http.Error(w, "failed to write response: "+err.Error(), http.StatusInternalServerError)
	}

}
