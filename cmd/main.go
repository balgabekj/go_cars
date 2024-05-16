package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/balgabekj/go_car/pkg/jsonlog"
	"github.com/balgabekj/go_car/pkg/model"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	"github.com/peterbourgon/ff/v3"
	"os"
	"sync"

	_ "github.com/lib/pq"
)

type config struct {
	port       int
	env        string
	fill       bool
	migrations string
	db         struct {
		dsn string
	}
}

//var (
//	version = vcs.Version()
//)

type application struct {
	config config
	models model.Models
	logger *jsonlog.Logger
	wg     sync.WaitGroup
}

func main() {
	fmt.Println("Started server")
	//var cfg config
	//flag.IntVar(&cfg.port, "port", 8080, "API server port")
	//flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	//flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://postgres:postgres@localhost:5432/gocars?sslmode=disable", "PostgreSQL DSN")
	//flag.Parse()

	fs := flag.NewFlagSet("demo-app", flag.ContinueOnError)

	var (
		cfg        config
		fill       = fs.Bool("fill", false, "Fill database with dummy data")
		migrations = fs.String("migrations", "", "Path to migration files folder. If not provided, migrations do not applied")
		port       = fs.Int("port", 8080, "API server port")
		env        = fs.String("env", "development", "Environment (development|staging|production)")
		dbDsn      = fs.String("dsn", "postgres://postgres:postgres@localhost:5432/gocars?sslmode=disable", "PostgreSQL DSN")
	)

	// Connect to DB
	logger := jsonlog.NewLogger(os.Stdout, jsonlog.LevelInfo)
	if err := ff.Parse(fs, os.Args[1:], ff.WithEnvVars()); err != nil {
		logger.PrintFatal(err, nil)
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}

	cfg.port = *port
	cfg.env = *env
	cfg.fill = *fill
	cfg.db.dsn = *dbDsn
	cfg.migrations = *migrations

	//logger.PrintInfo("starting application with configuration", map[string]string{
	//	"port":       fmt.Sprintf("%d", cfg.port),
	//	"fill":       fmt.Sprintf("%t", cfg.fill),
	//	"env":        cfg.env,
	//	"db":         cfg.db.dsn,
	//	"migrations": cfg.migrations,
	//})
	db, err := openDB(cfg)
	if err != nil {
		logger.PrintError(err, nil)
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.PrintFatal(err, nil)
		}
	}()
	app := &application{
		config: cfg,
		models: model.NewModels(db),
		logger: logger,
	}
	if err := app.serve(); err != nil {
		logger.PrintFatal(err, nil)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	if cfg.migrations != "" {
		driver, err := postgres.WithInstance(db, &postgres.Config{})
		if err != nil {
			return nil, err
		}
		err = db.Ping()
		if err != nil {
			return nil, err
		}
		m, err := migrate.NewWithDatabaseInstance(
			cfg.migrations,
			"postgres", driver)
		if err != nil {
			return nil, err
		}
		m.Up()
	}
	return db, nil
}
