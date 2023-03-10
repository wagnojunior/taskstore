package taskstore

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCreateTaskHandler(t *testing.T) {
	// Creates a test task and marshals it into a JSON format
	task := Task{
		ID:   1,
		Text: "First task",
		Tags: []string{"Tag_1", "Tag_2"},
		Due:  time.Now(),
	}
	taskJSON, err := json.Marshal(task)
	if err != nil {
		t.Fatal(err)
	}

	// Creates a POST request to the `/tas/` directory and with `taskJSON` as
	// body
	r, err := http.NewRequest("POST", "/task/", bytes.NewBuffer(taskJSON))
	if err != nil {
		t.Fatal(err)
	}

	// Creates a ResponseRecorder, which is an implementation of
	// `http.ResponseWriter`
	w := httptest.NewRecorder()

	// Crestes a task server
	ts := NewTaskServer()

	// Serves HTTP
	handler := http.HandlerFunc(ts.createTaskHandler)
	handler.ServeHTTP(w, r)

	// Checks the HTTP reponse status code
	status := w.Code
	if status != http.StatusOK {
		t.Errorf("handler returned the wrong status code: got %v wanted %v", status, http.StatusOK)
	}
}

func TestGetAllTasksHandler(t *testing.T) {
	// TODO
}

func TestGetTaskByTagHandler(t *testing.T) {
	// Creates test task and a task server
	task := Task{
		ID:   0,
		Text: "First task",
		Tags: []string{"tag_1", "tag_2"},
		Due:  time.Now(),
	}
	ts := NewTaskServer()
	_ = ts.store.CreateTask(task.Text, task.Tags, task.Due)

	// Creates a GET request to the `/tag/` directory
	r, err := http.NewRequest("GET", "/tag/tag_1", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Creates a ResponseRecorder, which is an implementation of
	// `http.ResponseWriter`
	w := httptest.NewRecorder()

	// Serves HTTP
	handler := http.HandlerFunc(ts.getTaskByTagHandler)
	handler.ServeHTTP(w, r)

	// Gets the task from the response writer and unmarshal it into a struct
	var got []Task
	err = json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}

	if !compareTask(task, got[0]) {
		t.Errorf("handler returned the wrong task: got %v wanted %v", got[0], task)
	}
}

func TestGetTaskByDueDateHandler(t *testing.T) {
	// Creates test task and a task server
	task := Task{
		ID:   0,
		Text: "Test task",
		Tags: []string{"tag_1", "tag_2"},
		Due:  time.Date(2023, time.September, 5, 0, 0, 0, 0, time.UTC),
	}
	ts := NewTaskServer()
	_ = ts.store.CreateTask(task.Text, task.Tags, task.Due)

	// Creates a GET request to the `/tag/` directory
	r, err := http.NewRequest("GET", "/due/2023/9/5", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Creates a ResponseRecorder, which is an implementation of
	// `http.ResponseWriter`
	w := httptest.NewRecorder()

	// Serves HTTP
	handler := http.HandlerFunc(ts.getTaskByDueDateHandler)
	handler.ServeHTTP(w, r)

	// Gets the task from the response writer and unmarshal it into a struct
	var got []Task
	err = json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}

	if !compareTask(task, got[0]) {
		t.Errorf("handler returned the wrong task: got %v wanted %v", got[0], task)
	}
}
