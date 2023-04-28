package repository

import "github.com/wagnojunior/taskstore/internal/models"

type DBRepo interface {
	CreateStore(name string) (*models.Store, error)
	DeleteStoreByID(id int) error

	CreateTask(storeID int, text string, tags []string, due string) (*models.Task, error)
	GetTaskByID(storeID, id int) (*models.Task, error)
	GetTaskByTags(storeID int, tags []string) (*[]models.Task, error)
	GetTaskByDue(storeID int, due string) (*[]models.Task, error)
	GetAllTasks(storeID int) (*[]models.Task, error)
	DeleteTaskByID(storeID, id int) error
	DeleteAllTasks(storeID int) (int, error)
}
