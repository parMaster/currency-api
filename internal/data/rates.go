package data

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// RateResponse is a response from the API
type RateResponse struct {
	Date  Date               `json:"date"`
	Base  string             `json:"base"`
	Rates map[string]float32 `json:"rates"`
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
