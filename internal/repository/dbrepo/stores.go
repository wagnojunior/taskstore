package dbrepo

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/wagnojunior/taskstore/internal/models"
)

// `CreateStore` creates a store with the given name, and returns a pointer to
// the store and an error
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

// `DeleteStoreByID` deletes a store with the given name and return an error
func (pr *PostgresRepo) DeleteStoreByID(id int) error {
	var delID int

	row := pr.DB.QueryRow(`
		DELETE FROM stores
		WHERE id = ($1)
		RETURNING id`,
		id)

	err := row.Scan(&delID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("store does not exist: %w", err)
		} else {
			return fmt.Errorf("delete store by ID: %w", err)
		}
	}

	log.Println("Store deleted by ID!")
	return nil
}
