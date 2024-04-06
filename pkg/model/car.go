package model

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type Car struct {
	ID     string `json:"id"`
	Model  string `json:"model"`
	Brand  string `json:"brand"`
	Year   int    `json:"year"`
	Color  string `json:"color"`
	IsUsed bool   `json:"isUsed"`
}

type CarModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (m CarModel) Insert(car *Car) error {
	query := `
        INSERT INTO car (model, brand, year, color, isUsed)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, car.Model, car.Brand, car.Year, car.Color, car.IsUsed).Scan(&car.ID)
}

func (m CarModel) Get(id string) (*Car, error) {
	query := `
        SELECT id, model, brand, year, color, isUsed
        FROM car
        WHERE id = $1
    `

	var car Car

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&car.ID, &car.Model, &car.Brand, &car.Year, &car.Color, &car.IsUsed)
	if err != nil {
		return nil, err
	}
	return &car, nil
}

func (m CarModel) Update(car *Car) error {
	query := `
        UPDATE car
        SET model = $1, brand = $2, year = $3, color = $4, isUsed = $5
        WHERE id = $6
    `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, car.Model, car.Brand, car.Year, car.Color, car.IsUsed, car.ID)
	return err
}

func (m CarModel) Delete(id string) error {
	query := `
        DELETE FROM car
        WHERE id = $1
    `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}
