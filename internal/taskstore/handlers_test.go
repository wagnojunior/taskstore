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
		Tags: []string{"Tag"},
		Due:  time.Now(),
	}

	taskJSON, err := json.Marshal(task)

	// Crestes a task server
	ts := NewTaskServer()

	// Creates a request to pass to the tested handler
	r, err := http.NewRequest("POST", "/task/", bytes.NewBuffer(taskJSON))
	if err != nil {
		t.Fatal(err)
	}

	// Creates a ResponseRecorder
	w := httptest.NewRecorder()

	// Serves HTTP
	handler := http.HandlerFunc(ts.createTaskHandler)
	handler.ServeHTTP(w, r)

	// Checks status
	status := w.Code
	if status != http.StatusOK {
		t.Errorf("handler returned the wrong status code: got %v want %v", status, http.StatusOK)
	}

}
