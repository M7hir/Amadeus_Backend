package main

import (
	"context"
	"database/sql"
	"flag"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"amadeus.m7hir.net/internal/data"
	"amadeus.m7hir.net/internal/deezer"
	"amadeus.m7hir.net/internal/jsonlog"
	"amadeus.m7hir.net/internal/mailer"
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
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	limiter struct {
		enabled bool
		rps     float64
		burst   int
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}

	cors struct {
		trustedOrigins []string
	}
}

type application struct {
	config     config
	logger     *jsonlog.Logger
	reccobeats *reccobeats.Client
	deezer     *deezer.Client
	models     data.Models
	mailer     mailer.Mailer
	wg         sync.WaitGroup
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

	// fmt.Println("DSN:", dbDSN)
	// fmt.Println("Port:", port)
	var cfg config

	flag.IntVar(&cfg.port, "port", port, "API server port")
	flag.StringVar(&cfg.db.dsn, "db-dsn", dbDSN, "Postgresql DSN")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")

	flag.StringVar(&cfg.smtp.host, "smtp-host", "sandbox.smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 2525, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "684e24253efbac", "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", "f9382eae204ba9", "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "Amadeus <no-reply@amadeus.m7hir.net>", "SMTP sender")

	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		cfg.cors.trustedOrigins = strings.Fields(val)
		return nil
	})

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
		models:     data.NewModels(db, logger),
		mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username,
			cfg.smtp.password, cfg.smtp.sender),
	}

	err = app.server()
	if err != nil {
		logger.PrintFatal(err, nil)
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
