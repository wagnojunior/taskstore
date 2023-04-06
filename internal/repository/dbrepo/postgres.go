package dbrepo

import (
	"fmt"

	"github.com/wagnojunior/taskstore/internal/models"
)

func (pr *PostgresRepo) CreateStore(name string) (*models.Store, error) {
	store := models.Store{
		Name: name,
	}

	row := pr.DB.QueryRow(`
		INSERT INTO stores (name)
		VALUES ($1) RETURNING id`,
		name)

	err := row.Scan(&store.ID)
	if err != nil {
		return nil, fmt.Errorf("create store: %w", err)
	}

	return &store, nil
}
