// This project is largely based on Eli Bendersky's series `REST Servers in Go`.
// https://eli.thegreenplace.net/2021/rest-servers-in-go-part-1-standard-library/

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/wagnojunior/taskstore/internal/taskstore"
)

func main() {
	// Starts Gorilla mux
	r := mux.NewRouter()
	r.StrictSlash(true)

	// Starts the task server
	server := taskstore.NewTaskServer()

	// Defines the routes
	r.HandleFunc("/task/", server.CreateTaskHandler).Methods(http.MethodPost)
	r.HandleFunc("/task/", server.GetAllTasksHandler).Methods(http.MethodGet)
	r.HandleFunc("/task/", server.DeleteAllTasksHandler).Methods(http.MethodDelete)
	r.HandleFunc("/task/{id:[0-9]+}/", server.GetTaskHandler).Methods(http.MethodGet)
	r.HandleFunc("/task/{id:[0-9]+}/", server.DeleteTaskHandler).Methods(http.MethodGet)
	r.HandleFunc("/tag/{tag}/", server.GetTaskByTagHandler).Methods(http.MethodGet)
	r.HandleFunc("/due/{year:[0-9]+}/{month:[0-9]+}/{day:[0-9]}/", server.GetTaskByDueHandler).
		Methods(http.MethodGet)

	fmt.Println("Starting server at port :3030...")
	log.Fatal(http.ListenAndServe("localhost:3030", taskstore.LoggingMiddleware(r.ServeHTTP)))
}
