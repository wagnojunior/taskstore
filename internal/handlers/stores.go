package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/wagnojunior/taskstore/internal/models"
)

type Store struct {
	StoreService *models.StoreService
}

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

	_, err = s.StoreService.Create(store.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}
