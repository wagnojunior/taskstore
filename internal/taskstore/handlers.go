package taskstore

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// `taskServer` is a constructor for the server type. It wraps a `TaskStore`
type taskServer struct {
	store *TaskStore
}

// `NewTaskServer` returns a new task server
func NewTaskServer() *taskServer {
	return &taskServer{
		store: New(),
	}
}

func (ts *taskServer) TaskHandler(w http.ResponseWriter, r *http.Request) {
	// Checks if the request has an ID associated with if
	if r.URL.Path == "/task/" { // No ID associated
		// Checks the requet type
		if r.Method == http.MethodPost { // Creates a new task
			ts.createTaskHandler(w, r)
		} else if r.Method == http.MethodGet { // Gets all tasks
			ts.getAllTasksHandler(w, r)
		} else if r.Method == http.MethodDelete { // Deletes all tasks
			ts.deleteAllTasksHandler(w, r)
		} else { //
			http.Error(w, fmt.Sprintf("expect method POST, GET, DELETE but got %v", r.Method),
				http.StatusMethodNotAllowed)
			return
		}
	}
}

// `createTaskHandler` handles the creation of new tasks
func (ts *taskServer) createTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("createTaskHandler at %s\n", r.URL.Path)

	// Creates a new decoder and disallow unknowm fields
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	// Decodes the JSON data to Go struct
	var task Task
	err := dec.Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Creates a new task and converts it back to JSON
	_ = ts.store.CreateTask(task.Text, task.Tags, task.Due)
	json, err := json.Marshal(task.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Writes to the http response writer
	w.Write(json)

}

func (ts *taskServer) getAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (ts *taskServer) deleteAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}
