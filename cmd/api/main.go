package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"amadeus.m7hir.net/internal/jsonlog"
)

const version = "1.0.0"

const baseUrl = "https://api.reccobeats.com"

type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *jsonlog.Logger
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	app := &application{
		config: cfg,
		logger: logger,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.logger.PrintInfo("starting server", map[string]string{"env": cfg.env, "addr": srv.Addr})

	if err := srv.ListenAndServe(); err != nil {
		app.logger.PrintFatal(err, map[string]string{"addr": srv.Addr})
	}

}
