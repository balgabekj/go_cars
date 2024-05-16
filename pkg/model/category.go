package model

import (
	"database/sql"
	"log"
)

type Category struct {
	Name string `json:"name"`
}

type CategoryModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}
