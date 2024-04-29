# Currency exchange rates API

Job interview test task

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
