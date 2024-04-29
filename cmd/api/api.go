package main

import (
	"fmt"
	"log"
	"net/http"
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
		log.Printf("[ERROR] failed to write response, %+v", err)
		http.Error(w, "failed to write response", http.StatusInternalServerError)
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
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Validation errors!\n"))
		for key, message := range validator.Errors {
			w.Write([]byte(fmt.Sprintf("%s: %s\n", key, message)))
		}
		return
	}

	rates := data.RateResponse{
		Date:  data.Date(date),
		Base:  "USD",
		Rates: map[string]float32{"UAH": 39.5, "EUR": 0.85, "RON": 4.8},
	}

	err = s.writeJSON(w, http.StatusOK, rates, nil)
	if err != nil {
		log.Printf("[ERROR] failed to write response, %+v", err)
		http.Error(w, "failed to write response", http.StatusInternalServerError)
	}

}

func (s *Server) Pair(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Write([]byte("Pair!\n"))
	w.Write([]byte(ps.ByName("pair")))
}
