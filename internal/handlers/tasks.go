package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/wagnojunior/taskstore/internal/models"
)

type Tasks struct {
	TaskService *models.TaskService
}

// `Create` handles POST requests
func (t Tasks) Create(w http.ResponseWriter, r *http.Request) {
	var task models.Task

	b, err := io.ReadAll(r.Body)
	json.Unmarshal(b, &task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = t.TaskService.Create(task.StoreID, task.Text, task.Tags,
		task.Due)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// `GetByID` handles GET requests
func (t Tasks) GetByID(w http.ResponseWriter, r *http.Request) {
	// Gets a `map[string]string` associated with the http resquest `r`
	reqData := mux.Vars(r)
	storeID, err := strconv.Atoi(reqData["store_id"])
	id, err := strconv.Atoi(reqData["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = t.TaskService.GetByID(storeID, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (t Tasks) GetAll(w http.ResponseWriter, r *http.Request) {
	// Gets a `map[string]string` associated with the http request `r`
	reqData := mux.Vars(r)
	storeID, err := strconv.Atoi(reqData["store_id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = t.TaskService.GetAll(storeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
