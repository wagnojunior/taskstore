package taskstore

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
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
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)

}

func (ts *taskServer) getAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("getAllTasksHandler at %s\n", r.URL.Path)

	// Gets all the tasks and marshals it to JSON
	allTasks := ts.store.GetAllTasks()
	json, err := json.Marshal(allTasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Writes to the http response writer
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func (ts *taskServer) getTaskHandler(w http.ResponseWriter, r *http.Request, id int) {
	log.Printf("getTaskHandler at %s\n", r.URL.Path)

	// Gets the tasks by ID and marshals it to JSON
	task, err := ts.store.GetTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Writes to the http response writer
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func (ts *taskServer) deleteTaskHandler(w http.ResponseWriter, r *http.Request, id int) {
	log.Printf("deleteTaskHandler at %s\n", r.URL.Path)

	// Deletes the task by ID
	err := ts.store.DeleteTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ts *taskServer) deleteAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("deleteAllTasksHandler at %s\n", r.URL.Path)

	// Deletes all tasks
	err := ts.store.DeleteAllTasks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ts *taskServer) getTaskByTagHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("getTaskByTagHandler at %s\n", r.URL.Path)

	// Checks if requet method is a GET
	if r.Method != http.MethodGet {
		http.Error(
			w,
			fmt.Sprintf("expected method GET /tag/<tag>, got %s", r.Method),
			http.StatusMethodNotAllowed)
	}

	// Get the `tag` from the URL
	path := strings.Trim(r.URL.Path, "/")
	subPath := strings.Split(path, "/")
	println(path, subPath)

	// Check if there are at least one tag in the path. If length of `subPath`
	// is less than 2, then there are no tags in the path
	if len(subPath) < 2 {
		http.Error(w, "expected a path with a least one tag, not none", http.StatusBadRequest)
	}
	tag := subPath[1]

	// Get the task by tag and marshal it into JSON
	tasks := ts.store.GetTaskByTag(tag)
	json, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Writes to the http response writer
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func (ts *taskServer) getTaskByDueDateHandler(w http.ResponseWriter, r *http.Request, year int,
	month time.Month, day int) {

	log.Printf("getTaskByDueDateHandler at %s\n", r.URL.Path)

	// Get the task by date and marshal it into JSON
	tasks := ts.store.GetTaskByDueDate(year, month, day)
	json, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Writes to the http response writer
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
