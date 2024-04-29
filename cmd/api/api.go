package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
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

func (s *Server) Rates(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Write([]byte("Rates!\n"))
}

func (s *Server) Pair(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Write([]byte("Pair!\n"))
	w.Write([]byte(ps.ByName("pair")))
}
