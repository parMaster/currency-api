package main

import (
	"errors"
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
		// if not found, use the API and store in the database
		client := client.New(s.cfg.ApiKey)
		rates, err = client.GetRates(strings.Split(s.cfg.Currencies, ","), date, []byte{})
		if err == nil {
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
	// historical data can be unavailable

	if len(rates.Rates) == 0 {
		return data.Rates{}, ErrNoContent
	}
	return rates, nil
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
