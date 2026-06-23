package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"amadeus.m7hir.net/internal/data"
	"amadeus.m7hir.net/internal/deezer"
	"amadeus.m7hir.net/internal/jsonlog"
	"amadeus.m7hir.net/internal/reccobeats"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const version = "1.0.0"

const reccobeatsBaseUrl = "https://api.reccobeats.com"
const deezerBaseUrl = "https://api.deezer.com"

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
}

type application struct {
	config     config
	logger     *jsonlog.Logger
	reccobeats *reccobeats.Client
	deezer     *deezer.Client
	models     data.Models
}

func main() {

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	err := godotenv.Load()
	if err != nil {
		logger.PrintFatal(err, map[string]interface{}{"port-error": "Cannot find the port"})

	}
	dbDSN := os.Getenv("DB_DSN")
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		logger.PrintFatal(err, map[string]interface{}{"port-error": "Cannot convert the port"})
	}

	fmt.Println("DSN:", dbDSN)
	fmt.Println("Port:", port)
	var cfg config

	flag.IntVar(&cfg.port, "port", port, "API server port")
	flag.StringVar(&cfg.db.dsn, "db-dsn", dbDSN, "Postgresql DSN")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	db, err := OpenDB(cfg)
	if err != nil {
		logger.PrintFatal(err, map[string]interface{}{"db-error": "connection failed"})
	}
	defer db.Close()
	logger.PrintInfo("database connection pool established", map[string]interface{}{"db": cfg.db.dsn})

	app := &application{
		config:     cfg,
		logger:     logger,
		reccobeats: reccobeats.NewClient(reccobeatsBaseUrl),
		deezer:     deezer.NewClient(deezerBaseUrl),
		models:     data.NewModels(db,logger),
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.PrintInfo("starting server", map[string]interface{}{"env": cfg.env, "version": version, "port": cfg.port})

	if err := srv.ListenAndServe(); err != nil {
		logger.PrintFatal(err, map[string]interface{}{"env": cfg.env})
	}

}

func OpenDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
