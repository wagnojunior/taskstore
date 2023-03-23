package models

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// `Task` defines the structure of a task.
type Task struct {
	StoreID int       `json:"store_id"`
	ID      int       `json:"id"`
	Text    string    `json:"text"`
	Tags    []string  `json:"tags"`
	Due     time.Time `json:"due"`
}

// `TaskService` defines the connection to the database
type TaskService struct {
	DB *sql.DB
}

// `Create` creates a new task, inserts it into the database, and returns it.
func (ts *TaskService) Create(storeID int, text string, tags []string, due time.Time) (*Task, error) {
	task := Task{
		StoreID: storeID,
		Text:    text,
		Tags:    tags,
		Due:     due,
	}

	row := ts.DB.QueryRow(`
		INSERT INTO tasks (store_id, text, tags, due)
		VALUES ($1, $2, $3, $4) RETURNING id`,
		storeID, text, tags, due)

	err := row.Scan(&task.ID)
	if err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}

	log.Println("Task created!")
	return &task, nil
}

// `GetByID` gets a task from the given store by ID
func (ts *TaskService) GetByID(storeID, id int) (*Task, error) {
	task := Task{
		StoreID: storeID,
		ID:      id,
	}
	var auxTags, auxDue string

	row := ts.DB.QueryRow(`
		SELECT text, tags, due
		FROM tasks
		WHERE store_id = ($1)
		AND id = ($2)`,
		storeID, id)

	err := row.Scan(&task.Text, &auxTags, &auxDue)
	if err != nil {
		return nil, fmt.Errorf("get task by ID: %w", err)
	}

	// TODO: convert auxTags to []string
	// TODO: convert auxDue to time.Time

	log.Println("Task retrieved by ID!")
	return &task, nil
}
