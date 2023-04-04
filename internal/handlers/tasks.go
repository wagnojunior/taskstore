package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/wagnojunior/taskstore/internal/models"
	"github.com/wagnojunior/taskstore/internal/utils"
)

type Tasks struct {
	TaskService *models.TaskService
}

// `Create` handles POST requests to create a task in a given store
func (t Tasks) Create(w http.ResponseWriter, r *http.Request) {
	var task models.Task

	b, err := io.ReadAll(r.Body)
	json.Unmarshal(b, &task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newTask, err := t.TaskService.Create(task.StoreID, task.Text, task.Tags,
		task.Due)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.RenderJSON(w, newTask)
}

// `GetByID` handles GET requests retrieve a task by ID from a given store
func (t Tasks) GetByID(w http.ResponseWriter, r *http.Request) {
	// Gets a `map[string]string` associated with the http resquest `r`
	reqData := mux.Vars(r)
	storeID, err := strconv.Atoi(reqData["store_id"])
	id, err := strconv.Atoi(reqData["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	gotTask, err := t.TaskService.GetByID(storeID, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.RenderJSON(w, gotTask)
}

// `GetAll` handles GET requests to retrieve all tasks from a given store
func (t Tasks) GetAll(w http.ResponseWriter, r *http.Request) {
	// Gets a `map[string]string` associated with the http request `r`
	reqData := mux.Vars(r)
	storeID, err := strconv.Atoi(reqData["store_id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	allTasks, err := t.TaskService.GetAll(storeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.RenderJSON(w, allTasks)
}

// `DeleteByID` handles POST requests to delete a task by ID from a given store
func (t Tasks) DeleteByID(w http.ResponseWriter, r *http.Request) {
	// Gets a `map[string]string associated with the http request `r`
	reqData := mux.Vars(r)
	storeID, err := strconv.Atoi(reqData["store_id"])
	id, err := strconv.Atoi(reqData["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = t.TaskService.DeleteByID(storeID, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (t Tasks) DeleAll(w http.ResponseWriter, r *http.Request) {
	// gets a `map[string]string` associated with the http request `r`
	reqData := mux.Vars(r)
	storeID, err := strconv.Atoi(reqData["store_id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	n, err := t.TaskService.DeleteAll(storeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type delIDs struct {
		Num int `json:"del_id"`
	}
	num := delIDs{
		Num: n,
	}
	utils.RenderJSON(w, num)

}
