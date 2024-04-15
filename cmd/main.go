package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/balgabekj/go_car/pkg/jsonlog"
	"github.com/balgabekj/go_car/pkg/model"
	"github.com/gorilla/mux"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type config struct {
	port string
	env  string
	db   struct {
		dsn string
	}
}

type application struct {
	config config
	models model.Models
	logger *jsonlog.Logger
}

func main() {
	fmt.Println("Started server")
	var cfg config
	flag.StringVar(&cfg.port, "port", ":8081", "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://postgres:postgres@localhost:5432/gocars?sslmode=disable", "PostgreSQL DSN")
	flag.Parse()

	// Connect to DB
	db, err := openDB(cfg)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	app := &application{
		config: cfg,
		models: model.NewModels(db),
	}

	app.run()
}
func (app *application) run() {
	fmt.Println("Running")
	r := mux.NewRouter()

	v1 := r.PathPrefix("/api/v1").Subrouter()

	v1.HandleFunc("/cars", app.createCarHandler).Methods("POST")
	v1.HandleFunc("/cars/{id}", app.getCarHandler).Methods("GET")
	v1.HandleFunc("/cars", app.getAllCarHandler).Methods("GET")
	v1.HandleFunc("/cars/{id}", app.updateCarHandler).Methods("PUT")
	v1.HandleFunc("/cars/{id}", app.deleteCarHandler).Methods("DELETE")

	v1.HandleFunc("/owners", app.createOwnerHandler).Methods("POST")
	v1.HandleFunc("/owners/{id}", app.getOwnerHandler).Methods("GET")
	v1.HandleFunc("/owners", app.getAllOwnersHandler).Methods("GET")
	v1.HandleFunc("/owners/{id}", app.updateOwnerHandler).Methods("PUT")
	v1.HandleFunc("/owners/{id}", app.deleteOwnerHandler).Methods("DELETE")

	log.Printf("Starting server on %s\n", app.config.port)
	err := http.ListenAndServe(app.config.port, r)
	log.Fatal(err)
}
func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
