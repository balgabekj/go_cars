package model

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Car struct {
	ID           string  `json:"id"`
	Model        string  `json:"model"`
	Brand        string  `json:"brand"`
	Year         *int    `json:"year"`
	Color        string  `json:"color"`
	Price        float64 `json:"price"`
	IsUsed       bool    `json:"isUsed"`
	UserID       string  `json:"userId"`
	CategoryName string  `json:"categoryName"`
}

type CarModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (m CarModel) Insert(car *Car) error {
	query := `
        INSERT INTO cars (model, brand, year, color, price, isUsed, userID, categoryName)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id
    `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, car.Model, car.Brand, car.Year, car.Color, car.Price, car.IsUsed, car.UserID, car.CategoryName).Scan(&car.ID)
}

func (m CarModel) Get(id string) (*Car, error) {
	query := `
        SELECT id, model, brand, year, color, price, isUsed, userID, categoryName
        FROM cars
        WHERE id = $1
    `

	var car Car

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&car.ID, &car.Model, &car.Brand, &car.Year, &car.Color, &car.Price, &car.IsUsed, &car.UserID, &car.CategoryName)
	if err != nil {
		return nil, err
	}
	return &car, nil
}

func (m CarModel) GetAll(brand string, minYear int, maxYear int, filters Filters) ([]Car, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id,model, brand, year,color, price, isUsed, userId, categoryName
		FROM cars
		WHERE (LOWER(brand) = LOWER($1) OR $1 = '')
		AND (year >= $2 OR $2 = 0)
		AND (year <= $3 OR $3 = 0)
		ORDER BY %s %s, id ASC
		LIMIT $4 OFFSET $5
	`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	args := []interface{}{brand, minYear, maxYear, filters.limit(), filters.offset()}
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	defer func() {
		if err := rows.Close(); err != nil {
			m.ErrorLog.Println(err)
		}
	}()

	totalRecords := 0

	var cars []Car

	for rows.Next() {
		var car Car
		err := rows.Scan(&totalRecords, &car.ID, &car.Model, &car.Brand, &car.Year, &car.Color, &car.Price, &car.IsUsed, &car.UserID, &car.CategoryName)
		if err != nil {
			return nil, Metadata{}, err
		}
		cars = append(cars, car)
	}

	if err := rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return cars, metadata, nil
}

func (m CarModel) Update(car *Car) error {
	query := `
        UPDATE cars
        SET model = $1, brand = $2, year = $3, color = $4, price = $5, isUsed = $6
        WHERE id = $7
    `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, car.Model, car.Brand, car.Year, car.Color, car.Price, car.IsUsed, car.ID)
	return err
}

func (m CarModel) Delete(id string) error {
	query := `
        DELETE FROM cars
        WHERE id = $1
    `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

func (m *CarModel) GetByCategory(categoryName string) ([]*Car, error) {
	query := `
		SELECT 
			c.id, 
			c.model, 
			c.brand, 
			c.year, 
			c.price, 
			c.color, 
			c.isUsed, 
			c.userId, 
			c.categoryName
		FROM 
			cars c
		JOIN 
			category cat ON c.categoryName = cat.name
		WHERE 
			c.categoryName = $1
		ORDER BY 
			c.id;
	`

	var cars []*Car

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, categoryName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var car Car
		err := rows.Scan(&car.ID, &car.Model, &car.Brand, &car.Year, &car.Price, &car.Color, &car.IsUsed, &car.UserID, &car.CategoryName)
		if err != nil {
			return nil, err
		}
		cars = append(cars, &car)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return cars, nil
}
