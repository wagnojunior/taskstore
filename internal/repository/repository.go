package repository

import "github.com/wagnojunior/taskstore/internal/models"

type DBRepo interface {
	CreateStore(name string) (*models.Store, error)
}
