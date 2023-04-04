package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/wagnojunior/taskstore/internal/utils"
)

// `Task` defines the structure of a task.
type Task struct {
	StoreID int      `json:"store_id"`
	ID      int      `json:"id"`
	Text    string   `json:"text"`
	Tags    []string `json:"tags"`
	Due     string   `json:"due"`
}

// `TaskService` defines the connection to the database
type TaskService struct {
	DB *sql.DB
}

// `Create` creates a new task, inserts it into the database, and returns it.
func (ts *TaskService) Create(storeID int, text string, tags []string, due string) (*Task, error) {
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

// `GetByID` gets a task by ID from the given store
func (ts *TaskService) GetByID(storeID, id int) (*Task, error) {
	task := Task{
		StoreID: storeID,
		ID:      id,
	}
	// `Task.Tags` is a slice of string, however it is saved as a
	// comma-separated string in the database. Therefore, when querying from
	// the database there is a type mismatch: string->[]string. `auxTags`
	// receives the data from the database for sebsequent conversion
	// string->[]string
	var auxTags string

	row := ts.DB.QueryRow(`
		SELECT text, tags, due
		FROM tasks
		WHERE store_id = ($1)
		AND id = ($2)`,
		storeID, id)

	err := row.Scan(&task.Text, &auxTags, &task.Due)
	if err != nil {
		return nil, fmt.Errorf("get task by ID: %w", err)
	}

	// Converts `auxTags` to a slice of strings, and adds it to `task`
	task.Tags = utils.StrToSlice(auxTags)

	log.Println("Task retrieved by ID!")
	return &task, nil
}

// `GetAll` gets all tasks from a determined store
func (ts *TaskService) GetAll(storeID int) (*[]Task, error) {
	var size int

	row := ts.DB.QueryRow(`
		SELECT COUNT (*)
		FROM tasks
		WHERE store_id = ($1)`,
		storeID)

	err := row.Scan(&size)
	if err != nil {
		return nil, fmt.Errorf("get task by ID: %w", err)
	}

	tasks := make([]Task, 0, size)
	// `Task.Tags` is a slice of string, however it is saved as a
	// comma-separated string in the database. Therefore, when querying from
	// the database there is a type mismatch: string->[]string. `auxTags`
	// receives the data from the database for sebsequent conversion
	// string->[]string
	var auxTags string

	rows, err := ts.DB.Query(`
		SELECT *
		FROM tasks
		WHERE store_id = ($1)`,
		storeID)
	if err != nil {
		return nil, fmt.Errorf("get task by ID: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var task Task

		err = rows.Scan(&task.ID, &task.Text, &auxTags, &task.Due,
			&task.StoreID)
		if err != nil {
			return nil, fmt.Errorf("get task by ID: %w", err)
		}

		// Converts `auxTags`  to a slice of strings, and adds it to `task`
		task.Tags = utils.StrToSlice(auxTags)
		tasks = append(tasks, task)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("get task by ID: %w", err)
	}

	log.Println("All tasks retrieved!")
	return &tasks, nil
}

// `DeleteByID` deletes a task by ID from the given store
func (ts *TaskService) DeleteByID(storeID, id int) error {
	var delID int

	row := ts.DB.QueryRow(`
		DELETE FROM tasks
		WHERE store_id = ($1)
		AND id = ($2)
		RETURNING id`,
		storeID, id)

	err := row.Scan(&delID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("task does not exist: %w", err)
		} else {
			return fmt.Errorf("delete task by ID: %w", err)
		}
	}

	log.Println("Task deleted by ID!")
	return nil
}

// `DeleteAll` deletes all taks from the given store and returns the number of
// tasks deleted and an error, if any
func (ts *TaskService) DeleteAll(storeID int) (int, error) {
	var delIDs []int

	rows, err := ts.DB.Query(`
		DELETE FROM tasks
		WHERE store_id = ($1)
		RETURNING id`,
		storeID)
	if err != nil {
		return 0, fmt.Errorf("delete all tasks: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int

		err = rows.Scan(&id)
		if err != nil {
			return 0, fmt.Errorf("delete all tasks: %w", err)
		}

		delIDs = append(delIDs, id)
	}
	err = rows.Err()
	if err != nil {
		return 0, fmt.Errorf("delete all tasks: %w", err)
	}

	log.Println("All tasks deleted")
	return len(delIDs), nil
}
