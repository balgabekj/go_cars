package model

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type Owner struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Number   string `json:"number"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type OwnerModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (o OwnerModel) Insert(owner *Owner) error {
	query := `
        INSERT INTO owner (name, number, email, password)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := o.DB.QueryRowContext(ctx, query, owner.Name, owner.Number, owner.Email, owner.Password).Scan(&owner.ID)
	if err != nil {
		return err
	}
	return nil
}

func (m OwnerModel) Get(id string) (*Owner, error) {
	query := `
        SELECT id, name, number, email, password
        FROM owner
        WHERE id = $1
    `

	var owner Owner

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&owner.ID, &owner.Name, &owner.Number, &owner.Email, &owner.Password)
	if err != nil {
		return nil, err
	}
	return &owner, nil
}

func (m OwnerModel) GetAll() ([]Owner, error) {
	query := `
        SELECT id, name, number, email, password
        FROM owner
    `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var owners []Owner
	for rows.Next() {
		var owner Owner
		err := rows.Scan(&owner.ID, &owner.Name, &owner.Number, &owner.Email, &owner.Password)
		if err != nil {
			return nil, err
		}
		owners = append(owners, owner)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return owners, nil
}

func (m OwnerModel) Update(owner *Owner) error {
	query := `
        UPDATE owner
        SET name = $1, number = $2, email = $3, password = $4
        WHERE id = $5
    `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, owner.Name, owner.Number, owner.Email, owner.Password, owner.ID)
	return err
}

func (m OwnerModel) Delete(id string) error {
	query := `
        DELETE FROM owner
        WHERE id = $1
    `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}
