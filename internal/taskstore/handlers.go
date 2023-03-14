package taskstore

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
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

// `createTaskHandler` handles the creation of new tasks
func (ts *taskServer) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
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
func (ts *taskServer) GetAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("getAllTasksHandler at %s\n", r.URL.Path)

	// Gets all the tasks and marshals it to JSON
	allTasks := ts.store.GetAllTasks()
	utils.RenderJSON(w, allTasks)
}

// getTaskHandler handles the getting of a task
func (ts *taskServer) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("getTaskHandler at %s\n", r.URL.Path)

	// Gets a `map[string]string` associated with the http resquest `r`, more specifically the value
	// associated with the key `id`. Then, converts from str to int
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Gets the tasks by ID and marshals it to JSON
	task, err := ts.store.GetTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	utils.RenderJSON(w, task)
}

// deleteTaskHandler handles the deleting of a task
func (ts *taskServer) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("deleteTaskHandler at %s\n", r.URL.Path)

	// Gets a `map[string]string` associated with the http resquest `r`, more specifically the value
	// associated with the key `id`. Then, converts from str to int
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Deletes the task by ID
	err = ts.store.DeleteTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// deleteAllTasksHandler handles the deleting of all tasks
func (ts *taskServer) DeleteAllTasksHandler(w http.ResponseWriter, r *http.Request) {
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
func (ts *taskServer) GetTaskByTagHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("TagHandler at %s\n", r.URL.Path)

	// Gets a `map[string]string` associated with the http resquest `r`, more specifically the value
	// associated with the key `tag`
	tag := mux.Vars(r)["tag"]

	// Get the task by tag and marshal it into JSON
	tasks := ts.store.GetTaskByTag(tag)
	utils.RenderJSON(w, tasks)
}

// /////////////////////////////////////////////////////////////////////////////////////////////////
// TAG HANDLER
// /////////////////////////////////////////////////////////////////////////////////////////////////

// DueHandler handles the getting of a task by due date
func (ts *taskServer) GetTaskByDueHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("getTaskByDueDateHandler at %s\n", r.URL.Path)

	// Gets a `map[string]string` associated with the http resquest `r`
	vars := mux.Vars(r)
	year, err := strconv.Atoi(vars["year"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	month, err := strconv.Atoi(vars["month"])
	if err != nil || month < int(time.January) || month > int(time.December) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	day, err := strconv.Atoi(vars["day"])
	if err != nil || day > utils.DaysIn(month, year) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the task by date and marshal it into JSON
	tasks := ts.store.GetTaskByDueDate(year, time.Month(month), day)
	utils.RenderJSON(w, tasks)
}
