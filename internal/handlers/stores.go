package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/wagnojunior/taskstore/internal/models"
	"github.com/wagnojunior/taskstore/internal/repository"
)

type Store struct {
	StoreService repository.DBRepo
}

// `Create` handles POST requests to create a store
func (s Store) Create(w http.ResponseWriter, r *http.Request) {
	// Creates a new decoder and disallow unknown fields
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	// Decodes the JSON data to a Go struct
	var store models.Store
	err := dec.Decode(&store)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = s.StoreService.CreateStore(store.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

// `DeleteByID` handles GET request to delete a store
func (s Store) DeleteByID(w http.ResponseWriter, r *http.Request) {
	// Gets a `map[string]string` associated with the http request `r`
	reqData := mux.Vars(r)
	id, err := strconv.Atoi(reqData["store_id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.StoreService.DeleteStoreByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
