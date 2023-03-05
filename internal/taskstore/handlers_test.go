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
		t.Errorf("handler returned the wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestGetAllTasksHandler(t *testing.T) {
	// TODO
}

func TestGetTaskByTagHandler(t *testing.T) {
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

	// Creates a GET request to the `/tag/` directory and with `taskJSON`
	// as body
	r, err := http.NewRequest("GET", "/tag/tag_1", bytes.NewBuffer(taskJSON))
	if err != nil {
		t.Fatal(err)
	}

	// Creates a ResponseRecorder, which is an implementation of
	// `http.ResponseWriter`
	w := httptest.NewRecorder()

	// Crestes a task server
	ts := NewTaskServer()

	// Serves HTTP
	handler := http.HandlerFunc(ts.getTaskByTagHandler)
	handler.ServeHTTP(w, r)

	// Checks the HTTP reponse status code
	status := w.Code
	if status != http.StatusOK {
		t.Errorf("handler returned the wrong status code: got %v want %v", status, http.StatusOK)
	}
}
