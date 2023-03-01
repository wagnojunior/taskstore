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

	// Increments `nextID` of the taskstore so the next task won't have the
	// same ID as the previous
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

// `DeleteAllTasks` deletes all tasks from the taskstore
func (ts *TaskStore) DeleteAllTasks() error {
	// Deletes all tasks by creating a new map
	ts.tasks = make(map[int]Task)
	return nil
}

// `GetAllTasks` returns all tasks in the taskstore, in no particular order
func (ts *TaskStore) GetAllTasks() []Task {
	// Creates an empty slice of tasks
	allTasks := make([]Task, 0, len(ts.tasks))
	for _, t := range ts.tasks {
		allTasks = append(allTasks, t)
	}

	return allTasks
}

// `GetTaskByTag` returns all tasks that match the given tag
func (ts *TaskStore) GetTaskByTag(tag string) []Task {
	// Creates an empty slice of tasks
	taskByTag := make([]Task, 0, len(ts.tasks))

	for _, t := range ts.tasks { // Loops through all tasks in the taskstore
		for _, tg := range t.Tags { // Loops through all tags in the task
			if tg == tag {
				taskByTag = append(taskByTag, t)
			}
		}
	}

	return taskByTag
}

// `GetTaskByDueDate` returns all tasks that match the diven date
func (ts *TaskStore) GetTaskByDueDate(year int, month time.Month, day int) []Task {
	// Creates an empty slice of tasks
	taskByDate := make([]Task, 0, len(ts.tasks))

	for _, t := range ts.tasks { // Loops through all tasks in the taskstore
		// Gets the year, month, dar for task `t`
		y, m, d := t.Due.Date()
		if year == y && month == m && day == d {
			taskByDate = append(taskByDate, t)
		}
	}

	return taskByDate
}
