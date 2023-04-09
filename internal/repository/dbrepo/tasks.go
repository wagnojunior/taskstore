package dbrepo

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/wagnojunior/taskstore/internal/models"
	"github.com/wagnojunior/taskstore/internal/utils"
)

// `CreateTask` creates a new task, inserts it into the database, and returns
// it.
func (pr *PostgresRepo) CreateTask(storeID int, text string, tags []string, due string) (*models.Task, error) {
	task := models.Task{
		StoreID: storeID,
		Text:    text,
		Tags:    tags,
		Due:     due,
	}

	row := pr.DB.QueryRow(`
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

// `GetTaskByID` gets a task by ID from the given store
func (pr *PostgresRepo) GetTaskByID(storeID, id int) (*models.Task, error) {
	task := models.Task{
		StoreID: storeID,
		ID:      id,
	}
	// `Task.Tags` is a slice of string, however it is saved as a
	// comma-separated string in the database. Therefore, when querying from
	// the database there is a type mismatch: string->[]string. `auxTags`
	// receives the data from the database for sebsequent conversion
	// string->[]string
	var auxTags string

	row := pr.DB.QueryRow(`
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

func (pr *PostgresRepo) GetTaskByTags(storeID int, tags []string) (*[]models.Task, error) {
	tasks := make([]models.Task, 0)
	// `Task.Tags` is a slice of string, however it is saved as a
	// comma-separated string in the database. Therefore, when querying from
	// the database there is a type mismatch: string->[]string. `auxTags`
	// receives the data from the database for sebsequent conversion
	// string->[]string
	var auxTags string
	var task models.Task

	for _, t := range tags {
		rows, err := pr.DB.Query(`
			SELECT *
			FROM tasks
			WHERE store_id = ($1)
			AND ($2) = ANY (tags)`,
			storeID, t)
		if err != nil {
			return nil, fmt.Errorf("get task by tags: %w", err)
		}
		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(&task.ID, &task.Text, &auxTags, &task.Due, &task.StoreID)
			if err != nil {
				return nil, fmt.Errorf("get task by tags: %w", err)
			}

			// Converts `auxTags`  to a slice of strings, and adds it to `task`
			task.Tags = utils.StrToSlice(auxTags)
			tasks = append(tasks, task)
		}
		err = rows.Err()
		if err != nil {
			return nil, fmt.Errorf("get task by ID: %w", err)
		}
	}

	// The same task can be queried multiple times if its tags match the search
	// more than once. In this case it is necessary to remove duplicates
	tasks = unique(tasks)

	log.Println("All tasks retrieved!")
	return &tasks, nil
}

// `GetAllTasks` gets all tasks from a determined store
func (pr *PostgresRepo) GetAllTasks(storeID int) (*[]models.Task, error) {
	tasks := make([]models.Task, 0)
	// `Task.Tags` is a slice of string, however it is saved as a
	// comma-separated string in the database. Therefore, when querying from
	// the database there is a type mismatch: string->[]string. `auxTags`
	// receives the data from the database for sebsequent conversion
	// string->[]string
	var auxTags string
	var task models.Task

	rows, err := pr.DB.Query(`
		SELECT *
		FROM tasks
		WHERE store_id = ($1)`,
		storeID)
	if err != nil {
		return nil, fmt.Errorf("get task by ID: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
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

// `DeleteTaskByID` deletes a task by ID from the given store
func (pr *PostgresRepo) DeleteTaskByID(storeID, id int) error {
	var delID int

	row := pr.DB.QueryRow(`
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

// `DeleteAllTasks` deletes all taks from the given store and returns the
// number of tasks deleted and an error, if any
func (pr *PostgresRepo) DeleteAllTasks(storeID int) (int, error) {
	var delIDs []int
	var id int

	rows, err := pr.DB.Query(`
		DELETE FROM tasks
		WHERE store_id = ($1)
		RETURNING id`,
		storeID)
	if err != nil {
		return 0, fmt.Errorf("delete all tasks: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
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

// `unique` returns a slice of unique tasks
func unique(tasks []models.Task) []models.Task {
	if len(tasks) < 2 {
		return tasks
	}

	unique := 0
	for i := 1; i < len(tasks); i++ {
		if tasks[unique].ID != tasks[i].ID {
			unique++
			tasks[unique] = tasks[i]
		}
	}

	return tasks[:unique+1]
}
