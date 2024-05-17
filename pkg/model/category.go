package model

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type Category struct {
	Name string `json:"name"`
}

type CategoryModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (m *CategoryModel) InsertCategory(category *Category) error {
	query := `
        INSERT INTO category (name)
        VALUES ($1)
    `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, category.Name)
	return err
}

func (m *CategoryModel) UpdateCategory(oldName string, category *Category) error {
	query := `
        UPDATE category
        SET name = $1
        WHERE name = $2
    `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, category.Name, oldName)
	return err
}

func (m *CategoryModel) DeleteCategory(name string) error {
	query := `
        DELETE FROM category
        WHERE name = $1
    `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, name)
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
