package model

import (
	"database/sql"
	"log"
)

type Owner struct {
	Name string `json:"name"`
}

type OwnerModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}
