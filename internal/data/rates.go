package data

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// PairResponse is a response from the API
type PairResponse struct {
	Date string    `json:"date"`
	Pair string    `json:"pair"`
	Rate FloatRate `json:"rate"`
}

// RateResponse is a response from the API
type RateResponse struct {
	Date  Date                 `json:"date"`
	Base  string               `json:"base"`
	Rates map[string]FloatRate `json:"rates"`
}

type Date time.Time

func (r Date) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", time.Time(r).Format("2006-01-02 15:04:05"))), nil
}

func (r *Date) UnmarshalJSON(data []byte) error {
	dataStr := strings.Trim(string(data), "\"")
	// ditch the timezone
	dataStr = strings.Split(dataStr, "+")[0]

	t, err := time.Parse("2006-01-02 15:04:05", dataStr)
	if err != nil {
		t = time.Time{}
	}
	*r = Date(t)
	return nil
}

func (r *Date) ParseDate(date string) error {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return err
	}
	*r = Date(t)
	return nil
}

func (r Date) String() string {
	return time.Time(r).Format("2006-01-02")
}

// Rates is a client response
type Rates struct {
	Date  Date                 `json:"date"`
	Base  string               `json:"base"`
	Rates map[string]FloatRate `json:"rates"`
}

type FloatRate float64

func (r *FloatRate) UnmarshalJSON(data []byte) error {
	dataStr := strings.Trim(string(data), "\"")
	t, err := strconv.ParseFloat(dataStr, 64)
	if err != nil {
		t = 0
	}
	*r = FloatRate(t)
	return nil
}
