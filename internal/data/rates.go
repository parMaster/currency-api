package data

import (
	"fmt"
	"strings"
	"time"
)

type Date time.Time

func (r Date) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", time.Time(r).Format("2006-01-02 15:04:05"))), nil
}

func (r *Date) UnmarshalJSON(data []byte) error {
	dataStr := strings.Trim(string(data), "\"")
	t, err := time.Parse("2006-01-02 00:00:00", dataStr)
	if err != nil {
		t = time.Time{}
	}
	*r = Date(t)
	return nil
}

type RateResponse struct {
	Date  Date               `json:"date"`
	Base  string             `json:"base"`
	Rates map[string]float32 `json:"rates"`
}
