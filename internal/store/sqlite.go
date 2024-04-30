package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/parmaster/currency-api/internal/data"
)

type SQLiteStorage struct {
	DB  *sql.DB
	ctx context.Context
}

func NewSQLite(ctx context.Context, path string) (*SQLiteStorage, error) {

	sqliteDatabase, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	go func() {
		<-ctx.Done()
		sqliteDatabase.Close()
	}()

	q := `
	CREATE TABLE IF NOT EXISTS rates (
		date TEXT,
		base TEXT,
		currency TEXT,
		rate REAL,
		PRIMARY KEY (date, base, currency)
	);
	CREATE TABLE IF NOT EXISTS log (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		dateTime TEXT,
		type TEXT,
		request TEXT
	);
	`
	_, err = sqliteDatabase.ExecContext(ctx, q)
	if err != nil {
		return nil, err
	}

	return &SQLiteStorage{DB: sqliteDatabase, ctx: ctx}, nil
}

func (s *SQLiteStorage) Write(d data.Rates) error {

	for currency, rate := range d.Rates {
		q := `REPLACE INTO rates VALUES ($1, $2, $3, $4)`
		_, err := s.DB.ExecContext(s.ctx, q, d.Date.String(), d.Base, currency, rate)
		if err != nil {
			return err
		}
	}

	return nil
}

type line struct {
	date     string
	base     string
	currency string
	rate     float64
}

// Read reads rates from the database
func (s *SQLiteStorage) Read(date time.Time) (res data.Rates, err error) {

	q := fmt.Sprintf("SELECT * FROM `rates` WHERE `date` = '%s'", date.Format("2006-01-02"))
	rows, err := s.DB.QueryContext(s.ctx, q)
	if err != nil {
		return res, err
	}
	defer rows.Close()

	res.Rates = make(map[string]data.FloatRate)

	line := line{}
	for rows.Next() {
		err = rows.Scan(&line.date, &line.base, &line.currency, &line.rate)
		if err != nil {
			return res, err
		}

		err = res.Date.ParseDate(line.date)
		if err != nil {
			return res, err
		}
		res.Base = line.base
		res.Rates[line.currency] = data.FloatRate(line.rate)
	}

	return
}

// cleanup drops the rates table, used for testing
func (s *SQLiteStorage) cleanup() {
	s.DB.Exec("DROP TABLE `rates`")
}