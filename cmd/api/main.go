package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-pkgz/lgr"
	"github.com/jessevdk/go-flags"
)

type Options struct {
	Port       int      `long:"port" short:"p" env:"PORT" description:"Listening port" default:"8080"`
	ApiKey     string   `long:"apikey" env:"APIKEY" description:"currencyfreaks.com API key" required:"true"`
	Currencies []string `long:"currencies" env:"CURRENCIES" description:"currency codes to use" default:"UAH,USD,EUR,RON"`
	Interval   int      `long:"interval" env:"INTERVAL" description:"update interval in seconds" default:"3600"`
	Debug      bool     `long:"dbg" env:"DEBUG" description:"Enable debug mode with verbose logging"`
	Version    bool     `short:"v" description:"Show version and exit"`
}

type Server struct {
	cfg Options
	ctx context.Context
}

func NewServer(cfg Options, ctx context.Context) *Server {
	return &Server{cfg: cfg, ctx: ctx}
}

func (s *Server) Run() {

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.cfg.Port),
		Handler:      s.router(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		<-s.ctx.Done()
		log.Printf("[INFO] Terminating http server")

		if err := srv.Close(); err != nil {
			log.Printf("[ERROR] failed to close http server, %v", err)
		}
	}()

	log.Printf("[DEBUG] starting server with options: %+v", s.cfg)

	err := srv.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Printf("[ERROR] server failed: %v", err)
	}
}

func main() {
	// Parsing command line arguments
	var cfg Options
	p := flags.NewParser(&cfg, flags.PassDoubleDash|flags.HelpFlag)
	if _, err := p.Parse(); err != nil {
		if err.(*flags.Error).Type != flags.ErrHelp {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
		p.WriteHelp(os.Stderr)
		os.Exit(2)
	}

	// Logger setup
	logOpts := []lgr.Option{
		lgr.LevelBraces,
		lgr.StackTraceOnError,
		lgr.Secret(cfg.ApiKey),
	}
	if cfg.Debug {
		logOpts = append(logOpts, lgr.Debug)
	}
	lgr.SetupStdLogger(logOpts...)

	// Version
	if cfg.Version {
		fmt.Printf("Version: %s\n", version)
		os.Exit(0)
	}
	log.Printf("[DEBUG] Pid: %d, ver: %s", os.Getpid(), version)

	// Graceful termination
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		// catch signal and invoke graceful termination
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-stop
		log.Println("Shutdown signal received\n*********************************")
		cancel()
	}()

	// Recover from panics
	defer func() {
		if x := recover(); x != nil {
			log.Printf("[WARN] run time panic: %+v", x)
		}
	}()

	// Starting the server
	NewServer(cfg, ctx).Run()
}
