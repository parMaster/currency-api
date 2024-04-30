package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/parmaster/currency-api/internal/data"
)

type Client struct {
	ApiUrl map[string]string
	ApiKey string
}

func New(apiKey string) *Client {
	return &Client{
		ApiUrl: map[string]string{
			"latest":     "https://api.currencyfreaks.com/v2.0/rates/latest",
			"historical": "https://api.currencyfreaks.com/v2.0/rates/historical",
		},
		ApiKey: apiKey,
	}
}

func (c *Client) request(endpoint string, parameters map[string]string) ([]byte, error) {
	params := url.Values{}
	params.Add(`apikey`, c.ApiKey)
	for k, v := range parameters {
		params.Add(k, v)
	}
	log.Printf("[DEBUG] CF request: %s?%s", c.ApiUrl[endpoint], params.Encode())

	response, err := http.Get(fmt.Sprintf("%s?%s", c.ApiUrl[endpoint], params.Encode()))
	if err != nil {
		return []byte{}, err
	}
	defer response.Body.Close()

	// dump response
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return []byte{}, err
	}
	log.Printf("[DEBUG] CF Api response: %s", string(body))

	return body, nil
}

func (c *Client) parseResponse(response []byte) (data.Rates, error) {
	rates := data.Rates{}
	if err := json.Unmarshal(response, &rates); err != nil {
		return data.Rates{}, err
	}

	return rates, nil
}

func (c *Client) GetLatest(symbols string) (data.Rates, error) {
	response, err := c.request(
		"latest",
		map[string]string{
			`symbols`: symbols,
		},
	)
	if err != nil {
		return data.Rates{}, err
	}

	return c.parseResponse(response)
}

func (c *Client) GetHistorical(symbols string, date time.Time) (data.Rates, error) {
	response, err := c.request(
		"historical",
		map[string]string{
			`symbols`: symbols,
			`date`:    date.Format("2006-01-02"),
		},
	)
	if err != nil {
		return data.Rates{}, err
	}

	return c.parseResponse(response)
}
