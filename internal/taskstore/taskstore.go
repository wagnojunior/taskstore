package taskstore

import (
	"fmt"
	"time"
)

// `Task` defines the structure of a task.
type Task struct {
	ID   int       `json:"id"`
	Text string    `json:"text"`
	Tags []string  `json:"tags"`
	Due  time.Time `json:"due"`
}

// `TaskStore` is a simple in-memory database of tasks.
type TaskStore struct {
	tasks  map[int]Task
	nextID int
}

// `New` creates and returns a new TaskStore.
func New() *TaskStore {
	ts := TaskStore{
		tasks:  make(map[int]Task),
		nextID: 0,
	}
	return &ts
}

// `CreateTask` creates a new task and returns its ID
func (ts *TaskStore) CreateTask(text string, tags []string, due time.Time) int {
	// Creates a task with the data passed as argument
	task := Task{
		ID:   ts.nextID,
		Text: text,
		Tags: tags,
		Due:  due,
	}

	// Adds `task` to the taskstore
	ts.tasks[ts.nextID] = task

	// Increments `nextID` of the taskstore so the next task won't have the same ID
	// as the previous
	ts.nextID++

	return task.ID
}

// `GetTask` returns a task from the taskstore by id. If the id does not exists, returns an error.
func (ts *TaskStore) GetTask(id int) (Task, error) {
	task, ok := ts.tasks[id]
	if !ok {
		return Task{}, fmt.Errorf("there is no task with id: %d", id)
	}
	return task, nil
}

// `DeleteTask` deletes a task from the taskstore by id. If the id does not exists, returns an
// error.
func (ts *TaskStore) DeleteTask(id int) error {
	// Checks if the id exists
	_, ok := ts.tasks[id]
	if !ok {
		return fmt.Errorf("there is no task with id: %d", id)
	}

	delete(ts.tasks, id)
	return nil
}

func (ts *TaskStore) DeleteAllTasks() error {
	// Deletes all tasks by creating a new map
	ts.tasks = make(map[int]Task)
	return nil
}
