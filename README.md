# Currency exchange rates API

This is a simple REST API that provides exchange rates for the specified list of currencies. The exchange rates are obtained from the currencyfreaks.com API and saved in a SQLite database.

## Prerequisites
Put an API key for the currencyfreaks.com API in the `config.ini` file:
```bash
apikey = your_api_key
```
or pass it as a command line argument:
```bash
./bin/api --apikey your_api_key
```
or set it as an environment variable:
```bash
APIKEY=your_api_key ./bin/api
```

## Docker build and run
Correct API key should be put in the `config.ini` file before building the docker container.

To build and run the docker container:
```bash
docker build -t currency-api .
docker run -it --rm -p 8080:8080 currency-api
```

API will be available at `http://localhost:8080/`

## Local build and run
```bash
make && ./bin/api
```
Will start the API server on port 8080

## Configuration
The configuration file is `config.ini` in the root of the project - it contains default settings which can be overridden by environment variables or command line arguments:

```bash
make && ./bin/api --apikey your_api_key --port 8081 --dbpath /tmp/currency-api.db --currencies UAH,USD,EUR,RON --debug
```
Full list of configuration options can be listed with `make && ./bin/api --help`

## Database
The database is created in the file specified in the configuration file. The database schema with sample data is created automatically on the first run of the API server.

## Logging
The API logs all requests to the database. Last 10 logs can be viewed with the `/v1/status/` endpoint.

## Testing
```bash
make test
```
Runs all available tests

## API endpoints
`/v1/rates/` - get latest exchange rates for the currencies specified in the config file
```json
{
	"date": "2024-04-30 00:00:00",
	"base": "USD",
	"rates": {
		"EUR": 0.933291146848332,
		"RON": 4.642474983710198,
		"UAH": 39.63845919120639,
		"USD": 1
	}
}
```

`/v1/rates/<date>/` - get exchange rates for the specified date (e.g. 2024-04-20)
```json
{
	"date": "2024-04-20 00:00:00",
	"base": "USD",
	"rates": {
		"EUR": 0.8,
		"RON": 4.7,
		"UAH": 39.4
	}
}
```

`/v1/pair/<pair>/` - get exchange rates for the specified currency pair (e.g. UAH-RON)
```json
{
	"date": "2024-04-30",
	"pair": "UAH-RON",
	"rate": 0.11712047033200801
}
```

`/v1/status/` - get the status of the API
```json
{
	"status": "ok",
	"version": "93952c1-main-20240501",
	"config": {
		"port": 8080,
		"dbpath": "file:/tmp/currency-api.db?mode=rwc\u0026_journal_mode=WAL",
		"currencies": "UAH,USD,EUR,RON",
		"interval": 3600,
		"debug": true
	},
	"logs": [
		"2024-05-01 01:45:39 | pair | pair: UAH-RON"
	]
}
```

Original job interview test task:

### Implement a REST API with the following functionality

1. Obtaining exchange rates on a third-party API (UAH / EUR / USD / RON). (eg: https://currencyfreaks.com/)

	1.1 If there was no exchange rate for the current day, get the exchange rate updated at 12:00

	1.2 Implement saving the exchange rate in the database (Sqlite for start, PostgreSQL for extra points)

2. Implement endpoints:

	2.1. Receiving rates for a specified date (all 4 currencies)

	2.2. Obtaining currency pairs: UAH-USD.

3. Implement validation of requests (dates, currency tickers)

4. Implement API access via API key

5. Implement logging of API requests in the database. Save in the logs - date, type of request (by date / by pair)

6. Use docker for the project.

7. When initializing the server, fill the database with initial data on exchange rates for several dates (any values)
