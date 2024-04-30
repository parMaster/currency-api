package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/parmaster/currency-api/internal/data"
)

/*
curl 'https://api.currencyfreaks.com/v2.0/rates/latest?apikey=APIKEY&symbols=PKR,GBP,EUR,USD'
{
    "date": "2023-03-21 13:26:00+00",
    "base": "USD",
    "rates": {
        "EUR": "0.9278605451274349",
        "GBP": "0.8172754173817152",
        "PKR": "281.6212943333344",
        "USD": "1.0"
    }
}
*/

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

// TODO:
func (c *Client) request(symbols []string, date string) ([]byte, error) {
	params := url.Values{}
	params.Add(`apikey`, c.ApiKey)
	params.Add(`symbols`, strings.Join(symbols, `,`))
	apiType := `latest`
	if date != "" {
		params.Add(`date`, date)
		apiType = `historical`
	}
	log.Printf("[DEBUG] CF request: %s?%s", c.ApiUrl[apiType], params.Encode())
	response, err := http.Get(fmt.Sprintf("%s?%s", c.ApiUrl[apiType], params.Encode()))
	if err != nil {
		return []byte{}, err
	}
	defer response.Body.Close()

	// dump response
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return []byte{}, err
	}
	// fmt.Println(string(body))

	return body, nil
}

func (c *Client) GetRates(symbols []string, date time.Time, testResponse []byte) (data.Rates, error) {
	var err error
	dateStr := date.Format("2006-01-02")
	if date.IsZero() {
		dateStr = ""
	}
	response := testResponse
	if len(testResponse) == 0 {
		response, err = c.request(symbols, dateStr)
		if err != nil {
			return data.Rates{}, err
		}
	}

	rates := data.Rates{}
	if err := json.Unmarshal(response, &rates); err != nil {
		return data.Rates{}, err
	}

	return rates, nil
}
