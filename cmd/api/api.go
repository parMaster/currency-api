package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
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

func (s *Server) Status(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Write([]byte("Status!\n"))
	w.Write([]byte("Version: " + version + "\n"))
	w.Write([]byte(fmt.Sprintf("Currencies: %v\n", s.cfg.Currencies)))
}

func (s *Server) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Write([]byte("Welcome!\n"))
}

type Rate struct {
	Date  string            `json:"date"`
	Base  string            `json:"base"`
	Rates map[string]string `json:"rates"`
}

// Rates returns exchange rates, stored in the database
// GET /v1/rates[/date]
func (s *Server) Rates(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	validator := validator.New()

	dateStr := ps.ByName("date")
	var date time.Time
	date, err := time.Parse("2006-01-02", dateStr)
	validator.Check(dateStr == "" || err == nil, "date", "invalid date format, use 2006-01-02")

	if !validator.Valid() {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Validation errors!\n"))
		for key, message := range validator.Errors {
			w.Write([]byte(fmt.Sprintf("%s: %s\n", key, message)))
		}
		return
	}

	if !date.IsZero() {
		w.Write([]byte("Rates for date: " + date.Format("2006-01-02") + "\n"))
	}
	w.Write([]byte("Rates!\n"))
}

func (s *Server) Pair(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Write([]byte("Pair!\n"))
	w.Write([]byte(ps.ByName("pair")))
}
