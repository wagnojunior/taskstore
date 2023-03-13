package taskstore

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/wagnojunior/taskstore/internal/utils"
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

// /////////////////////////////////////////////////////////////////////////////
// TASK HANDLERS
// /////////////////////////////////////////////////////////////////////////////

func (ts *taskServer) TaskHandler(w http.ResponseWriter, r *http.Request) {
	// Checks if the request has an ID associated with if
	if r.URL.Path == "/task/" { // Request without an ID
		// Checks the requet type
		if r.Method == http.MethodPost { // Creates a new task
			ts.createTaskHandler(w, r)
		} else if r.Method == http.MethodGet { // Gets all tasks
			ts.getAllTasksHandler(w, r)
		} else if r.Method == http.MethodDelete { // Deletes all tasks
			ts.deleteAllTasksHandler(w, r)
		} else { //
			http.Error(w, fmt.Sprintf("expect method POST, GET, DELETE but got %v", r.Method), http.StatusMethodNotAllowed)
			return
		}
	} else { // Request with an ID
		path := strings.Trim(r.URL.Path, "/")
		subPath := strings.Split(path, "/")
		if len(subPath) < 2 {
			http.Error(w, "expected /task/<id> in task handler",
				http.StatusBadRequest)
		}

		// Gets the ID
		id, err := strconv.Atoi(subPath[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Checks the method
		if r.Method == http.MethodDelete { // Deletes a task
			ts.deleteTaskHandler(w, r, id)
		} else if r.Method == http.MethodGet { // Gets a task
			ts.getTaskHandler(w, r, id)
		} else {
			http.Error(w, fmt.Sprintf("expect method POST, GET, DELETE but got %v", r.Method), http.StatusMethodNotAllowed)
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
	utils.RenderJSON(w, task.ID)

}

// getAllTasksHandler handles the getting all tasks
func (ts *taskServer) getAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("getAllTasksHandler at %s\n", r.URL.Path)

	// Gets all the tasks and marshals it to JSON
	allTasks := ts.store.GetAllTasks()
	utils.RenderJSON(w, allTasks)
}

// getTaskHandler handles the getting of a task
func (ts *taskServer) getTaskHandler(w http.ResponseWriter, r *http.Request, id int) {
	log.Printf("getTaskHandler at %s\n", r.URL.Path)

	// Gets the tasks by ID and marshals it to JSON
	task, err := ts.store.GetTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	utils.RenderJSON(w, task)
}

// deleteTaskHandler handles the deleting of a task
func (ts *taskServer) deleteTaskHandler(w http.ResponseWriter, r *http.Request, id int) {
	log.Printf("deleteTaskHandler at %s\n", r.URL.Path)

	// Deletes the task by ID
	err := ts.store.DeleteTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// deleteAllTasksHandler handles the deleting of all tasks
func (ts *taskServer) deleteAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("deleteAllTasksHandler at %s\n", r.URL.Path)

	// Deletes all tasks
	err := ts.store.DeleteAllTasks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// /////////////////////////////////////////////////////////////////////////////////////////////////
// TAG HANDLER
// /////////////////////////////////////////////////////////////////////////////////////////////////

// TagHandler handles the getting of a task by tag
func (ts *taskServer) TagHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("TagHandler at %s\n", r.URL.Path)

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
		http.Error(w, "expected a path with a least one tag, got none", http.StatusBadRequest)
	}
	tag := subPath[1]

	// Get the task by tag and marshal it into JSON
	tasks := ts.store.GetTaskByTag(tag)
	utils.RenderJSON(w, tasks)
}

// /////////////////////////////////////////////////////////////////////////////////////////////////
// TAG HANDLER
// /////////////////////////////////////////////////////////////////////////////////////////////////

// DueHandler handles the getting of a task by due date
func (ts *taskServer) DueHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("getTaskByDueDateHandler at %s\n", r.URL.Path)

	// Checks if requet method is a GET
	if r.Method != http.MethodGet {
		http.Error(
			w,
			fmt.Sprintf("expected method GET /due/<year>/<month>/<day>, got %s", r.Method),
			http.StatusMethodNotAllowed)
	}

	// Gets the `tag` from the URL
	path := strings.Trim(r.URL.Path, "/")
	subPath := strings.Split(path, "/")

	// Checks if there are at least one tag in the path. If length of `subPath`
	// is less than 4, then there are missing fields
	if len(subPath) < 4 {
		http.Error(w, "expected a path with format /due/<year>/<month>/<day>",
			http.StatusBadRequest)
	}

	// Checks the validity of each field in the path
	year, err := strconv.Atoi(subPath[1])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	month, err := strconv.Atoi(subPath[2])
	if err != nil || month < int(time.January) || month > int(time.December) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	day, err := strconv.Atoi(subPath[3])
	if err != nil || day > utils.DaysIn(time.Month(month), year) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the task by date and marshal it into JSON
	tasks := ts.store.GetTaskByDueDate(year, time.Month(month), day)
	utils.RenderJSON(w, tasks)
}
